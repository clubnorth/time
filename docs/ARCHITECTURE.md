# 时间轴 · 项目文档与技术架构

> 一个移动端优先的垂直时间轴应用，用于浏览和管理时序条目。支持多层级弹出式交互流程。

---

## 1. 项目概览

| 属性 | 值 |
|------|-----|
| 项目名 | `time` |
| 类型 | 单页应用 (SPA) |
| 技术栈 | Vue 3 + Vite |
| 语言 | JavaScript (ES Modules) |
| 运行时 | Node.js ≥ 18 |
| 包管理器 | npm |
| 入口文件 | `index.html` → `src/main.js` |
| 构建输出 | `dist/` |

### 核心依赖

```
vue              ^3.5.38     UI 框架 (Composition API)
vite             ^8.1.0      构建工具
@vitejs/plugin-vue ^6.0.7    Vue SFC 编译插件
```

---


---

## 1.5 技术架构总览

| 层 | 语言 | 框架/库 |
|---|---|---|
| 前端 | JavaScript (ES Modules) | Vue 3 (Composition API) + Vite |
| 后端 | Go | chi (HTTP 路由) + sqlc (查询代码生成) |
| 数据库 | SQLite | modernc.org/sqlite (纯 Go 驱动, 无 CGO) |
| API | RESTful JSON | 前端 fetch / 后端 encoding/json |
| HTTP 客户端 | 原生 fetch | 不引入 axios |

**选型理由**：
- Go 编译为单文件二进制，部署只需复制一个文件
- chi 完全兼容 `net/http` 标准库接口，无黑盒行为
- SQLite 零配置、单文件，个人使用无并发瓶颈；modernc 纯 Go 实现，编译不依赖 CGO
- sqlc 从手写 SQL 自动生成类型安全的 Go 代码，AI 可直读 SQL，无 ORM 隐式查询
- 前端继续当前 Vue 3 + Vite + JavaScript，不做框架升级

### 数据库

单表设计，覆盖全部 12 种条目类型：

```sql
CREATE TABLE entries (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    type        TEXT    NOT NULL,
    title       TEXT    NOT NULL DEFAULT '',
    description TEXT    NOT NULL DEFAULT '',
    category    TEXT    NOT NULL DEFAULT '',
    valence     TEXT,
    recorded_at TEXT    NOT NULL,
    created_at  TEXT    NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at  TEXT    NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX idx_entries_type      ON entries(type);
CREATE INDEX idx_entries_recorded  ON entries(recorded_at);
```

### API

统一响应格式：`{ "code": 0, "message": "ok", "data": {...} }`

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/entries?month=2025-06` | 获取当月条目 |
| POST | `/api/entries` | 创建条目 |
| PUT | `/api/entries/:id` | 编辑条目 |
| DELETE | `/api/entries/:id` | 删除条目 |

### 后端项目结构（待实现）

```
server/
├── main.go
├── config.json
├── db/
│   ├── schema.sql
│   ├── queries.sql
│   └── sqlc.yaml
├── handler/
│   └── entry.go
├── service/
│   └── entry.go
└── middleware/
    ├── cors.go
    └── logger.go
```

### 前端改动（待实现）

```
src/
├── api/
│   ├── client.js          # fetch 封装
│   └── entries.js         # CRUD 函数
├── App.vue                # 替换 mock 数据为 API 调用
└── components/            # 不变
```

后端 `handler → service → db(sqlc)` 三层分离，前端仅新增 `api/` 目录替换 App.vue 数据源，组件层零改动。

## 2. 项目结构

```
time/
├── index.html
├── package.json
├── vite.config.js
├── docs/
│   └── ARCHITECTURE.md
├── public/
│   ├── favicon.svg
│   └── icons.svg
└── src/
    ├── main.js
    ├── style.css
    ├── App.vue
    └── components/
        ├── YearMonthHeader.vue      # 顶部固定年月选择器
        ├── TimelineEntryGroup.vue   # 日期聚合块
        ├── TimelineEntryCard.vue    # 时刻条目卡片
        ├── AddButton.vue            # 底部加号按钮
        ├── AddEntryPanel.vue        # 第二层：条目类型选择面板
        └── ThoughtFormPanel.vue     # 第三层：念头录入表单
