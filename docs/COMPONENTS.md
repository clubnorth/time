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

**角色：** 固定在顶部的年月选择器。v1.1 新增筛选功能。

| Prop | 类型 | 说明 |
|------|------|------|
| `modelValue` | String | 当前月份 `YYYY-MM` |
| `months` | Array | 可选月份列表 `[{label, value}]` |
| `activeFilters` | Set | 当前激活的条目类型筛选集合 |

| Emit | 参数 | 说明 |
|------|------|------|
| `update:modelValue` | value | 选中月份变化 |
| `update:activeFilters` | filters | 筛选条件变化 |
| `toggle-stats` | — | 打开统计页 |
| `toggle-todo` | — | 打开待办页 |

**状态：** `showDropdown` — 控制下拉菜单显隐；`showFilterPanel` — 控制筛选面板显隐。

**v1.1 新增筛选功能：**
- 搜索图标旁增加漏斗按钮（SVG funnel icon，未激活时灰色，有筛选时高亮蓝色）
- 点击弹出多选筛选面板，列出所有条目类型（随记·太阳/阴雨、自律、禁糖、运动、资产、读书、影视、待办）
- 每项带 checkbox + 色点标识，选中即显示
- 底部"全部显示" / "全部隐藏"快捷按钮
- 筛选面板使用 Teleport 挂载到 body

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

## EditEntryForm.vue

**角色：** 条目编辑面板，v1.1 新增。长按卡片弹出，按条目类型展示专属编辑字段。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |
| `entry` | Object | 待编辑的条目数据 |

| Emit | 参数 | 说明 |
|------|------|------|
| `close` | — | 关闭编辑面板 |
| `saved` | entry | 保存成功，通知父组件刷新 |

**按类型分派字段：** 随记（备注）、自律/禁糖（时间）、运动（类型+数量+时间→自动重算卡路里）、读书（书名/作者/封面/分类胶囊/状态/评分）、影视（标题/导演/年份/封面/类型胶囊/状态/评分）、资产（余额）、尿酸（数值）。

**标签胶囊：** 支持 `parseTags` 解析逗号分隔标签，已选标签以彩色胶囊展示（× 移除），+ 按钮添加新标签。

**评分星星：** 复用与 ReadingFormPanel 相同的 1-10 分可点击星星组件。

**保存：** `PUT /api/entries/:id` → emit `saved` → App.vue 刷新 timelineData。

---

## TodoTimePicker.vue

**角色：** 待办专用时间选择器，v1.1 新增。10 分钟步进，滚动吸附。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |
| `modelValue` | String | 当前选中时间 `HH:MM` |

| Emit | 参数 | 说明 |
|------|------|------|
| `update:modelValue` | value | 时间变化 |
| `confirm` | value | 确认选择 |

**特性：**
- 小时/分钟双列滚动轮（`scroll-snap-type: y mandatory`）
- 分钟列以 10 分钟为间隔（00, 10, 20, 30, 40, 50）
- 当前值高亮区（蓝色背景条，上下居中）
- 底部"取消"/"确定"按钮

---

## TimePickerModal.vue

**角色：** 底部弹出（bottom sheet）时间选择器，v1.1 重新设计。上半部分日历面板选择日期，下半部分滚动轮选择时分。

| Prop | 类型 | 说明 |
|------|------|------|
| `visible` | Boolean | 控制显隐 |
| `modelValue` | String | 当前日期时间值 ISO 格式 |

| Emit | 参数 | 说明 |
|------|------|------|
| `close` | — | 取消选择（关闭面板） |
| `confirm` | value | 确认时间（ISO 字符串） |

**布局（v1.1 新设计）：**
- **上半区 — 日历面板**：
  - 月导航栏：左箭头 SVG + 月份文字（居中对齐）+ 右箭头 SVG
  - 星期头行：一 二 三 四 五 六 日（灰色小字）
  - 日期网格：7 列 × 5-6 行，当天日期蓝色 accent 圆底白字，选中日期蓝色实心圆底
  - 农历/节气标注（小字灰色，仅展示当月 1-2 个关键日期）
