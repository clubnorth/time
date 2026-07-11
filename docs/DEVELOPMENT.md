# 时间轴 · 开发指南

---

## 环境要求

| 工具 | 最低版本 |
|------|---------|
| Node.js | 18+ |
| Go | 1.26+ |
| npm | 9+ |

---

## 快速开始

```powershell
# 1. 克隆项目
git clone <repo-url> time
cd time

# 2. 安装前端依赖
npm install

# 3. 启动后端（Go）
cd server
go build -o time-server.exe .
.\time-server.exe

# 4. 新终端，启动前端
cd ..
npm run dev -- --host
```

- 前端：http://localhost:5173
- 后端：http://localhost:8080

---

## 项目结构详解

### 前端 (`src/`)

```
src/
├── main.js              # createApp + mount，纯入口
├── App.vue              # 根组件：所有状态 + 所有事件处理 + 路由控制
├── config.js            # API_BASE = 'http://localhost:8080'
├── style.css            # 全局 CSS 变量 + reset + 响应式断点 + 共享动画
├── utils/
│   └── date.js          # formatRecordTime(Date) → {recordedAt, monthKey}
├── composables/
│   └── useTimePicker.js # 可复用时间选择器逻辑（currentTime, pickers, adjusters）
└── components/
    └── ...               # 14 个 Vue SFC 组件
```

**关键设计决策：**

- **无 Vue Router**：通过 `showStats` / `showTodo` 等 ref 控制页面切换，用 `v-if` 实现伪路由
- **无 Pinia**：所有状态集中在 `App.vue` 中，通过 props/emits 传递
- **Teleport to="body"**：所有弹窗/面板通过 Teleport 挂载到 body，避免 z-index 和 overflow 问题
- **CSS 变量响应式**：通过 `:root` 上的 CSS 变量 `--content-width` 配合 `@media` 实现响应式

### 后端 (`server/`)

```
server/
├── main.go              # 入口 + 配置 + 路由
├── config.json          # {"host","port","db_path","static_dir","migrations_dir"}
├── handler/
│   ├── entry.go         # EntryHandler: GetEntries, GetAllEntries, CreateEntry, GetSetting, SetSetting
│   │                    #              RecalculateEntries, ExportData, ImportData
│   └── todo.go          # 挂载在 EntryHandler: ListTodos, CreateTodo, UpdateTodo, DeleteTodo
├── service/
│   ├── entry.go         # EntryService: SQL 查询封装 + RecalculateEntries
│   ├── migration.go     # RunMigrations + ExportAll/ImportAll
│   └── todo.go          # EntryService: Todo SQL 封装
├── middleware/
│   ├── cors.go          # CORS: Allow-Origin: *
│   └── logger.go        # 请求耗时日志
├── db/
│   ├── schema.sql       # 建表 DDL（参考用）
│   ├── queries.sql      # 命名查询（未使用，预留 sqlc 支持）
│   └── migrations/      # 版本化迁移文件
│       └── 001_init.sql # 初始建表
└── data/
    └── time.db          # SQLite 数据库
```

**关键设计决策：**

- **Todo handler 挂载在 EntryHandler**：共享同一个 `*service.EntryService` 实例，避免额外的 struct
- **SQLite CGo-free**：使用 `modernc.org/sqlite`，无需 C 编译器
- **Schema 自动迁移**：启动时执行 `db/migrations/` 下未应用的 SQL 文件，通过 `schema_migrations` 表跟踪版本。配置项 `migrations_dir` 指定迁移目录。

---

## 开发约定

### 前端

**命名：**
- 组件文件：PascalCase（`TimelineEntryCard.vue`）
- 工具/composable：camelCase（`useTimePicker.js`）
- CSS class：kebab-case（`.timeline-container`）
- Vue transition name：kebab-case verb（`name="form"`）

**组件结构顺序：**
1. `<template>` — 模板
2. `<script setup>` — Composition API
3. `<style scoped>` — 组件样式

**事件命名：**
- 子→父：kebab-case verb（`@create`, `@close`, `@select`）
- 避免 `@click` 等原生事件穿透组件

**状态管理：**
- 跨组件共享状态放在 `App.vue`
- 单一组件局部状态用 `ref()` / `reactive()`
- 可复用逻辑提取到 `composables/`

### 后端

**命名：**
- 包名：小写单词（`handler`, `service`, `middleware`）
- 导出类型：PascalCase（`EntryHandler`, `EntryService`）
- 导出函数：PascalCase（`NewEntryService`）
- 私有变量/函数：camelCase（`respond`, `respondError`）

**文件组织：**
- 一个文件 = 一个职责
- handler 负责请求解析和响应
- service 负责数据库操作
- middleware 负责横切关注点

**错误处理：**
- Service 层返回 Go error
- Handler 层转为 `{code:1, message: err.Error()}` JSON
- 不吞错误，不使用 `panic`

---

## 常用命令

```powershell
# 前端
npm run dev          # 开发模式（HMR）
npm run build        # 生产构建 → dist/
npm run preview      # 预览构建产物

# 后端
cd server
go run .             # 开发运行
go build -o time-server.exe .   # 编译
go vet ./...         # 静态检查

# Git
git log --oneline -10
git diff --stat
```

---

## 添加新条目类型

以添加"喝水记录"为例：

### 后端
1. `server/db/migrations/` — 无需修改（通用字段已足够，如有新字段则新增迁移文件）
2. 无需修改 handler/service（`POST /api/entries` 接受任意 type）

