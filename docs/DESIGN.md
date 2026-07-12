# 时间轴 · 设计系统

> 版本 1.1 · "Ink on Paper" 主题

---

## 设计理念

核心隐喻：**私人日记本**。屏幕即纸张，条目即墨迹，分类色彩区分不同活动类型。

---

## 配色系统

### 语义 Token

| Token | Hex | 用途 |
|-------|-----|------|
| `--color-paper` | `#FCFAF7` | 页面背景（暖纸白） |
| `--color-card` | `#FFFFFF` | 卡片/面板底色（纯白） |
| `--color-ink` | `#1A1A1A` | 主文字（墨黑） |
| `--color-graphite` | `#6E6E6E` | 次级文字 / 辅助信息 |
| `--color-pencil` | `#E0DCD6` | 边框 / 分割线 / 时间轴线 |
| `--color-border-light` | `#F0EDEA` | 浅色边框 |
| `--color-surface-dim` | `#F7F5F2` | 卡片灰色底板 |

### 分类色标

| 分类 | 色值 | 对应条目 |
|------|------|---------|
| 绿色系 Green | `#7BA88A` | 正向念头、自律、禁糖、运动 |
| 琥珀色 Amber | `#D4A574` | 资产记录、黄色分类 |
| 红色系 Rose | `#D4787A` | 负向念头、红色分类 |

### 条目类型专用色

| 条目 | 色值 |
|------|------|
| 念头 | `#9DB5C9` |
| 资产记录 | `#D4B87A` |
| 运动 | `#84B8A4` |
| 读书 | `#A099C4` |
| 电影 | `#CB99B0` |
| 剧集 | `#D4A882` |
| 日记 | `#8CAD8C` |
| 自律 | `#8CA4BD` |
| 尿酸 | `#BE999B` |
| 禁糖 | `#C88C8C` |
| 户外时间 | `#92B4C7` |

---

## 排版

### 字体栈

```css
font-family: -apple-system, BlinkMacSystemFont, "Segoe UI",
  "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei",
  "Helvetica Neue", Arial, sans-serif;
```

设计决策：中英文混排应用，系统字体是最佳选择。PingFang SC 和 Microsoft YaHei 提供优秀的中文阅读体验。

### 字号层级

| 用途 | 手机 | 平板(≥600px) | 大屏(≥768px) |
|------|------|-------------|-------------|
| 年月标题 | 22px / 700 | 24px / 700 | 24px / 700 |
| 日期标签 | 16px / 700 | 18px / 700 | 20px / 700 |
| 卡片标题 | 14px / 600 | 15px / 600 | 15px / 600 |
| 卡片正文 | 13px | 14px | 14px |
| 时间戳 | 12px | 13px | 13px |

---

## 圆角

| 用途 | 值 | CSS 变量 |
|------|----|---------|
| 小元素（热力图格、色点） | 3px-6px | `--radius-sm: 6px` |
| 卡片 / 输入框 / 按钮 | 10px | `--radius-md: 10px` |
| 模态框 / 面板 | 16px | `--radius-lg: 16px` |

特殊：时间轴卡片左侧无边角（`border-radius: 0 var(--radius-md) var(--radius-md) 0`），与时间轴线对齐。

---

## 评分星星

所有记录类型（读书、影视）共用统一的星级评分组件。

| 属性 | 值 |
|------|-----|
| 完整星 SVG 颜色 | `#D4A574`（琥珀金） |
| 空心星 SVG 颜色 | `var(--color-pencil)` |
| 计分规则 | 2 分 = 1 颗星星 |
| 半星规则 | 奇数分显示半星（如 7 分 = ★★★½） |
| 最大星星数 | 5 颗（满分 10 分） |
| 交互 | 点击星星左半 = 整数分，右半 = .5 分 |

### 星星 SVG 组件结构

```html
<svg viewBox="0 0 20 20" class="rating-star">
  <defs>
    <clipPath id="half-clip">
      <rect x="0" y="0" width="10" height="20" />
    </clipPath>
  </defs>
  <!-- 空心底 -->
  <path d="M10 1.5l2.4 4.9 5.4.8-3.9 3.8.9 5.4-4.8-2.5-4.8 2.5.9-5.4-3.9-3.8 5.4-.8z"
        fill="none" stroke="var(--color-pencil)" />
  <!-- 填充层（clipPath 控制半星） -->
</svg>
```