```

---

## 3. 组件架构

### 3.1 组件树

```
App.vue
├── YearMonthHeader.vue         (fixed top)
├── TimelineEntryGroup.vue × N  (grid items, display:contents)
│   └── TimelineEntryCard.vue × M
├── AddButton.vue               (fixed bottom)
├── AddEntryPanel.vue           (bottom sheet overlay, 80% height)
└── ThoughtFormPanel.vue        (bottom sheet overlay, 78% height, on top)
```

### 3.2 页面交互流程

```
┌───────────┐  点击底部+   ┌───────────────┐  点击念头正/负  ┌──────────────────────┐
│  第一页    │ ──────────→ │   第二页       │ ──────────────→ │  第三页               │
│  时间轴    │             │  条目类型选择   │                 │  念头录入表单         │
│           │ ←── 新建 ─── │  (12种类型)    │ ←── 取消 ────── │  时间选择器 + 备注    │
└───────────┘             └───────────────┘                 └──────────────────────┘
                                                                    │
                                             新建 (同时关闭第二/三页) ↓
                                                              ┌───────────┐
                                                              │  第一页    │
                                                              │  (回到)   │
                                                              └───────────┘
```

**关键交互逻辑**：

| 操作 | 行为 |
|------|------|
| 点击底部 + | 打开第二页（条目类型列表，80% 高度） |
| 点击第二页空白区域 / X / 选择非念头条目 | 关闭第二页，回第一页 |
| 点击念头条目右侧「正」/「负」 | 第二页保持打开，第三页弹出覆盖上方（78% 高度），标题自动为「念头·正」或「念头·负」 |
| 点击第三页空白区域 / 「取消」 | 关闭第三页，回到第二页（第二页仍然开着） |
| 第三页「新建」 | 同时关闭第二页和第三页，回到第一页 |
| 第二/三页堆叠 | 第二页 80%、第三页 78%，都从底部弹出，第三页上方露出一条第二页的圆角边缘 |

### 3.3 组件职责

#### `App.vue`
- **职责**：根容器、数据层、滚动交互逻辑、面板状态管理
- **数据**：持有 `allTimelineData`（6 个月 mock 数据）、`availableMonths`（月份选项）
- **面板状态**：`showAddPanel`（第二页）、`showThoughtForm` + `thoughtKind`（第三页）
- **关键函数**：
  - `handleAdd()` → 打开第二页
  - `handleSelect(item)` → 若 `item.kind` 存在则打开第三页（不改第二页），否则关闭第二页
  - `handleThoughtCreate(data)` → 关闭第二/三页，回第一页
- **布局**：Grid `46px 4px 20px 9px 1fr`，竖线 `left: 84px`

#### `AddEntryPanel.vue`（第二页）
- **职责**：弹出式条目类型选择面板
- **位置**：`position: fixed; z-index: 1000`，底部弹出，占视口 80%
- **条目列表**：12 种类型（念头、备忘录、日记、待办事项、读书、电影、剧集、运动、自律、尿酸、禁止糖分、户外时间）
- **布局**：每行 = 彩色圆形图标(40px) + 标题/副标题文字 + 右侧操作按钮
- **念头特殊处理**：右侧为两个胶囊形按钮「正」(绿底) 和「负」(红底)，其余条目为圆形 + 按钮
- **颜色系统**（低饱和度）：
  - 念头 `#9DB5C9`、备忘录 `#D4B87A`、日记 `#8CAD8C`、待办 `#CF8B86`
  - 读书 `#A099C4`、电影 `#CB99B0`、剧集 `#D4A882`、运动 `#84B8A4`
  - 自律 `#8CA4BD`、尿酸 `#BE999B`、禁止糖分 `#C88C8C`、户外 `#92B4C7`
- **关闭方式**：点击遮罩 / 右上角 X / 选择条目
- **动画**：底部滑入，opacity + translateY transition

#### `ThoughtFormPanel.vue`（第三页）
- **职责**：念头正/负录入表单
- **位置**：`position: fixed; z-index: 1100`，底部弹出，占视口 78%（比第二页矮 2%，产生堆叠）
- **Props**：`visible` (Boolean), `kind` ('positive' | 'negative')
- **内容模块**：
  - 顶部导航栏：左侧「取消」→ 关闭第三页；居中标题「念头·正/负」；右侧「新建」(备注为空时禁用)
  - 记录时间：圆角点击区域（居中显示时间），点击弹出自定义时间选择器
  - 备注输入：大尺寸 textarea，placeholder「请输入」，200 字上限，右下角实时计数
