# 时间轴 v1.1.0 · 开发对话完整记录

> 基于 v1.0.0，本次对话新增功能 15 项，修复 Bug 30+ 项

---

## 一、新增功能

### 1. 待办事项系统重构
- 四分类自动归类：你忘记做了 / 今日 / 未来三日 / 以后
- 待办卡片左滑操作（修改 / 删除）
- 完成联动时间轴：勾选完成 → 自动写入"完成待办"记录
- 待办添加入口：添加记录面板新增"待办事项"卡片
- 新建后自动跳转待办页面

### 2. 时间选择器重设计（TimePickerModal）
- 底部弹出式日历 + 时分双列滚轮
- 月份左右箭头切换，7 列网格日期布局
- 蓝色主色调，选中日期蓝色圆底白色数字
- 时分渐变遮罩效果，中心高亮蓝色

### 3. 待办专用时间选择器（TodoTimePicker）
- 分钟仅允许 00/20/40 三档（后改为 00/10/20/30/40/50 六档）
- 点击数字平滑滚动到蓝色高亮行
- 日期按钮回写父组件事件修复

### 4. DeepSeek AI 优化
- 读书记录：自动查询作者、国籍、字数、出版年份、3 个标签
- 影视记录：自动查询导演、主演、国家、上映年份、时长/集数、3 个标签
- 提示词多次迭代优化，禁止返回"未知"
- API_BASE 路径修复

### 5. 评分系统
- 读书、影视添加 1-10 分星星评分
- 2 分 = 1 颗金星，奇数为半星
- 卡片显示：★★★★★ 7分

### 6. 图书分类
- 推理悬疑 / 明日方舟 / 人物历史 / 社会科学 / 小说文学 / 自我成长 / 杂项拾遗
- 7 种彩色胶囊分类选择

### 7. 时间轴卡片编辑功能
- 长按 3 秒弹出气泡：修改 / 删除 / 取消
- 各类型独立编辑字段：
  - 资产：余额数字
  - 尿酸：尿酸值
  - 运动：运动类型 + 数值 + 卡路里重算
  - 随记：太阳/阴雨 + 备注
  - 自律/禁糖：连续天数
  - 读书/影视：名称、导演、主演、国家、年份、时长、标签胶囊、评分

### 8. 标签胶囊编辑器
- 每个标签单独胶囊展示，右上角 × 移除
- + 按钮弹出输入框，最多 4 汉字
- 读书/影视回显标签单独解析（parseTags 函数）
- 国家字段防括号叠套

### 9. 读书分类下拉框
- 编辑页分类从文本输入改为 select 下拉

### 10. 打卡/补卡规则重构
- 规则①：普通打卡仅校验今日，已有记录弹窗禁止
- 规则②：补卡支持手动选日期，校验重复 + 未来日期拦截
- 删除打卡后自动 recalculate 重算连续天数
- 删除中间某天，次日从第 1 天重新计数

### 11. 自律/禁糖卡片长按弹窗
- 仅保留"删除"+"取消"，移除"修改"

### 12. 热力图滑动优化
- 移除自定义 JS pointer 事件处理器
- 改用浏览器原生 `overflow-x: auto` + `-webkit-overflow-scrolling: touch`
- 手机平板滑动手感顺滑无延迟

### 13. 时间轴筛选功能
- 顶栏右侧漏斗按钮
- 多选 9 种记录类型，彩色圆点 + 蓝色对号
- 确定后时间轴只渲染已选类型

### 14. 年份回显修复
- 解析时自动 `replace(/年$/, '')` 去掉尾部"年"
- 保存时统一加"年"，防止"2026年年"叠套

### 15. 运动卡路里公式更新
- 基于 69kg 体重、心率 160，标准 MET 公式
- 编辑运动条目保留卡路里重算

---

## 二、Bug 修复记录

