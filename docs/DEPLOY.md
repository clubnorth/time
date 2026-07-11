# 时间轴 · 部署指南

> 目标：单文件部署到 Linux 服务器，局域网内通过手机/iPad 访问

---

## 架构

```
┌──────────────────────────────────────────────┐
│                Go 二进制 (time-server)        │
│                                              │
│  ┌─────────────┐    ┌──────────────────────┐ │
│  │  API 处理器  │    │ 静态文件服务          │ │
│  │  /api/*     │    │ dist/* → SPA fallback │ │
│  └─────────────┘    └──────────────────────┘ │
│         │                    │               │
│         ▼                    ▼               │
│  ┌──────────┐    ┌─────────────────┐         │
│  │  SQLite   │    │  Vue 前端 bundle │         │
│  │  time.db  │    │  index.html +   │         │
│  └──────────┘    │  assets/*        │         │
│                  └─────────────────┘         │
└──────────────────────────────────────────────┘
```

前端构建产物 `dist/` 与 Go 二进制放在同一服务器上，Go 同时处理 API 和静态文件。

---

## 1. 构建

### 在 Windows 上交叉编译到 Linux

```powershell
# 1. 构建前端
npm run build

# 2. 交叉编译 Go → Linux amd64
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"
go build -C server -o time-server .

# 3. 确认产物
ls server/time-server
```

> `CGO_ENABLED=0` 是关键：`modernc.org/sqlite` 是纯 Go 实现，无需 C 编译器即可交叉编译。

### 在 Linux 上直接编译

```bash
cd time/server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o time-server .
```

---

## 2. 部署到 Linux 服务器

### 目录结构

```
/home/user/time/
├── time-server        # Go 二进制
├── dist/              # 前端构建产物
├── db/
│   ├── schema.sql     # 建表 SQL（参考用）
│   └── migrations/    # 迁移文件
│       └── 001_init.sql
├── data/              # SQLite 数据（自动创建）
│   └── time.db
└── config.json        # 运行配置
```

### 上传文件

```bash
# 从 Windows 上传到 Linux（使用 scp）
scp server/time-server user@192.168.1.100:/home/user/time/
scp -r dist user@192.168.1.100:/home/user/time/
scp server/db/schema.sql user@192.168.1.100:/home/user/time/db/
scp -r server/db/migrations user@192.168.1.100:/home/user/time/db/
scp server/config.json user@192.168.1.100:/home/user/time/
```

### 配置文件

```json
{
  "host": "0.0.0.0",
  "port": 8080,
  "db_path": "./data/time.db",
  "static_dir": "./dist"
}
```

| 字段 | 说明 | 默认值 |
|------|------|-------|
| `host` | 绑定地址 | `0.0.0.0`（所有网卡） |
| `port` | 监听端口 | `8080` |
| `db_path` | SQLite 数据库路径 | `./data/time.db` |
| `static_dir` | 前端文件目录 | `./dist` |

---

## 3. 启动

```bash
cd /home/user/time

# 给二进制添加执行权限
chmod +x time-server

# 直接运行
./time-server

# 或后台运行
nohup ./time-server > server.log 2>&1 &

# 或使用 systemd（见下方）
```

---

## 4. 访问

局域网内任意设备浏览器打开：

```
http://<服务器IP>:8080
```

例如：`http://192.168.1.100:8080`

---

## 5. systemd 服务（推荐）

创建 `/etc/systemd/system/time.service`：

```ini
[Unit]
Description=Time Server
After=network.target

[Service]
Type=simple
User=youruser
WorkingDirectory=/home/youruser/time
ExecStart=/home/youruser/time/time-server
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动：

```bash
sudo systemctl daemon-reload
sudo systemctl enable time
sudo systemctl start time
sudo systemctl status time
```

---

## 6. 防火墙

```bash
# 开放端口
sudo ufw allow 8080/tcp

