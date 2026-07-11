# 时间轴 · 项目架构文档

> 版本 1.2.0 · 个人时间管理与待办事项应用 · 移动端 + 平板优先

---

## 1. 项目概述

**时间轴 (Time)** 是一款面向个人的时间管理 Web 应用，支持：

- 时间轴记录：随记(太阳/阴雨)、运动、自律打卡、禁糖打卡、尿酸、资产、读书、电影、待办
- 待办事项：完整的 CRUD 操作，支持分类(category)和截止日期(due_date)
- DeepSeek AI 集成：读书自动获取书籍信息(作者/国籍/字数/出版日期/标签)，影视自动获取影视信息(导演/主演/国家/上映年份/时长/标签)
- 补卡回填：选择日期补录遗漏打卡，自动重算连续天数（面板保持打开，仅在确认/取消时关闭）
- 统计页面：GitHub 风格热力图 + 统计卡片（touch-action: pan-y，区分水平/垂直滑动）
- 面板排序：添加面板条目可拖拽排序，顺序持久化
- 数据导出/导入：GET /api/data/export + POST /api/data/import（JSON 快照，含 entries/todos/settings）
- 数据库迁移系统：启动时自动执行 `db/migrations/` 下未应用的 SQL 文件，通过 `schema_migrations` 表记录版本
- 运动卡路里计算：基于 MET 公式（69kg 体重 / 心率 160），覆盖跑步/骑行/俯卧撑/跳绳/游泳/徒步/自由活动

### 技术栈

| 层 | 技术 | 版本 |
|----|------|------|
| 前端 | Vue 3 (Composition API) | ^3.5.38 |
| 构建 | Vite | ^8.1.0 |
| 后端 | Go | 1.26.4 |
| 数据库 | SQLite (纯 Go) | modernc.org/sqlite v1.32.0 |
| 通信 | RESTful JSON API | 端口 8080 |
| 开发 | Vite HMR | 端口 5173 |

---

## 2. 项目结构

```
time/
├── index.html                     # SPA 入口 HTML
├── vite.config.js                 # Vite 配置 + recalculate 中间件
├── package.json                   # 前端依赖
├── README.md                      # 项目说明
├── src/                           # ── 前端源码 ──
│   ├── main.js                    # Vue 应用入口
│   ├── App.vue                    # 根组件：状态管理 + 路由控制
│   ├── config.js                  # 集中常量 (API_BASE)
│   ├── style.css                  # 全局样式 + CSS 变量 + 响应式断点
│   ├── utils/
│   │   └── date.js                # 日期格式化 (formatRecordTime)
│   ├── composables/
│   │   └── useTimePicker.js       # 时间选择器复用逻辑
│   └── components/
│       ├── YearMonthHeader.vue    # 顶栏：年月选择 + 下拉菜单
│       ├── AddButton.vue          # 底部 FAB "+" 按钮
│       ├── AddEntryPanel.vue      # 添加记录面板（含排序模式）
│       ├── TimelineEntryCard.vue  # 时间轴单条卡片（分类色标）
│       ├── TimelineEntryGroup.vue # 日期分组容器
│       ├── ThoughtFormPanel.vue   # 念头创建表单
│       ├── AssetFormPanel.vue     # 资产创建表单
│       ├── UricFormPanel.vue      # 尿酸创建表单
│       ├── ExerciseFormPanel.vue  # 运动创建表单
│       ├── ReadingFormPanel.vue   # 读书表单
│       ├── MovieFormPanel.vue     # 影视表单
│       ├── TodoFormPanel.vue      # 待办创建表单
│       ├── TimePickerModal.vue    # 通用时间选择弹窗
│       ├── ConfirmModal.vue       # 确认对话框（替代 alert）
│       ├── StatisticsPage.vue     # 统计页：热力图 + 卡片
│       └── TodoPage.vue           # 待办事项页
├── server/                        # ── Go 后端 ──
│   ├── main.go                    # 入口：配置加载 + 路由注册 + 启动
│   ├── config.json                # 运行时配置 (host, port, db_path, static_dir, migrations_dir)
│   ├── go.mod / go.sum            # Go 模块依赖
│   ├── db/
│   │   ├── schema.sql             # 建表 DDL（参考用，实际运行用迁移文件）
│   │   └── queries.sql            # 命名查询（参考用）
│   │   └── migrations/            # 数据库迁移文件
│   │       └── 001_init.sql       # 初始建表（entries/settings/todos）
│   ├── handler/
│   │   ├── entry.go               # Entry/Settings/Todo/Export/Import/Recalculate 请求处理
│   │   ├── todo.go                # Todo 请求处理（挂载在 EntryHandler）
│   │   ├── book.go                # DeepSeek 书籍信息查询
│   │   └── media.go               # DeepSeek 影视信息查询
│   ├── service/
│   │   ├── entry.go               # Entry/Settings/Recalculate 数据层
│   │   ├── migration.go           # RunMigrations + ExportAll/ImportAll
│   │   └── todo.go                # Todo 数据层
│   ├── middleware/
│   │   ├── cors.go                # CORS 中间件
│   │   └── logger.go              # 请求日志中间件
│   └── data/
│       └── time.db                # SQLite 数据库文件
└── docs/                          # ── 项目文档 ──
    ├── ARCHITECTURE.md            # 本文档
    ├── API.md                     # API 参考
    ├── DATABASE.md                # 数据库设计
    ├── DEVELOPMENT.md             # 开发指南
    ├── COMPONENTS.md              # 组件文档
    ├── DESIGN.md                  # 设计系统
    └── CHANGELOG.md               # 变更日志
```