- **自定义时间选择器**（picker modal）：
  - 居中弹出卡片（340px 宽，16px 圆角），半透明遮罩
  - 5 列：年/月/日/时/分，每列独立 ▲▼ 按钮
  - 月份 1↔12 自动进位年，日期根据当月天数自适应
  - 底部居中「确定」按钮（胶囊形）
  - 动画：scale 弹出/消失

### 3.4 原有组件

#### `YearMonthHeader.vue`
- **职责**：顶部固定年月显示 + 下拉月份切换
- **位置**：`position: fixed; top: 0; max-width: 480px; z-index: 50`
- **Props**：`modelValue`, `months`

#### `TimelineEntryGroup.vue`
- **职责**：日期聚合容器（星期 + 圆点 + 月日 + 条目列表）
- **布局**：`display: contents`
- **颜色**：周一到周四 `#3a3a3a`，周五 `#f5a623`，周末 `#4caf50`

#### `TimelineEntryCard.vue`
- **职责**：单条时刻条目（时间 + 空心圆环 + 内容卡片）
- **分类标记**：3×14px 竖线（黄/红/绿）

#### `AddButton.vue`
- **职责**：底部加号按钮
- **样式**：52px 深色圆形，白边，半露出于底部栏上

---

## 4. 数据流

```
┌──────────────────────────────────────────────────┐
│  allTimelineData (static mock)                    │
│  { '2025-01': [...], '2025-02': [...], ... }      │
└──────────────────┬───────────────────────────────┘
                   ▼
         ┌────────────────┐
         │   allGroups     │  computed: 全部月份倒序合并
         └────────┬───────┘
                  │ v-for
                  ▼
    ┌──────────────────────────┐
    │  TimelineEntryGroup × N  │
    │  └─ TimelineEntryCard × M│
    └──────────────────────────┘
```

### 面板状态流转

```
                 showAddPanel      showThoughtForm      thoughtKind
初始                 false             false              'positive'
点击底部 +           true              false              'positive'
点击念头正            true              true               'positive'
点击念头负            true              true               'negative'
点击第三页取消         true              false              (不变)
点击第三页新建         false             false              (不变)
点击第二页X/空白       false             false              (不变)
```

### 月份自动切换流程（不变）

1. `window scroll` → `handleScroll()` 遍历 `[data-date]`
2. 找第一个 `rect.top > 55px` 的圆点
3. 提取月份更新 `currentMonth`（`skipScrollWatch` 防止回环）
4. `YearMonthHeader` 通过 `v-model` 自动响应

---

## 5. 布局体系

### 5.1 CSS Grid 列定义（不变）

```
grid-template-columns: 46px 4px 20px 9px 1fr
```
竖线在 `left: 84px`（第 3 列中心）。

### 5.2 Z-Index 层级（更新）

| 层级 | 元素 |
|------|------|
| 0 | 时间轴竖线 `.timeline-line` |
| 1 | Grid 常规内容 |
| 2 | 圆点 `.solid-dot`, `.hollow-dot` |
| 50 | `YearMonthHeader`、`.bottom-zone` |
| 1000 | `AddEntryPanel` 遮罩 + 面板（第二页） |
| 1100 | `ThoughtFormPanel` 遮罩 + 面板（第三页） |
| 1200 | 自定义时间选择器 popup |

### 5.3 弹出面板布局

```
AddEntryPanel (第二页):
  .panel-overlay: fixed inset:0, z-index:1000
  .panel-sheet:  height:80%, 底部对齐, border-radius:20px 20px 0 0

ThoughtFormPanel (第三页):
  .form-overlay: fixed inset:0, z-index:1100
  .form-sheet:   height:78%, 底部对齐, border-radius:16px 16px 0 0
  → 比第二页矮 2%，上方露出第二页圆角 → 堆叠效果
```

---

## 6. 色彩系统

### 6.1 主色调

| 用途 | 色值 |
|------|------|
| 页面背景 | `#ffffff` |
| 卡片背景 | `#fafafa` |
| 正文文字 | `#1a1a1a` / `#2c2c2c` |
| 辅助文字 | `#666` / `#888` |
| 浅灰文字 | `#999` / `#bbb` |
| 竖线 / 边框 | `#d0d0d0` / `#ebebeb` / `#f0f0f0` |
| 遮罩 | `rgba(0,0,0,0.25)` |