| # | 现象 | 根因 | 修复 |
|---|------|------|------|
| B1 | 点击日历日期，父组件时间不更新 | TodoTimePicker `@click="day = d"` 未 emit `adjustDay` | 改为 `selectDay(d)` 调用 |
| B2 | 编辑待办保存后时间重置为当天 | `closeTimePicker()` 只关弹窗未同步 `currentTime` | close 前 `currentTime = new Date(pick*)` |
| B3 | 禁止糖分重复打卡不拦截 | `limit=500` 被服务端截为30，打卡记录在第31+条漏检 | `fetchAllEntries()` 全程扫描 |
| B4 | `toISOString()` 导致 UTC 时区日期偏移 | 中文时区 UTC+8，凌晨查日期差一天 | 改用 `getFullYear/Month/Date` 本地日期 |
| B5 | 标签保存后重开全部合并为一个胶囊 | `parseField` 去 HTML 后无分隔符无法拆分 | 新增 `parseTags()` 用正则精确提取 |
| B6 | 国家字段 `【【美国】】` 括号叠套 | 已有括号时保存再加括号 | 读入时 `replace(/[【】]/g, '')` 去括号 |
| B7 | 运动/资产编辑页数据显示"你今天跑步跑了..." | 传入了卡片渲染后的描述而非原始值 | `handleEditEntry` 改用 `full.description` |
| B8 | 待办新建时间使用过期默认值 | watch 只设 `currentTime` 未同步 pickState | 同步 pickYear/Month/Day/Hour/Minute |
| B9 | `RecalculateEntries` 连续天数可能计算错误 | `d1.Sub(d2).Hours() == 24` 浮点比较 | 改用日期字符串 + 23-25 小时范围 |
| B10 | 打卡/补卡 API 失败静默吞错 | catch 只 console.log | 弹窗"打卡失败，请重试" |
| B11 | 待办删除无二次确认 | 直接调用 API | `confirm()` 弹窗确认 |
| B12 | 编辑运动条目卡路里丢失 | handleSave 未写入 valence | 添加 MET 公式重算 |
| B13 | 读书/影视 DeepSeek 调用路径硬编码 | 缺少 `${API_BASE}` 前缀 | 改为相对路径引用 |
| B14 | 长按计时器组件销毁后触发 | 未在 `onUnmounted` 清理 | 添加 clearTimeout |
| B15 | 时分滚轮共用一个 debounce timer | 同时滚动互相取消 | 拆分为 hour/minScrollTimer |
| B16 | UTC 日期解析偏移一天 | `new Date(date)` 按 UTC 解析 | `new Date(date + 'T00:00:00')` |
| B17 | 统计页换年份热力图滚动位不重置 | 无 watch 处理 | watch(year) → scrollLeft |
| B18 | json.Marshal 错误被忽略 | `jsonBody, _ := json.Marshal(body)` | 添加 err 判断 |
| B19 | 死代码 getConsecutiveDay | 重构后不再调用 | 删除 |
| B20 | 已完成待办仍可修改 | 无 completed 判断 | `v-if="!todo.completed"` |
| B21 | 分钟滚轮吸附错位 | scrollTop 多了 80px | 统一 `index * 40` |
| B22 | 筛选对号不即时显示 | 用 `activeFilters`（props）判断 | 改用 `tempFilters` |
| B23 | `pointerleave` 导致鼠标长按失效 | 鼠标微动触发取消 | 移除 pointerleave |
| B24 | 补卡面板提前关闭 | handleBackfill 立即 `showAddPanel=false` | 移至确认成功后关闭 |
| B25 | 影视记录描述"再看"错别字 | 应为"在看" | 全部改为"在看" |
| B26 | 国家字段多余括号叠套 | 同上 B6 |
| B27 | 编辑运动卡路里公式不准确 | 同上 B12 |
| B28 | 统计页无读书/影视热力图 | categories 缺少两个类型 | 添加 reading + movie |
| B29 | 念头图标排序模式消失 | `v-if` 带 `!sortMode` 条件 | 移除 sortMode 限制 |
| B30 | 排序箭头丑化 | 方形框 + ↑↓ 文字 | 圆形 + SVG chevron |
| B31 | 补卡按钮消失 | AddEntryPanel 未渲染 backfill emit | 重新添加 action-group |
| B32 | 已打钩待办仍可修改 | 同上 B20 |
| B33 | 书籍添加缺在读验证 | 已读需存在对应在读 | 扫描全量在读记录 |
| B34 | 待办添加未自动跳转待办页 | create 后未设 showTodo=true | 添加 showTodo=true |

---

## 三、代码审查专项修复

| # | 级别 | 问题 | 修复 |
|---|------|------|------|
| C1 | 严重 | TodoFormPanel picker 不同步 | 同步 pickYear/Month/Day/Hour/Minute |
| C2 | 严重 | RecalculateEntries 日期比较浮点 | 改用日期字符串 23-25h 范围 |
| C3 | 严重 | 读书/影视 API 路径硬编码 | 加 API_BASE |
| M1 | 中 | 长按 timer 未清理 | onUnmounted clearTimeout |
| M2 | 中 | 滚轮共用 debounce | 拆分独立 timer |
| M4 | 中 | UTC 偏移 | T00:00:00 |
| M6 | 中 | recalculate 硬编码 | 加 API_BASE |
| M7 | 中 | API 失败静默 | 弹窗提示 |
| M8 | 中 | 删除无确认 | confirm() |
| M9 | 中 | 编辑运动丢卡路里 | MET 重算 |
| L2 | 低 | getConsecutiveDay 死代码 | 删除 |
| L4 | 低 | json.Marshal 未处理 | 加 err 判断 |
| L5 | 低 | 热力图年份切换 | watch scroll |