- **下半区 — 滚动轮**：
  - 时 / 分 双列滚动轮（`scroll-snap-type: y mandatory`）
  - 中间高亮区（`rgba(0,122,255,0.08)` 蓝色半透明背景条）
  - 分钟 5 分钟间隔（00, 05, 10, ..., 55）
  - 滚动吸附到最近项
- **底部栏**：左侧"取消" + 右侧蓝色"确定"按钮
- 动画：`translateY(100%) → 0` 配合 opacity 过渡

**内部方法：** `daysInMonth(year, month)`, `buildCalendar(year, month)` 构建当月日历网格。

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

**角色：** 全页统计数据展示。v1.1 改为原生滚动，新增读书/影视热力图。

| Emit | 参数 | 说明 |
|------|------|------|
| `back` | — | 返回时间轴 |

**v1.1 滚动改造：**
- 移除自定义 pointer 拖拽滚动，改用原生 `overflow-x: auto` + `-webkit-overflow-scrolling: touch`
- `scroll-behavior: smooth` 提供平滑滚动体验
- 热力图行在移动端原生 momentum 滚动

**数据获取：**
- `onMounted` → `fetchAll()` → 分页拉取所有条目
- `categoryData` (computed) → 为 8 种分类分别构建热力图

**8 种分类：** 随记·太阳、随记·阴雨、自律、禁糖、运动、资产、读书、影视

**v1.1 新增读书/影视热力图：**
- 读书热力图：按阅读状态（在读/已读）区分格子颜色深度
- 影视热力图：按观看状态（想看/在看/已看/弃剧）区分格子颜色深度
- 深色 = 已读/已看，浅色 = 在读/在看，最浅 = 想看

**热力图构建：** `buildHeatmapWeeks(entries, year)`
- 输入：某类型的条目列表 + 年份
- 输出：7 列周数组，每格包含 `{date, filled, isToday, label, dateNum, intensity}`
- 每月 1 号显示月份标签，每月末显示日期号
- `intensity` 字段用于读书/影视的渐变颜色映射

**统计卡片：** `topStats` (computed)
- 总记录数 / 本月记录数 / 记录天数 / 已使用天数

**交互：** 热力图区域使用原生横向滚动（`overflow-x: auto`），CSS `touch-action: pan-y` 确保垂直滑动穿透。

---

## TodoPage.vue

**角色：** 待办事项全页管理。v1.1 增强左滑操作和完成同步。

| Emit | 参数 | 说明 |
|------|------|------|
| `back` | — | 返回时间轴 |
| `todo-synced` | entry | 完成待办时同步到时间轴 |

**4 种自动分类（v1.1 优化）：**
- **逾期 overdue** — 截止日期已过，红色梯度（`#D4787A` → `#E8A0A2`）
- **今天 today** — 截止日期为今天，橙色梯度
- **近期 upcoming** — 未来 7 天内，蓝色梯度
- **稍后 later** — 7 天以上或无截止日期，灰色梯度
- 分类切换使用标签栏（TabBar），流畅 CSS 过渡动画

**交互：**
- 添加：通过 `TodoFormPanel`（半小时间隔时间选择器，50 字符限制）
- 切换完成：点击左侧圆圈 → `PATCH /api/todos/:id` + 同步到时间轴
- **左滑操作（v1.1）：** 使用 `@touchstart` / `@touchmove` / `@touchend` 实现流畅左滑，reveal 编辑（蓝色）和删除（红色）操作按钮，需滑动超过阈值（60px）才固定展开，点击空白区域收回
- 编辑：双击文本行内编辑（`blur` / `enter` 保存）
- **完成同步（v1.1）：** 勾选待办完成时自动在时间轴中创建对应条目（`POST /api/entries`），emit `todo-synced` 通知 App.vue 刷新时间轴

**数据获取：** `<script setup>` 顶层立即执行 `fetchTodos()`。

**列表优化（v1.1）：** 使用 `v-memo` 优化长列表渲染性能，仅当 todo 状态或内容变化时重新渲染。

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
