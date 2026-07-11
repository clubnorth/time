# 时间轴 · 组件文档

> Vue 3 SFC · Composition API · scoped CSS

---

## 组件树

```
App.vue (根)
├── YearMonthHeader         v-if=!showStats && !showTodo
├── TimelineEntryGroup ×N   v-if=!showStats && !showTodo
│   └── TimelineEntryCard ×N
├── StatisticsPage          v-if=showStats
├── TodoPage                v-if=showTodo
├── AddEntryPanel           v-if=showAddPanel
├── ThoughtFormPanel        v-if=showThoughtForm
├── AssetFormPanel          v-if=showAssetForm
├── UricFormPanel           v-if=showUricForm
├── ExerciseFormPanel       v-if=showExerciseForm
├── ReadingFormPanel        v-if=showReadingForm
├── MovieFormPanel          v-if=showMovieForm
├── TodoFormPanel           v-if=showTodoForm
├── ConfirmModal            v-if=showConfirm
└── TimePickerModal         v-if=bfShowPicker
```

---

## App.vue

**角色：** 根组件，应用控制器。

**职责：**
- 时间轴数据获取与状态管理（`timelineData`, `availableMonths`, `currentMonth`）
- 页面路由控制（`showStats`, `showTodo`, `showAddPanel` 等）
- 所有创建/补卡/查询事件处理
- 滚动监听：月份检测 + 无限加载
- 日期分组计算（`allGroups`）

**关键数据：**

| ref / computed | 类型 | 说明 |
|---------------|------|------|
| `timelineData` | array | 原始条目列表 |
| `allGroups` | computed | 按日期分组的卡片渲染数据 |
| `currentMonth` | string | 当前选中月份 `YYYY-MM` |
| `availableMonths` | array | 可选月份列表 |
| `hasMore` | boolean | 是否还有更多数据 |
| `isLoading` | boolean | 加载中状态 |

---

## YearMonthHeader.vue

**角色：** 固定在顶部的年月选择器。

| Prop | 类型 | 说明 |
|------|------|------|
| `modelValue` | String | 当前月份 `YYYY-MM` |
| `months` | Array | 可选月份列表 `[{label, value}]` |

| Emit | 参数 | 说明 |
|------|------|------|
| `update:modelValue` | value | 选中月份变化 |

**状态：** `showDropdown` — 控制下拉菜单显隐。

---

## TimelineEntryGroup.vue

**角色：** 日期分组容器。为每个日期渲染日期标签 + 该日所有条目卡片。

| Prop | 类型 | 说明 |
|------|------|------|
| `group` | Object | `{ date, weekday, dateNum, entries[] }` |

**子组件：** `TimelineEntryCard`

**计算属性：** `dotClass` — 根据星期返回不同颜色（周五=金色，周末=绿色，其他=默认）

**DOM 标记：** 在 `solid-dot` 上设置 `data-date` 属性，供 `App.vue` 的滚动监听使用。

---

## TimelineEntryCard.vue

**角色：** 时间轴单条卡片。展示时间、分类色点、标题、描述。

| Prop | 类型 | 说明 |
|------|------|------|
| `entry` | Object | `{ time, title, description, category, isAsset, isUric }` |

**视觉特征：**
- 卡片左侧 2px 色标（`.cat-green` / `.cat-red` / `.cat-yellow`）
- 6px 圆形分类点
- 资产类型使用 `v-html` 渲染彩虹渐变金额
- 读书/影视类型使用 `v-html` 渲染封面图 + 结构化信息
- CSS 动画：`@keyframes rainbow-flow`
- 长按删除：长按卡片弹出确认 popover，支持删除操作

---

## AddEntryPanel.vue

**角色：** 底部弹出面板，列出所有可记录条目类型，支持排序模式和补卡。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |

| Emit | 参数 | 说明 |
|------|------|------|
| `close` | — | 关闭面板 |
| `select` | item | 选中条目（含 `kind` 字段区分正负向） |
| `backfill` | item | 触发补卡（打开 TimePickerModal，面板保持打开） |

**两种模式：**
1. **普通模式** (`sortMode = false`)：显示条目列表，每行带操作按钮
   - 随记：太阳/阴雨图标按钮（正向/负向）
   - 自律/禁糖：日历补卡按钮 + 绿色勾号打卡按钮
   - 待办：跳转 TodoFormPanel 创建待办事项
   - 其他（读书/影视/运动/资产/尿酸）：灰色圆形 + 号按钮