---

## 3. 架构分片

```
┌──────────────┐     HTTP/JSON     ┌──────────────┐     SQL      ┌──────────┐
│  Vue 3 前端   │ ◄──────────────► │  Go 后端 API  │ ◄──────────► │  SQLite   │
│  端口 5173    │                  │  端口 8080    │              │  time.db  │
│  Vite HMR     │  Vite proxy:     │  net/http     │              │  CGo-free │
└──────────────┘  /api → :8080     └──────────────┘              └──────────┘
                                        │
                          POST /api/entries/recalculate
                          GET  /api/data/export
                          POST /api/data/import
                          (全部在 Go 后端实现)
```

### 分层架构（后端）

```
main.go              → 入口 + 配置 + 启动
  ├── handler/       → HTTP 请求解析 + 参数验证 + 响应封装
  │   ├── entry.go   → EntryHandler (条目 CRUD + 设置读写 + 导出导入 + 连续天数重算)
  │   ├── todo.go    → ListTodos/CreateTodo/UpdateTodo/DeleteTodo
  │   ├── book.go    → BookInfo (DeepSeek 书籍信息查询)
  │   └── media.go   → MediaInfo (DeepSeek 影视信息查询)
  ├── service/       → 业务逻辑 + 数据库操作
  │   ├── entry.go   → EntryService (分页/按月查询/创建/设置/连续天数重算)
  │   ├── todo.go    → Todo 数据层
  │   └── migration.go → RunMigrations + ExportAll/ImportAll
  └── middleware/    → CORS + Logger
```

**数据库迁移：** 启动时自动执行 `server/db/migrations/` 目录下未应用的 SQL 文件，通过 `schema_migrations` 表记录版本。`server/db/schema.sql` 作为参考文件保留。`config.json` 中 `migrations_dir` 字段指定迁移文件路径。

### 组件分层（前端）

```
App.vue                            → 根组件（状态 + 路由）
  ├── YearMonthHeader              → 只在水位页面显示
  ├── TimelineEntryGroup ×N        → 日期分组
  │   └── TimelineEntryCard ×N     → 单条卡片
  ├── StatisticsPage               → 统计页（全屏切换）
  ├── TodoPage                     → 待办页（全屏切换）
  ├── AddEntryPanel                → 底部弹出面板（排序 + 补卡按钮）
  ├── ThoughtFormPanel             → 随记表单
  ├── AssetFormPanel               → 资产表单
  ├── UricFormPanel                → 尿酸表单
  ├── ExerciseFormPanel            → 运动表单
  ├── ReadingFormPanel             → 读书表单
  ├── MovieFormPanel               → 影视表单
  ├── TodoFormPanel                → 待办创建表单
  ├── ConfirmModal                 → 确认弹窗
  └── TimePickerModal              → 底部弹出时间选择器（补卡用）
```

---

## 4. 数据库

四张表：`entries`（时间轴条目）、`todos`（待办）、`settings`（键值配置）、`schema_migrations`（迁移版本记录）。

详见 [DATABASE.md](./DATABASE.md)。

---

## 5. API 总览

| 方法 | 路径 | 处理函数 | 功能 |
|------|------|---------|------|
| GET | `/api/entries?limit=&before=` | `GetAllEntries` | 分页获取条目 |
| GET | `/api/entries?month=` | `GetEntries` | 按月查询条目 |
| POST | `/api/entries` | `CreateEntry` | 创建条目 (201) |
| DELETE | `/api/entries/{id}` | `DeleteEntry` | 删除条目 |
| POST | `/api/entries/recalculate?type=` | `RecalculateEntries` | 重算连续天数 |
| POST | `/api/book-info` | `BookInfo` | DeepSeek 书籍信息查询 |
| POST | `/api/media-info` | `MediaInfo` | DeepSeek 影视信息查询 |
| GET | `/api/data/export` | `ExportData` | 导出全量数据 |
| POST | `/api/data/import` | `ImportData` | 导入数据（清空后写入） |
| GET | `/api/settings/{key}` | `GetSetting` | 读取配置 |
| PUT | `/api/settings/{key}` | `SetSetting` | 写入配置 |
| GET | `/api/todos` | `ListTodos` | 待办列表 |
| POST | `/api/todos` | `CreateTodo` | 创建待办 (201) |
| PUT | `/api/todos/{id}` | `UpdateTodo` | 更新待办（标题/完成） |
| DELETE | `/api/todos/{id}` | `DeleteTodo` | 删除待办 |