### 6.2 面板色彩

| 用途 | 色值 |
|------|------|
| 面板底色 | `#fafafa` / `#f5f5f7` |
| 卡片底色 | `#ffffff` |
| 卡片边框 | `#f0f0f0` |
| 新建按钮 | `#3b82f6` (蓝色) |
| 加号按钮(第二页) | `#f0f0f0` 底 + `#555` 加号 |
| 正按钮(念头) | `#e6f0e6` 底 + `#4a7c4a` 文字 |
| 负按钮(念头) | `#f0e6e6` 底 + `#7c4a4a` 文字 |
| 确定按钮(时间) | `#2c2c2c` 底 + `#ffffff` 文字 |

### 6.3 条目图标色（12 类，低饱和度）

| 条目 | 色值 |
|------|------|
| 念头 | `#9DB5C9` |
| 备忘录 | `#D4B87A` |
| 日记 | `#8CAD8C` |
| 待办事项 | `#CF8B86` |
| 读书 | `#A099C4` |
| 电影 | `#CB99B0` |
| 剧集 | `#D4A882` |
| 运动 | `#84B8A4` |
| 自律 | `#8CA4BD` |
| 尿酸 | `#BE999B` |
| 禁止糖分 | `#C88C8C` |
| 户外时间 | `#92B4C7` |

---

## 7. 交互功能清单

| 功能 | 状态 | 实现 |
|------|------|------|
| 年份/月份下拉切换 | ✅ | `YearMonthHeader` v-if 下拉 |
| 滚动自动切换月份 | ✅ | `window scroll` + `data-date` 检测 |
| 点击月份跳转 | ✅ | `watch(currentMonth)` + `scrollIntoView` |
| 数据倒序 | ✅ | `allGroups` computed |
| 日期圆点分色 | ✅ | 周五黄 / 周末绿 / 其他深灰 |
| 分类竖线标记 | ✅ | 3×14px 药丸竖线 |
| 底部加号按钮 | ✅ | 52px 深色圆，半露出 |
| 第二页：条目类型选择 | ✅ | `AddEntryPanel` 80% 底部弹出 |
| 12 种条目类型 + 图标 | ✅ | 彩色圆图标 + SVG 符号 |
| 念头正/负胶囊按钮 | ✅ | 绿底「正」/ 红底「负」 |
| 第三页：念头录入表单 | ✅ | `ThoughtFormPanel` 78% 底部弹出 |
| 页面堆叠效果 | ✅ | 第三页比第二页矮 2% |
| 自定义时间选择器 | ✅ | 5 列年/月/日/时/分 + 确定按钮 |
| 备注 200 字限制 | ✅ | textarea maxlength + 实时计数 |
| 新建按钮空禁用 | ✅ | `:disabled="!note.trim()"` |
| 移动端适配 | ✅ | max-width: 480px |
| 桌面端手机壳 | ✅ | `@media (min-width: 481px)` 灰底+阴影 |

---

## 8. 开发命令

```bash
npm install          # 安装依赖
npm run dev          # 启动开发服务器 (http://localhost:5173)
npm run build        # 生产构建 → dist/
npm run preview      # 预览生产构建
```

---

## 9. 后续扩展

- **后端接入**：`allTimelineData` → API 调用，数据库存储
- **数据库选型**：推荐 SQLite（Go + modernc.org/sqlite 纯静态编译），迁移用版本号管理
- **Vue Router**：如多页面引入路由
- **状态管理**：数据复杂度提升后引入 Pinia
- **PWA**：Service Worker 离线支持

---

## 2026-07-01 迭代记录

### 后端
- **游标分页**：`GET /api/entries?limit=30&before=<recorded_at>` 支持无限滚动
- **Content-Type**：`application/json; charset=utf-8` 修复中文乱码
- **通用表设计**：`entries` 表通过 `type` 字段支持全部条目类型（thought / asset / uric / exercise / discipline / nosugar），无需额外建表
- **Go 环境**：安装至 `D:\ProgramData\go`，`CGO_ENABLED=0` 静态编译

### 前端新增组件
| 组件 | 功能 |
|------|------|
| AssetFormPanel.vue | 资产记录录入，数字输入框 + 时间选择器 |
| UricFormPanel.vue | 尿酸记录录入，数字输入框 + 时间选择器 |
| ExerciseFormPanel.vue | 运动记录录入，胶囊按钮选择运动类型 + 数字输入 |