2. **排序模式** (`sortMode = true`)：显示圆形 chevron SVG 箭头按钮 + 深色"完成"按钮
   - 已移除：日记、户外时间、电视条目
   - 已新增：待办条目
   - 排序箭头改为圆形 chevron SVG（`chevron-up` / `chevron-down`）
   - `moveUp(idx)` / `moveDown(idx)` 交换位置
   - `saveOrder()` → `PUT /api/settings/panel-order`

**生命周期：** `onMounted` → `loadOrder()` → `GET /api/settings/panel-order` → 恢复排序

**补卡按钮：** 仅对 `discipline` 和 `nosugar` 显示，包含日历图标的 `backfill-btn`，点击时 `$emit('backfill', item)`，然后 `App.vue` 弹出 `TimePickerModal`（AddEntryPanel 保持打开）。

---

## ThoughtFormPanel.vue

**角色：** 随记创建表单。标题固定为"随记"，提交时根据 `kind` 分别标记"随记·太阳"（正向）或"随记·阴雨"（负向）。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |
| `kind` | String | `'positive'` / `'negative'` |

| Emit | 参数 | 说明 |
|------|------|------|
| `cancel` | — | 取消 |
| `create` | `{kind, time, note}` | 提交 |

**使用 `useTimePicker`：** 所有时间选择器状态和方法通过 composable 注入。

**表单字段：**
- 记录时间（可点击弹出 `TimePickerModal`）
- 备注（textarea，最大 200 字，实时计数）

**重置：** `watch(visible)` → 清空表单 + 重置时间为当前。

---

## AssetFormPanel.vue

**角色：** 资产余额记录表单。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |

| Emit | 参数 | 说明 |
|------|------|------|
| `cancel` | — | 取消 |
| `create` | `{time, amount}` | 提交 |

**表单字段：** 记录时间 + 资产余额（number input）

---

## UricFormPanel.vue

**角色：** 尿酸值记录表单。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |

| Emit | 参数 | 说明 |
|------|------|------|
| `cancel` | — | 取消 |
| `create` | `{time, amount}` | 提交 |

---

## ExerciseFormPanel.vue

**角色：** 运动记录表单。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |

| Emit | 参数 | 说明 |
|------|------|------|
| `cancel` | — | 取消 |
| `create` | `{time, sport, amount}` | 提交 |

**运动类型：** 跑步、骑行、俯卧撑、跳绳、游泳、徒步、自由活动

**计算属性：**
- `sportLabel` — 动态标签（如"跑步距离"）
- `sportUnit` — 动态单位（如"千米"）

**重置：** `watch(visible)` → 清空 + 重置运动类型为"跑步"。

---

## ReadingFormPanel.vue

**角色：** 读书记录表单，支持 DeepSeek 书籍信息自动查询。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |

| Emit | 参数 | 说明 |
|------|------|------|
| `cancel` | — | 取消 |
| `create` | `{time, title, author, cover, isbn, category, status, rating, note}` | 提交 |

**状态胶囊：** 在读 / 已读，通过胶囊按钮切换。

**7 种分类胶囊：** 文学、历史、哲学、科幻、技术、心理学、其他，各有独立 hex 色值。

**DeepSeek 查询：** 输入书名或 ISBN → 点击查询按钮 → `POST /api/book-info` → 自动填充标题、作者、封面图、出版社等信息。

**评分：** 1-10 分可点击星星组件（2 分 = 1 星，奇数 = 半星）。

---

## MovieFormPanel.vue

**角色：** 影视记录表单，支持 DeepSeek 影视信息自动查询。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |

| Emit | 参数 | 说明 |
|------|------|------|
| `cancel` | — | 取消 |
| `create` | `{time, title, director, year, cover, type, status, rating, note}` | 提交 |

**类型胶囊：** 电影（`#CB99B0`）、剧集（`#D4A882`）、动漫（`#A099C4`），通过胶囊按钮切换。

**观看状态胶囊：** 想看 → 在看 → 已看 → 弃剧，只有"已看"状态启用评分功能。

**DeepSeek 查询：** 输入影视名称 → 点击查询按钮 → `POST /api/media-info` → 自动填充标题、导演、年份、封面图等信息。

**评分：** 1-10 分可点击星星组件（同读书评分系统）。

---

## TodoFormPanel.vue

**角色：** 待办事项创建表单。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |

| Emit | 参数 | 说明 |
|------|------|------|
| `cancel` | — | 取消 |
| `create` | `{time, content, dueDate}` | 提交 |

**时间选择器：** 半小时间隔（00/30 分）的时间选择器。