### 前端
1. `src/components/AddEntryPanel.vue` — 在 `items` 数组中添加条目定义
2. 添加对应 SVG 图标（`v-if="item.id === 'water'"`）
3. `src/App.vue` — 在 `handleSelect()` 中处理新 type
4. 创建 `WaterFormPanel.vue`（或复用现有表单模式）
5. `src/App.vue` — 在 `allGroups` computed 中添加 description 格式化逻辑

---

## 运动卡路里计算公式

基于 MET（代谢当量）公式，参数：69kg 体重，平均心率 160。

```
kcal = MET × 体重(kg) × 时间(小时)
```

| 运动类型 | 单位 | kcal/单位 | MET 估算 |
|---------|------|----------|---------|
| 跑步 | km | 76 | 11.0 (10 km/h) |
| 骑行 | km | 28 | 9.0 (22 km/h) |
| 俯卧撑 | 个 | 0.17 | 4.5 (30个/min) |
| 跳绳 | 个 | 0.11 | 11.5 (120个/min) |
| 游泳 | m | 0.18 | 8.0 (50m/min) |
| 徒步 | km | 90 | 6.5 (5 km/h) |
| 自由活动 | min | 9.2 | 8.0 |

代码位置：`src/App.vue` → `handleExerciseCreate()` → `metCalcs` 对象。

---

## DeepSeek AI 集成

通过 DeepSeek API 实现书籍信息和影视信息的自动查询。

### API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/book-info` | 查询书籍信息（ISBN 或书名） |
| POST | `/api/media-info` | 查询影视信息（电影/剧集/动漫） |

### 配置

在 `server/config.json` 中配置 `deepseek_key`：

```json
{
  "host": "localhost",
  "port": 8080,
  "db_path": "data/time.db",
  "static_dir": "../dist",
  "migrations_dir": "db/migrations",
  "deepseek_key": "sk-xxxxxxxxxxxxxxxx"
}
```

后端通过 DeepSeek API 调用，返回结构化 JSON，前端解析后填充表单字段（标题、作者/导演、封面图等）。

---

## 评分系统

所有记录类型（读书、影视）均支持 1-10 分评分。

### 星星渲染规则

- 每 2 分 = 1 颗完整星星
- 奇数分值 = 最后半颗星（如 7 分 → ★★★½）
- 最大 5 颗星（10 分）

### 交互

可点击星星选择分值，点击第 N 颗星左侧为整数分，右侧为半星（`.5` 偏移检测）。

---

## 读书记录工作流

### 状态生命周期

```
在读 (reading) ──完成──▶ 已读 (read)
```

1. **创建记录**：在 `ReadingFormPanel` 中输入书名或 ISBN，点击 DeepSeek 查询按钮自动填充书籍信息
2. **状态选择**：通过胶囊按钮切换"在读"或"已读"
3. **分类**：7 种分类胶囊（文学、历史、哲学、科幻、技术、心理学、其他）
4. **评分**：完成阅读后通过星星组件评分（1-10）
5. **已读后**：可在时间轴中查看格式化条目（`v-html` 渲染封面图 + 书名 + 作者）

### 后端数据

`POST /api/book-info` 返回 `{title, author, cover, isbn, publisher, summary}`，前端存入 entry 的 `description` 字段（JSON 序列化）。

---

## 影视记录工作流

### 类型

| 类型 | 胶囊色 | 说明 |
|------|--------|------|
| 电影 | 玫瑰 `#CB99B0` | 单部影片 |
| 剧集 | 琥珀 `#D4A882` | 多季电视剧 |
| 动漫 | 紫色 `#A099C4` | 动画/漫画改编 |

### 观看状态

通过胶囊按钮切换：想看 → 在看 → 已看 → 弃剧

- `POST /api/media-info` 查询影视信息，前端进行 watch-status 校验（已看才能评分）
- 只有状态为"已看"时才启用评分功能

### 后端数据

`POST /api/media-info` 返回 `{title, director, year, cover, type, rating, summary}`，前端存入 entry 的 `description` 字段。

---

## 数据导出/导入工作流

### 日常备份

```bash
curl http://localhost:8080/api/data/export > backup_$(date +%Y%m%d).json
```

### 从生产同步到开发

```bash
# 服务器导出
ssh user@server "curl http://localhost:8080/api/data/export" > prod_data.json

# 开发机导入（清空本地数据，写入导入数据）
curl -X POST http://localhost:8080/api/data/import \
  -H "Content-Type: application/json" \
  -d @prod_data.json
```

### Schema 变更时重建数据

```bash
# 1. 导出
curl http://localhost:8080/api/data/export > backup.json

# 2. 停服 + 清库
# 停止服务，删除 server/data/time.db

# 3. 重启（自动执行全部迁移，创建新库）
# 启动服务

# 4. 导入
curl -X POST http://localhost:8080/api/data/import \
  -H "Content-Type: application/json" \
  -d @backup.json
```

---

## 调试技巧

- **API 调试**：使用 `curl` 或浏览器 DevTools Network 面板
- **数据库查看**：`sqlite3 server/data/time.db` → `.tables` → `SELECT * FROM entries LIMIT 5;`
- **热力图调试**：在 `StatisticsPage.vue` 中 `console.log(weeks)` 查看构建结果
- **CSS 变量调试**：DevTools Elements → Computed 面板查看 `--content-width` 实际值

---

## 常见问题

**Q: Go 编译报 "version does not match"**
A: 更新 `go.mod` 中 `go 1.xx` 为当前安装的 Go 版本。

**Q: 前端请求 CORS 报错**
A: 确认后端已启动在 8080 端口。

**Q: 数据库文件在哪里？**
A: `server/data/time.db`。删除此文件即可重置所有数据。
