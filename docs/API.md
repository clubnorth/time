# 时间轴 · API 参考

> 基础地址：`http://localhost:8080`

---

## 通用约定

### 响应格式

所有 API 统一返回 JSON：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | `0` = 成功，`1` = 错误 |
| `message` | string | 状态描述 |
| `data` | any | 响应数据（错误时为 `null`） |

### 时间格式

所有时间字段使用 `YYYY-MM-DD HH:MM:SS` 格式，本地时间。

---

## 条目 API

### 获取条目列表（分页）

```
GET /api/entries?limit=30&before=2026-07-10%2023%3A59%3A00
```

**参数：**

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `limit` | int | 否 | 返回条数（默认 30，最大 100） |
| `before` | string | 否 | 游标分页，返回 `recorded_at < before` 的条目 |

**响应：**

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "type": "thought",
      "title": "念头·正",
      "description": "今天天气很好",
      "category": "green",
      "valence": "positive",
      "recorded_at": "2026-07-10 14:30:00",
      "created_at": "2026-07-10 14:30:01",
      "updated_at": "2026-07-10 14:30:01"
    }
  ]
}
```

### 获取条目列表（按月）

```
GET /api/entries?month=2026-07
```

**参数：**

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `month` | string | 是 | 格式 `YYYY-MM` |

### 创建条目

```
POST /api/entries
Content-Type: application/json
```

**请求体：**

```json
{
  "type": "thought",
  "title": "念头·正",
  "description": "今天天气很好",
  "category": "green",
  "valence": "positive",
  "recorded_at": "2026-07-10 14:30:00"
}
```

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `type` | string | **是** | 条目类型（见下文枚举） |
| `title` | string | 否 | 标题 |
| `description` | string | 否 | 描述/数值 |
| `category` | string | 否 | 颜色分类：`green` / `red` / `yellow` |
| `valence` | string | 否 | 附加值（正负向/卡路里） |
| `recorded_at` | string | **是** | 记录时间 |

**条目类型 (`type`) 枚举：**

| 值 | 含义 |
|----|------|
| `thought` | 随记 |
| `asset` | 资产记录 |
| `exercise` | 运动 |
| `discipline` | 自律打卡 |
| `nosugar` | 禁糖打卡 |
| `uric` | 尿酸记录 |
| `reading` | 读书 |
| `movie` | 影视 |
| `todo` | 待办 |

**响应：** 201 Created，返回创建的条目（含自增 `id`）。

### 删除条目

```
DELETE /api/entries/{id}
```

**响应：** `{"code": 0, "data": {"status": "ok"} }`

---

## 设置 API

### 读取配置

```
GET /api/settings/{key}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "key": "panel-order",
    "value": "[\"thought\",\"asset\",\"exercise\"]"
  }
}
```

不存在时 `value` 为空字符串。

### 写入配置

```
PUT /api/settings/{key}
Content-Type: application/json

{ "value": "[\"thought\",\"exercise\",\"asset\"]" }
```

**说明：** 使用 `INSERT OR REPLACE` 语义，自动创建或覆盖。

**内置 key：**

| key | 用途 |
|-----|------|
| `panel-order` | 添加面板条目的 JSON 排序数组 |

---

## DeepSeek 查询 API

### 书籍信息查询

```
POST /api/book-info
Content-Type: application/json

{ "book_name": "三体" }
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "author": "刘慈欣",
    "nationality": "中国",
    "word_count": "90",
    "first_publish_date": "2006",
    "tags": ["科幻", "硬科幻", "雨果奖"]
  }
}
```

### 影视信息查询

```
POST /api/media-info
Content-Type: application/json

{ "name": "流浪地球", "type": "movie" }
```

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 影视作品名称 |
| `type` | string | 否 | `movie`（电影）/ `series`（剧集）/ `anime`（动漫） |

**响应（电影）：**

```json
{
  "code": 0,
  "data": {
    "director": "郭帆",
    "cast": "吴京/屈楚萧/李光洁",
    "country": "中国",
    "first_release_date": "2019",
    "duration": "125",
    "tags": ["科幻", "灾难", "刘慈欣"]
  }
}
```

**响应（剧集/动漫）：**

```json
{
  "code": 0,
  "data": {
    "director": "导演名",
    "cast": "主演1/主演2",
    "country": "中国",
    "first_release_date": "2024",
    "episodes": "45min/40集",
    "tags": ["古装", "爱情", "仙侠"]
  }
}
```

> 以上两个端点均调用 DeepSeek API 获取结构化信息，需要 `DEEPSEEK_API_KEY` 环境变量。

---

## 待办 API

### 获取待办列表

```
GET /api/todos
```

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "title": "买菜",
      "completed": false,
      "category": "today",
      "due_date": "2026-07-10 20:00:00",
      "created_at": "2026-07-10 09:00:00",
      "updated_at": "2026-07-10 09:00:00"
    }
  ]
}
```