**文本限制：** 待办内容最大 50 字符，实时计数显示。

**到期日：** 可选设置截止日期，用于 TodoPage 的自动分类（逾期/今天/近期/稍后）。

---

## TimePickerModal.vue

**角色：** 底部弹出（bottom sheet）5 列时间选择器（年 / 月 / 日 / 时 / 分），展开自底部，带取消/确认按钮。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |
| `pickYear/Month/Day/Hour/Minute` | Number | 当前选中值 |

| Emit | 参数 | 说明 |
|------|------|------|
| `close` | — | 取消选择（关闭面板） |
| `confirm` | — | 确认时间 |
| `adjustYear/Month/Day/Hour/Minute` | delta | `+1` 或 `-1` |

**视觉设计：**
- Bottom sheet 布局：`border-radius: 20px 20px 0 0`，`max-width: var(--content-width)`
- 每列：SVG chevron 上下箭头按钮 + 大号数字显示 + 单位标签
- 头部：左侧"取消" + 中间"选择时间" + 右侧"确定"
- 动画：`translateY(100%) → 0` 配合 opacity 过渡

**内部方法：** `pad(n)` — 补零显示（月/日/时/分）。

---

## ConfirmModal.vue

**角色：** 确认对话框（替代原生 `alert()`）。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |
| `message` | String | 提示文本 |

| Emit | 参数 | 说明 |
|------|------|------|
| `close` | — | 点击"确定" |

**视觉：** 绿色勾号图标 + 消息文本 + 深色"确定"按钮。

---

## StatisticsPage.vue

**角色：** 全页统计数据展示。

| Emit | 参数 | 说明 |
|------|------|------|
| `back` | — | 返回时间轴 |

**数据获取：**
- `onMounted` → `fetchAll()` → 分页拉取所有条目
- `categoryData` (computed) → 为 8 种分类分别构建热力图

**8 种分类：** 随记·太阳、随记·阴雨、自律、禁糖、运动、资产、读书、影视

**热力图构建：** `buildHeatmapWeeks(entries, year)`
- 输入：某类型的条目列表 + 年份
- 输出：7 列周数组，每格包含 `{date, filled, isToday, label, dateNum}`
- 每月 1 号显示月份标签，每月末显示日期号

**统计卡片：** `topStats` (computed)
- 总记录数 / 本月记录数 / 记录天数 / 已使用天数

**交互：** 热力图区域支持 pointer 拖拽水平滚动。CSS `touch-action: pan-y` 确保垂直滑动穿透到页面滚动，仅水平滑动触发热力图滚动。

---

## TodoPage.vue

**角色：** 待办事项全页管理。

| Emit | 参数 | 说明 |
|------|------|------|
| `back` | — | 返回时间轴 |
| `todo-synced` | entry | 完成待办时同步到时间轴 |

**4 种自动分类：**
- **逾期 overdue** — 截止日期已过，红色梯度（`#D4787A` → `#E8A0A2`）
- **今天 today** — 截止日期为今天，橙色梯度
- **近期 upcoming** — 未来 7 天内，蓝色梯度
- **稍后 later** — 7 天以上或无截止日期，灰色梯度

**交互：**
- 添加：通过 `TodoFormPanel`（半小时间隔时间选择器，50 字符限制）
- 切换完成：点击左侧圆圈 → `PATCH /api/todos/:id` + 同步到时间轴
- 滑动操作：左滑显示编辑/删除操作按钮
- 编辑：双击文本行内编辑（`blur` / `enter` 保存）

**数据获取：** `<script setup>` 顶层立即执行 `fetchTodos()`。

---

## 共享 Composable

### `useTimePicker.js`

**角色：** 时间选择器状态与逻辑复用。

**导出：**
- 响应式状态：`currentTime`, `showTimeModal`, `pickYear`, `pickMonth`, `pickDay`, `pickHour`, `pickMinute`
- 计算属性：`displayTime`, `daysInMonth`
- 方法：`openTimePicker`, `closeTimePicker`, `confirmTime`, `adjustMonth`, `adjustDay`, `adjustHour`, `adjustMinute`
- 工具：`pad`

**使用者：** `ThoughtFormPanel`, `AssetFormPanel`, `UricFormPanel`, `ExerciseFormPanel`, `App.vue`（补卡）

### `date.js`

**导出函数：** `formatRecordTime(date: Date) → { recordedAt: string, monthKey: string }`

- `recordedAt`：`"YYYY-MM-DD HH:MM:00"`
- `monthKey`：`"YYYY-MM"`