---

## 分类胶囊颜色

### 读书分类（7 种）

| 分类 | Hex | 文字色 |
|------|-----|--------|
| 文学 | `#B8A9C9` | `#4A3B5C` |
| 历史 | `#C9B8A9` | `#5C4A3B` |
| 哲学 | `#A9B8C9` | `#3B4A5C` |
| 科幻 | `#A9C9C0` | `#3B5C4A` |
| 技术 | `#A9B0C9` | `#3B3B5C` |
| 心理学 | `#C9B8B8` | `#5C4A4A` |
| 其他 | `#C0C0C0` | `#4A4A4A` |

### 影视类型胶囊

| 类型 | Hex | 说明 |
|------|-----|------|
| 电影 | `#CB99B0` | 玫瑰色 |
| 剧集 | `#D4A882` | 琥珀色 |
| 动漫 | `#A099C4` | 紫色 |

### 编辑面板标签胶囊

v1.1 新增：EditEntryForm 中 tag 胶囊交互样式。

| 属性 | 值 |
|------|-----|
| 胶囊底色 | 按标签名 hash 分配分类色标颜色（绿/琥珀/玫瑰） |
| 文字颜色 | 对应颜色的深色变体（如 `#4A3B5C`） |
| × 关闭按钮 | 16px × 16px 圆形，`border-radius: 50%`，hover 时背景加深 |
| 关闭图标 | 细线 SVG ×（`stroke-width: 1.5`） |
| 添加按钮 | 浅灰圆角按钮 `+ ` 文字，点击弹出输入框 |
| 胶囊间距 | `gap: 6px` flex-wrap 布局 |
| 圆角 | `--radius-sm: 6px` |
| 字号 | 12px / 500 |

---

## 待办分类红色梯度

TodoPage 中待办事项按紧急程度自动分类，颜色由深至浅表示紧迫度递减：

| 分类 | 颜色 | 条件 |
|------|------|------|
| 逾期 (overdue) | `#D4787A` → `#E8A0A2` | `dueDate < today` |
| 今天 (today) | `#E8A0A2` → `#F0C0C0` | `dueDate === today` |
| 近期 (upcoming) | `#A9B8C9` → `#C8D4E0` | `dueDate` 在未来 7 天内 |
| 稍后 (later) | `#C0C0C0` → `#E0E0E0` | `dueDate` 在 7 天后或无截止日期 |

渐变色用于分类标题栏背景，卡片左侧色标使用对应分类的首色。

---

## 响应式断点

```css
/* 基准（手机） */
:root { --content-width: 480px; }

/* 大屏手机 / iPad mini 竖屏 */
@media (min-width: 600px) { :root { --content-width: 544px; } }

/* iPad mini 横屏 / 小平板 */
@media (min-width: 768px) { :root { --content-width: 640px; } }

/* 桌面 */
@media (min-width: 960px) { :root { --content-width: 720px; } }
```

所有固定宽度元素（顶栏、底栏、面板、模态框）统一引用 `var(--content-width)`。

---

## 间距系统

| 断点 | 水平内边距 | 时间轴顶部 | 时间轴底部 |
|------|----------|----------|----------|
| 手机 | 24px | 64px | 60px |
| ≥600px | 28px | 72px | 64px |
| ≥768px | 36px | 80px | 68px |

网格列宽（时间轴布局）：
- 手机：`46px 4px 20px 9px 1fr`
- ≥600px：`52px 6px 24px 12px 1fr`
- ≥768px：`56px 6px 28px 12px 1fr`

---

## 动画与过渡

| 元素 | 效果 | 时长 | 缓动 |
|------|------|------|------|
| 底部面板弹入 | `translateY(100%) → 0` + opacity | 300ms | `cubic-bezier(0.32,0.72,0,1)` |
| 表单弹窗 | 同上 | 280ms | 同上 |
| 模态框弹出 | `scale(0.92) → 1` + opacity | 220ms | 同上 |
| 下拉菜单 | opacity fade | 200ms | ease |
| 资产金额 | 彩虹色流动（渐变位移） | 3s 循环 | linear |
| 卡片 hover | scale 微缩 | 120ms-150ms | ease |

### 彩虹动画

```css
@keyframes rainbow-flow {
  to { background-position: 200% center; }
}
```

仅用于资产余额数字，使用渐变色系（匹配项目调色板）。