> 所有端点均在 Go 后端实现。Vite dev 中间件额外保留 `/api/entries/recalculate` 用于直接读写 SQLite 的开发场景（`vite.config.js` 中的 `recalculatePlugin`），但生产环境由 Vite proxy（`/api → localhost:8080`）转发到 Go 后端处理。

详见 [API.md](./API.md)。

---

## 6. 核心数据流

### 6.1 创建条目

```
AddEntryPanel(@select)
  → App.handleSelect(item)
    → 根据 type 路由到对应 FormPanel
      → App.handleXxxCreate(data)
        → formatRecordTime(time) 格式化
        → POST /api/entries
        → loadInitial() 刷新时间轴
        → currentMonth = monthKey, showStats = false
```

### 6.2 快捷打卡（自律/禁糖）

```
AddEntryPanel(@select)
  → App.handleSelect({id:'discipline'})
    → getConsecutiveDay(type)
      → 检查今天是否已打卡（返回 -1 则弹确认窗）
      → 计算连续天数
    → POST /api/entries (description = 连续天数)
    → loadInitial()
```

### 6.3 补卡（回填）

```
AddEntryPanel(@backfill) → backfill 按钮点击
  → App.handleBackfill(item)
    → bfOpen() 弹出 TimePickerModal（AddEntryPanel 保持打开）
    → 用户调整时间（面板不关闭）
    → 用户点击"确定" → handleBackfillConfirm()
      或点击"取消" → bfCloseAndHide()（关闭选择器 + 面板）
    → handleBackfillConfirm:
      → GET /api/entries?limit=200 检查重复
      → 重复则 showConfirm = true
      → 不重复则 POST /api/entries
      → POST /api/entries/recalculate?type= 重算连续天数（Go 后端）
      → loadInitial()
      → showAddPanel = false
```

### 6.4 面板排序

```
AddEntryPanel
  → toggleSort() → sortMode = true
  → moveUp/moveDown 交换位置
    → saveOrder()
      → PUT /api/settings/panel-order
  下次 onMounted:
    → loadOrder()
      → GET /api/settings/panel-order
      → 解析 JSON 重新排列 items
```

### 6.5 统计页加载

```
StatisticsPage.onMounted
  → fetchAll()
    → while(true): GET /api/entries?limit=100&before=
    → 拼接 allEntries
  → categoryData(computed)
    → buildHeatmapWeeks(entries, year)
    → 构建 7 列周数组
  → 热力图自动滚动到末尾
```

### 6.6 滚动月份检测

```
window.addEventListener('scroll', handleScroll)
  → 遍历 [data-date] 元素，找第一个 top > 55 的
  → 更新 currentMonth → watch 触发 scrollIntoView
  → 接近底部时触发 loadMore()
```

### 6.7 DeepSeek 书籍查询（读书）

```
ReadingFormPanel → 用户输入书名
  → handleBookInfo(name)
    → POST /api/book-info { book_name }
    → Go handler/book.go → DeepSeek API
    → 返回 { author, nationality, word_count, first_publish_date, tags }
  → 自动填充表单字段
  → 用户确认后 POST /api/entries (type=reading)
```

### 6.8 DeepSeek 影视查询（电影/剧集/动漫）

```
MovieFormPanel → 用户输入影视名称 + 选择类型(movie|series|anime)
  → handleMediaInfo(name, type)
    → POST /api/media-info { name, type }
    → Go handler/media.go → DeepSeek API
    → 返回 { director, cast, country, first_release_date, duration|episodes, tags }
  → 自动填充表单字段
  → 用户确认后 POST /api/entries (type=movie)
```

### 6.9 待办创建

```
TodoFormPanel → 用户输入标题 + 选择分类 + 截止日期
  → handleTodoCreate(data)
    → POST /api/todos { title, category, due_date }
    → 创建成功后回首页
```

---

## 7. 响应式设计

### 断点

| 宽度 | 内容区最大宽 | 适用设备 |
|------|------------|---------|
| < 600px | 480px | 手机 |
| 600-768px | 544px | 大屏手机 / iPad mini 竖屏 |
| 768-960px | 640px | iPad mini 横屏 / 小平板 |
| > 960px | 720px | 桌面 |

所有 `max-width` 统一使用 CSS 变量 `var(--content-width)`，通过 `:root` 的 `@media` 规则驱动。

详见 [DESIGN.md](./DESIGN.md)。

---

## 8. 启动方式

```powershell
# 后端（需要 Go 1.26+）
cd server
go build -o time-server.exe .
.\time-server.exe

# 前端开发（需要 Node.js 18+）
npm install
npm run dev -- --host

# 生产构建
npm run build        # 输出到 dist/
npm run preview      # 预览构建产物
```

- 前端：http://localhost:5173
- 后端：http://localhost:8080
