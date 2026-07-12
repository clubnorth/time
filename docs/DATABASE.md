# 时间轴 · 数据库设计

> SQLite 3 · 纯 Go 驱动 (`modernc.org/sqlite`) · CGo-free

---

## ER 图

```
┌──────────────────────────────────┐
│           entries                 │
├──────────────────────────────────┤
│ id           INTEGER  PK  AUTO   │
│ type         TEXT     NOT NULL   │──┐
│ title        TEXT                 │  │ 常见值:
│ description  TEXT                 │  │ thought / asset / exercise
│ category     TEXT                 │  │ discipline / nosugar / uric
│ valence      TEXT     NULL        │  │ reading / movie / todo
│ created_at   TEXT                 │  │
│ updated_at   TEXT                 │  │
├──────────────────────────────────┤  │
│ INDEX: idx_entries_type          │◄─┘
│ INDEX: idx_entries_recorded      │
└──────────────────────────────────┘


┌──────────────────────────────────┐
│           settings               │
├──────────────────────────────────┤
│ key          TEXT     PK          │
│ value        TEXT     NOT NULL   │
└──────────────────────────────────┘


┌──────────────────────────────────┐
│            todos                 │
├──────────────────────────────────┤
│ id           INTEGER  PK  AUTO   │
│ title        TEXT     NOT NULL   │
│ completed    INTEGER  DEFAULT 0  │
│ category     TEXT     DEFAULT '' │
│ due_date     TEXT     DEFAULT '' │
│ created_at   TEXT                 │
│ updated_at   TEXT                 │
└──────────────────────────────────┘


┌──────────────────────────────────┐
│       schema_migrations          │
├──────────────────────────────────┤
│ version      INTEGER  PK         │
│ name         TEXT     NOT NULL   │
│ applied_at   TEXT     NOT NULL   │
└──────────────────────────────────┘
```

---

## 1. `entries` — 时间轴条目

```sql
CREATE TABLE IF NOT EXISTS entries (
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

CREATE INDEX IF NOT EXISTS idx_entries_type      ON entries(type);
CREATE INDEX IF NOT EXISTS idx_entries_recorded  ON entries(recorded_at);
```

### 字段说明

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | INTEGER | PK AUTOINCREMENT | 自增主键 |
| `type` | TEXT | NOT NULL | 条目类型枚举 |
| `title` | TEXT | DEFAULT '' | 展示标题 |
| `description` | TEXT | DEFAULT '' | 内容/数值/连续天数 |
| `category` | TEXT | DEFAULT '' | 颜色标签：`green`/`red`/`yellow` |
| `valence` | TEXT | NULL | 附加信息（正负向/卡路里值） |
| `recorded_at` | TEXT | NOT NULL | 用户指定的记录时间 |
| `created_at` | TEXT | DEFAULT NOW | 创建时间（自动） |
| `updated_at` | TEXT | DEFAULT NOW | 更新时间（自动） |

### `description` 的语义（按 type）

| type | description 含义 |
|------|-----------------|
| `thought` | 念头内容文本 |
| `asset` | 资产余额数值 |
| `exercise` | 运动量（千米/个/米/分钟） |
| `discipline` | 连续自律天数 |
| `nosugar` | 连续禁糖天数 |
| `uric` | 尿酸值 |
| `reading` / `movie` | 作品描述 + 元数据 |
| `todo` | 待办事项标题 |

### 索引策略

- `idx_entries_type` — 统计查询：按类型筛选
- `idx_entries_recorded` — 核心查询：按时间排序与分页

---

## 2. `settings` — 键值配置

```sql
CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `key` | TEXT | PRIMARY KEY | 配置键名 |
| `value` | TEXT | NOT NULL | JSON 字符串值 |

### 内置 key

| key | value 格式 | 说明 |
|-----|-----------|------|
| `panel-order` | `["thought","asset",...]` | 添加面板条目的显示顺序 |

写入使用 `INSERT OR REPLACE` 语义（upsert）。

---

## 3. `todos` — 待办事项

```sql
CREATE TABLE IF NOT EXISTS todos (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT    NOT NULL,
    completed  INTEGER NOT NULL DEFAULT 0,
    category   TEXT    NOT NULL DEFAULT '',
    due_date   TEXT    NOT NULL DEFAULT '',
    created_at TEXT    NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at TEXT    NOT NULL DEFAULT (datetime('now', 'localtime'))
);
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | INTEGER | PK AUTOINCREMENT | 自增主键 |
| `title` | TEXT | NOT NULL | 待办内容 |
| `completed` | INTEGER | DEFAULT 0 | 0 = 未完成，1 = 已完成 |
| `category` | TEXT | DEFAULT '' | 分类：`today`/`week`/`month`/`year` |
| `due_date` | TEXT | DEFAULT '' | 截止日期时间 |
| `created_at` | TEXT | DEFAULT NOW | 创建时间 |
| `updated_at` | TEXT | DEFAULT NOW | 更新时间 |