---

## 毛玻璃效果

底部栏和顶栏使用：

```css
background: #fcfcfc;
backdrop-filter: blur(16px);
-webkit-backdrop-filter: blur(16px);
```

---

## 可访问性

- `prefers-reduced-motion` 全局媒体查询，关闭所有动画
- `aria-label` 标注所有功能性按钮（排序、关闭、添加）
- 可见键盘焦点：保留浏览器默认 outline
- 触摸目标最小 32×32px（Apple HIG 标准）

---

## 时间选择器日历

v1.1 重新设计的 TimePickerModal，采用日历 + 滚动轮双区布局。

### 上半区 — 日历面板

| 属性 | 值 |
|------|-----|
| 月导航箭头 | SVG chevron 细线图标（`stroke-width: 1.5`），颜色 `var(--color-graphite)` |
| 月份文字 | 16px / 600，居中，`var(--color-ink)` |
| 星期头行 | 一 二 三 四 五 六 日，11px / 400，`var(--color-graphite)`，底部细线分割 |
| 日期格子 | 等宽 7 列 grid，格子 `aspect-ratio: 1`，14px / 400 |
| 当天日期 | 蓝色 accent（`#007AFF`）描边圆底，白字 |
| 选中日期 | 蓝色 accent 实心圆底，白字，`border-radius: 50%` |
| 默认日期 | `var(--color-ink)` 无背景 |
| 非当月日期 | `var(--color-pencil)` 淡色 |
| 农历标注 | 10px，`var(--color-pencil)`，仅展示节气或初一/十五 |
| 日历高度 | 固定约 280px，overflow hidden |

### 下半区 — 滚动轮

| 属性 | 值 |
|------|-----|
| 列数 | 2 列（时 / 分），等宽 |
| 滚动容器 | `scroll-snap-type: y mandatory`，`overflow-y: auto` |
| 每项高度 | 40px，上下各留 80px padding 形成可滚动区域 |
| 吸附行为 | `scroll-snap-align: center` |
| 字体 | 时/分数字：20px / 500，选中项加大到 22px / 600 |
| 非选中项 | `var(--color-pencil)` 淡色，选中项 `var(--color-ink)` |
| 分隔冒号 | 列间固定 `:` 字符，20px，`var(--color-graphite)` |

### 高亮区

| 属性 | 值 |
|------|-----|
| 位置 | 时/分列垂直居中，覆盖 3 行高度（约 120px） |
| 背景色 | `rgba(0, 122, 255, 0.08)` 蓝色半透明 |
| 上下边界 | 1px 细线 `rgba(0, 122, 255, 0.15)` |
| 圆角 | `--radius-sm: 6px` |
| 层级 | `z-index: 1`，数字文字 `z-index: 2` 浮于上方 |

---

## 筛选按钮

v1.1 新增：YearMonthHeader 中的漏斗筛选按钮及弹出面板。

### 漏斗图标按钮

| 属性 | 值 |
|------|-----|
| 尺寸 | 32px × 32px 圆形按钮 |
| 图标 | SVG funnel 图标，`stroke-width: 1.5` |
| 默认颜色 | `var(--color-graphite)` 灰色 |
| 激活颜色 | `#007AFF` 蓝色（有筛选条件时） |
| Hover | 背景 `var(--color-surface-dim)` 圆形阴影 |
| 位置 | YearMonthHeader 右侧，搜索图标旁 |

### 筛选面板

| 属性 | 值 |
|------|-----|
| 布局 | 多选项列表，每行：色点(8px 圆) + 标签名 + checkbox |
| 面板样式 | 白色卡片，`border-radius: var(--radius-lg)`，`box-shadow: 0 4px 20px rgba(0,0,0,0.08)` |
| 间距 | 选项行 `padding: 10px 16px`，底部操作栏 `padding: 12px 16px` |
| 底部按钮 | "全部显示" / "全部隐藏" 文字按钮，`var(--color-graphite)` |
| 动画 | scale(0.95)→1 + opacity fade，200ms |
| 定位 | Teleport to body，绝对定位在漏斗按钮下方 |

---

## 文件组织

| 文件 | 作用 |
|------|------|
| `src/style.css` | CSS 变量定义、全局 reset、响应式断点、共享过渡动画 |
| 各组件 `<style scoped>` | 组件专属样式，引用 CSS 变量 |