按 `created_at DESC` 排序。

### 创建待办

```
POST /api/todos
Content-Type: application/json

{ "title": "买菜", "category": "today", "due_date": "2026-07-10 20:00:00" }
```

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `title` | string | **是** | 待办内容（最长50字） |
| `category` | string | 否 | 分类：`today` (今天) / `week` (本周) / `month` (本月) / `year` (今年)，默认 `today` |
| `due_date` | string | 否 | 截止日期时间，格式 `YYYY-MM-DD HH:MM:SS` |

**响应：** 201 Created，返回完整 Todo 对象。

### 更新待办

```
PUT /api/todos/{id}
Content-Type: application/json
```

**更新标题：**

```json
{ "title": "买菜和水果" }
```

**切换完成状态：**

```json
{ "completed": true }
```

**同时更新两者：**

```json
{ "title": "买菜和水果", "completed": false }
```

### 删除待办

```
DELETE /api/todos/{id}
```

---

## 数据导出/导入 API

### 导出全部数据

```
GET /api/data/export
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "entries": [ ... ],
    "todos": [ ... ],
    "settings": [
      { "key": "panel-order", "value": "[\"thought\",\"asset\"]" }
    ]
  }
}
```

返回所有条目、待办、设置的完整 JSON 快照，可保存为文件备份。

### 导入数据

```
POST /api/data/import
Content-Type: application/json
```

将导出的 JSON 作为请求体发送。会先清空当前所有数据，再写入导入的数据。保留原始 ID 和时间戳。

> 用于：生产环境 → 开发环境的数据迁移、schema 变更后的数据重建。

---

## Go 后端 recalculate 端点

### 重算连续天数

```
POST /api/entries/recalculate?type=discipline
```

对指定 `type` 的所有条目按 `recorded_at` 升序重算连续天数，将结果写入 `description` 字段。

**参数：**

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `type` | string | 是 | 条目类型（`discipline` / `nosugar`） |

**响应：**

```json
{ "code": 0, "message": "ok" }
```

> 生产环境：此端点由 Go `EntryHandler.RecalculateEntries` 处理。
> 开发环境：Vite dev 中间件（`vite.config.js` `recalculatePlugin`）直接读写 SQLite 提供此端点，同时 `server.proxy` 将 `/api` 请求转发到 `localhost:8080`，Go 后端也处理此路由。

---

## 错误码

| HTTP 状态码 | 含义 |
|-------------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 405 | 方法不允许 |
| 500 | 服务器内部错误 |

---

## 示例流程

### 快捷自律打卡

```bash
# 1. 检查上次打卡
curl "http://localhost:8080/api/entries?limit=30"

# 2. 创建打卡
curl -X POST http://localhost:8080/api/entries \
  -H "Content-Type: application/json" \
  -d '{"type":"discipline","title":"自律","description":"7","category":"green","recorded_at":"2026-07-10 08:00:00"}'
```

### 补卡 + 重算

```bash
# 1. 补录三天前的打卡
curl -X POST http://localhost:8080/api/entries \
  -H "Content-Type: application/json" \
  -d '{"type":"nosugar","title":"禁止糖分","description":"1","category":"green","recorded_at":"2026-07-07 08:00:00"}'

# 2. 重算连续天数（Go 后端 / 开发环境下也可走 Vite 中间件由 proxy 转发）
curl -X POST http://localhost:8080/api/entries/recalculate?type=nosugar
```

### 导出并导入数据

```bash
# 导出全量数据
curl http://localhost:8080/api/data/export > backup.json

# 导入数据（清空后写入，保留原始 ID）
curl -X POST http://localhost:8080/api/data/import \
  -H "Content-Type: application/json" \
  -d @backup.json
```
