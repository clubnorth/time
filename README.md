# 时间轴 Time · v1.1.0

> 个人时间管理与待办事项应用 · 手机 + 平板优先

---

## 功能

- **时间轴** — 随记（太阳/阴雨）、运动、自律打卡、禁糖打卡、尿酸、资产、读书（DeepSeek AI）、影视（DeepSeek AI）
- **统计** — 8 类 GitHub 风格热力图 + 统计卡片
- **待办** — 四分类自动归类、截止时间、左滑操作、完成联动时间轴
- **补卡** — 回填遗漏记录，自动重算连续天数
- **排序** — 面板条目自由排序，保存到数据库
- **评分** — 1-10 分星星评分，DeepSeek 自动查询 + 标签

## 技术栈

| 层 | 技术 |
|----|------|
| 前端 | Vue 3 + Vite 8 |
| 后端 | Go 1.26 + SQLite |
| AI | DeepSeek API 集成 |
| 通信 | RESTful JSON API |
| 端口 | 前端 5173 · 后端 8080 |

## 快速启动

```powershell
# 后端
cd server
go build -o time-server.exe .
.\time-server.exe

# 前端
npm install
npm run dev -- --host
```

浏览器打开 http://localhost:5173

## 文档

| 文档 | 内容 |
|------|------|
| [架构设计](docs/ARCHITECTURE.md) | 项目结构、组件树、数据流 |
| [API 参考](docs/API.md) | REST API + DeepSeek 端点 |
| [数据库设计](docs/DATABASE.md) | ER 图、字段说明、迁移系统 |
| [开发指南](docs/DEVELOPMENT.md) | 环境搭建、约定、DeepSeek 集成 |
| [组件文档](docs/COMPONENTS.md) | 15 个 Vue 组件 Props/Emits |
| [设计系统](docs/DESIGN.md) | 配色、断点、动画、星星评分 |
| [部署指南](docs/DEPLOY.md) | Linux 部署、systemd、数据迁移 |
| [变更日志](docs/CHANGELOG.md) | 版本历史 |

## 项目结构

```
time/
├── src/           # Vue 前端 (15 组件)
├── server/        # Go 后端 (API + DeepSeek)
├── docs/          # 软件工程文档
├── vite.config.js # Vite 配置 + 代理
└── package.json
```