> `completed` 在数据库中存 INTEGER，前端通过 JSON 序列化自动转为 boolean。

---

## 4. `schema_migrations` — 迁移版本跟踪

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    name       TEXT    NOT NULL,
    applied_at TEXT    NOT NULL DEFAULT (datetime('now','localtime'))
);
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `version` | INTEGER | PRIMARY KEY | 迁移编号（从文件名前缀提取） |
| `name` | TEXT | NOT NULL | 迁移文件名（如 `001_init.sql`） |
| `applied_at` | TEXT | DEFAULT NOW | 应用时间 |

此表由 `service.RunMigrations()` 自动创建和管理。每次启动时检查 `db/migrations/` 目录下的 SQL 文件，跳过已应用的版本，只执行新增的迁移。

---

## 5. 关键查询

### 时间轴分页（前端核心查询）

```sql
-- 初始加载
SELECT * FROM entries
ORDER BY recorded_at DESC, id DESC
LIMIT 30;

-- 加载更多（游标分页）
SELECT * FROM entries
WHERE recorded_at < '2026-06-01 00:00:00'
ORDER BY recorded_at DESC, id DESC
LIMIT 20;
```

### 按月查询

```sql
SELECT * FROM entries
WHERE strftime('%Y-%m', recorded_at) = '2026-07'
ORDER BY recorded_at DESC, id DESC;
```

### 统计页全量拉取

```sql
SELECT * FROM entries
WHERE recorded_at < '2026-07-10 00:00:00'
ORDER BY recorded_at DESC, id DESC
LIMIT 100;
-- 循环直到返回 < 100 条
```

### 连续天数重算（Go service 层，应用层逻辑）

并非单条 SQL，而是在 `service/entry.go` 的 `RecalculateEntries()` 方法中遍历：
1. `SELECT id, recorded_at FROM entries WHERE type = ? ORDER BY recorded_at ASC`
2. 遍历每一行，向前查找连续日期，计算从最早断点到该行的连续天数
3. `UPDATE entries SET description = ? WHERE id = ?`

---

## 6. 数据库文件路径

- 开发环境：`server/data/time.db`
- 配置入口：`server/config.json` → `db_path`
- 默认值：`./data/time.db`（相对于 server 工作目录）

---

## 7. 模式迁移（版本化）

迁移文件位于 `server/db/migrations/`，命名格式为 `NNN_description.sql`（如 `001_init.sql`）。启动时 `service.RunMigrations(db, migrationsDir)` 自动执行未应用的迁移文件。

### 迁移列表

> 当前 v1.1.0 迁移版本：003

| 文件 | 版本 | 说明 |
|------|------|------|
| `001_init.sql` | 1 | 初始建表（entries/settings/todos/schema_migrations） |
| `002_todo_category.sql` | 2 | todos 表添加 `category` 字段（today/week/month/year） |
| `003_todo_due_date.sql` | 3 | todos 表添加 `due_date` 字段（截止日期时间） |

### 迁移原理

1. 确保 `schema_migrations` 表存在
2. 读取 `db/migrations/` 下所有 `.sql` 文件，按版本号排序
3. 跳过已在 `schema_migrations` 中记录的版本
4. 逐个执行未应用的 SQL，执行成功后插入版本记录

### 添加新迁移的流程

1. 在 `server/db/migrations/` 下创建 `004_your_change.sql`
2. 写入 DDL（`ALTER TABLE ...`、`CREATE TABLE IF NOT EXISTS ...` 等）
3. 重启后端，迁移自动执行

### 破坏性变更的数据重建

```bash
# 导出 → 停服 → 清库 → 启动（自动跑全部迁移） → 导入
curl http://localhost:8080/api/data/export > backup.json
# 停止服务，删除 time.db，重启
curl -X POST http://localhost:8080/api/data/import -H "Content-Type: application/json" -d @backup.json
```