# 或 firewalld
sudo firewall-cmd --add-port=8080/tcp --permanent
sudo firewall-cmd --reload
```

---

## 7. 更新部署

```bash
# 在开发机上
npm run build
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"
go build -C server -o time-server .

# 停止服务
ssh user@192.168.1.100 "sudo systemctl stop time"

# 上传新文件
scp server/time-server user@192.168.1.100:/home/user/time/
scp -r dist/* user@192.168.1.100:/home/user/time/dist/

# 重启服务
ssh user@192.168.1.100 "sudo systemctl restart time"
```

---

## 8. nginx 反向代理（可选）

如果需要用 80 端口或加 HTTPS：

```nginx
server {
    listen 80;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    client_max_body_size 10m;
}
```

此时 `config.json` 中 `host` 可设为 `127.0.0.1`。

---

## 9. 数据备份

```bash
# 只需备份 SQLite 数据库文件
cp /home/user/time/data/time.db /backup/time_$(date +%Y%m%d).db

# 定时备份（crontab）
0 2 * * * cp /home/user/time/data/time.db /backup/time_$(date +%Y%m%d).db
```

---

## 10. 数据迁移工作流

### 日常开发：开发与生产隔离

开发机（Windows）和服务器（Linux）使用各自独立的数据文件：

```
开发机: server/data/time.db     ← 测试数据
服务器: ~/time/server/data/time.db  ← 生产数据
```

上传 `server/` 目录时**不要覆盖 `data/time.db`**，使用：

```bash
# 仅上传代码，跳过数据目录
rsync -av --exclude='data/' server/ user@server:~/time/server/
```

### 备份生产数据

```bash
# 在服务器上
curl http://192.168.31.191:8080/api/data/export > backup_$(date +%Y%m%d).json
```

### 从生产同步到开发

```bash
# 1. 在服务器上导出
ssh user@server "curl http://localhost:8080/api/data/export" > prod_data.json

# 2. 在开发机导入
curl -X POST http://localhost:8080/api/data/import \
  -H "Content-Type: application/json" \
  -d @prod_data.json
```

### Schema 变更时的数据迁移

当 `db/migrations/` 下新增迁移文件时：

```
db/migrations/
├── 001_init.sql       ← 初始建表
├── 002_add_tag.sql    ← 新增字段
└── 003_add_xxx.sql    ← 未来变更
```

1. **新增字段/表（向后兼容）：** 直接提交新的迁移 SQL，服务重启后自动执行
2. **破坏性变更（不向后兼容）：** 导出数据 → 停服 → 清库 → 启动（自动跑全部迁移） → 导入数据

```bash
# 破坏性变更的标准流程
# 1. 导出生产数据
curl http://server:8080/api/data/export > backup.json

# 2. 停服 + 清库
sudo systemctl stop time
rm ~/time/server/data/time.db

# 3. 更新代码 + 重启（自动创建新库并跑全部迁移）
# 上传新文件...

# 4. 启动后导入数据
sudo systemctl start time
curl -X POST http://server:8080/api/data/import \
  -H "Content-Type: application/json" \
  -d @backup.json
```

> `schema_migrations` 表本身不会被清除（它在导入前已由迁移系统创建，导入只处理 entries/todos/settings），迁移状态在破环性变更后仍保持正确。

---

## 11. 故障排查

| 问题 | 检查 |
|------|------|
| 页面加载但 API 报错 | 检查 `config.json` 中 `static_dir` 路径是否正确 |
| 无法访问 | 检查防火墙；`host` 是否设为 `0.0.0.0` |
| 数据库不存在 | `./data/` 目录会自动创建 |
| 补卡功能报错 | Go 版本需 ≥ 1.26（recalculate 端点在 Go 中实现） |
| 跨域报错 | CORS 中间件已设 `Allow-Origin: *` |
| 导出/导入失败 | 检查服务是否在 8080 端口运行 |
| 迁移失败 | 查看服务器日志 `cat server.log`，检查 `migrations_dir` 路径 |