### 页面叠加
- 资产/尿酸/运动录入面板均为 78% 高度，叠加在 AddEntryPanel（80%）之上，露出圆角边缘
- z-index: AddEntryPanel 1000 → 录入面板 1100 → 时间选择器 1200

### 快捷打卡
- **自律** / **禁止糖分**：绿色对勾按钮，一键打卡
- **连续天数**：自动计算，中断归零，同一天不可重复打卡

### 卡片显示
| 类型 | 标题 | 描述格式 |
|------|------|----------|
| 念头 | 念头·正/负 | 备注文字 |
| 资产记录 | 资产记录 | 你现在的余额是 [彩虹数字] 元 |
| 尿酸记录 | 尿酸记录 | 你今天的尿酸值是 XX mol |
| 运动 | 运动·跑步 | 你今天跑步跑了 X 千米 本次消耗 XX 卡 |
| 自律 | 自律 | 今天是连续自律第 X 天 |
| 禁止糖分 | 禁止糖分 | 今天是连续无糖第 X 天 |

### 运动卡路里
- 基于 MET 公式，体重 70kg：
  - 跑步 km×56 / 骑行 km×21 / 俯卧撑 个×0.15 / 跳绳 个×0.12
  - 游泳 m×0.28 / 徒步 km×56 / 自由活动 min×3.5
- 卡路里存入 `valence` 字段，卡片追加"本次消耗 XX 卡"

### 底部栏
- 左侧图标+文字"统计"，右侧图标+文字"待办事项"，中间 + 号按钮
- 样式待实现

### UI 修复
- 卡片 `word-break: break-word` 防止长文本溢出
- 时间箭头 `position: absolute` 右对齐，文字居中
- 数字输入框按钮解除 `.trim()` bug（number 类型改用 `String()` 包裹）
- 全部新建成功后自动关闭第二页回到时间轴
- .rainbow 使用 `:deep()` 穿透 scoped CSS 到 v-html 内容

### 已知问题
- 部分组件中文偶现乱码（PowerShell `Set-Content -Encoding UTF8` BOM 污染），需用 Node.js `WriteAllText(noBOM)` 重建


---

## 2026-07-02 迭代记录

### 统计页面
- 新增 StatisticsPage.vue，点击底部栏"统计"按钮进入
- 顶部 2x2 网格展示：总记录 / 本月 / 今日 / 最长连续天数
- 周/月/年分段切换标签，年份箭头导航不可超过当前年
- 日历热力图：7行(周日到周六) x N列(周)，20px方块间距3px
  - 最左 = 该年第一周，最右 = 当前周
  - 月初方块标注"一""二"等，月末方块标注日期数字
  - 有记录 = 分类颜色，无记录 = 浅灰，今天 = 深色边框
  - 年份不可选未来年份

### 底部栏
- 移除统计页的 YearMonthHeader
- 底部栏在统计页保留显示
- 统计/待办事项按钮：图标在上文字在下，居中排列

### 快捷打卡
- 自律/禁止糖分：一天限打卡一次
- getConsecutiveDay 返回 -1 表示今天已打卡

### 运动卡路里
- MET 公式（体重70kg），存入 valence 字段
- 卡片追加"本次消耗 XX 卡"

### 已知问题 & 待解决

1. **热力图滚动**：scrollbar-width: none 隐藏滚动条后浏览器不再转发拖拽手势给滚动容器，需自行实现 pointerdown/pointermove 手动模拟 scrollLeft。

2. **溢出链**：overflow-x: auto 生效需要父容器用 overflow: hidden 约束宽度，否则子内容撑开父级导致无溢出。

3. **v-if + v-for 同元素**：Vue 3 中 v-if 优先级高于 v-for，需用 template 标签包裹。

4. **scoped CSS + v-html**：.rainbow 类不被 scoped 选择器匹配（v-html 内容无 data-v-xxx 属性），需用 :deep() 穿透。

5. **input type=number**：v-model 绑定返回 Number 类型，.trim() 报错，需 String().trim()。

6. **中文乱码**：PowerShell Set-Content -Encoding UTF8 添加 BOM 导致文件中文损坏，需用 Node.js WriteAllText(无BOM) 或纯 ASCII \u 转义。

7. **CRLF/LF**：Windows CRLF 换行导致正则 \n 不匹配，需用实际 byte sequence 或统一 LF。
