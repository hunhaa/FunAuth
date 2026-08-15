# FunAuth 使用指南

FunAuth 是一个支持网易我的世界 3.9 版本的验证服务器和租赁服客户端工具。

## 📦 文件说明

- `funauth-server` - 验证服务器（Linux amd64）
- `funauth-client` - 租赁服客户端（Linux amd64）

## 🚀 快速开始

### 1. 启动验证服务器

```bash
# 赋予执行权限
chmod +x funauth-server

# 启动服务器（默认端口 8080）
./funauth-server

# 或指定端口
./funauth-server --port 9000
```

**服务器功能：**
- `/api/phoenix/login` - 登录认证
- `/api/phoenix/transfer_check_num` - 传输校验码
- `/api/open/g79/rental_search` - 租赁服搜索
- `/api/open/g79/rental_details` - 租赁服详情
- `/api/open/g79/rental_enter_world` - 进入租赁服
- 更多 OpenAPI 接口...

**支持的协议版本：**
- Phoenix 协议（大厅/登录）
- G79 协议（租赁服）
- SDK Version: 3.9.0

### 2. 使用租赁服客户端

```bash
# 赋予执行权限
chmod +x funauth-client

# 启动客户端
./funauth-client
```

**客户端菜单选项：**
```
FunAuth 租赁服客户端 v3.9.0
==========================
1. 搜索租赁服
2. 查看租赁服详情
3. 进入租赁服
4. 查看可用租赁服列表
5. 退出
请选择操作：
```

## 🔧 配置说明

### 环境变量（可选）

```bash
# 服务器配置
export FUNAUTH_PORT=8080        # 服务器端口
export FUNAUTH_HOST=0.0.0.0     # 监听地址

# 启动服务器
./funauth-server
```

## 📝 使用示例

### 搜索租赁服

**通过客户端：**
```bash
./funauth-client
# 选择选项 1
# 输入租赁服名称关键词
```

**通过 API：**
```bash
curl -X POST http://localhost:8080/api/open/g79/rental_search \
  -H "Content-Type: application/json" \
  -d '{"keyword":"生存","page":1}'
```

### 进入租赁服

**通过客户端：**
```bash
./funauth-client
# 选择选项 3
# 输入租赁服 ID
# 获取服务器地址和认证信息
```

**通过 API：**
```bash
# 1. 获取租赁服详情
curl -X POST http://localhost:8080/api/open/g79/rental_details \
  -H "Content-Type: application/json" \
  -d '{"rental_id":"你的租赁服ID"}'

# 2. 进入租赁服
curl -X POST http://localhost:8080/api/open/g79/rental_enter_world \
  -H "Content-Type: application/json" \
  -d '{"rental_id":"你的租赁服ID"}'
```

## 🔐 认证方式

FunAuth 支持两种认证方式：

### 1. Cookie 认证（推荐）
适用于已有网易账号 Cookie 的用户

```bash
curl -X POST http://localhost:8080/api/open/g79/rental_search \
  -H "Cookie: your_netease_cookie_here" \
  -H "Content-Type: application/json" \
  -d '{"keyword":"生存"}'
```

### 2. Bearer Token 认证
适用于通过登录接口获取 token 的用户

```bash
# 先登录获取 token
curl -X POST http://localhost:8080/api/phoenix/login \
  -H "Content-Type: application/json" \
  -d '{"username":"your_username","password":"your_password"}'

# 使用 token 访问其他接口
curl -X POST http://localhost:8080/api/open/g79/rental_search \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json" \
  -d '{"keyword":"生存"}'
```

## 🌐 API 接口列表

### Phoenix 协议（登录/大厅）
- `POST /api/phoenix/login` - 用户登录
- `POST /api/phoenix/transfer_check_num` - 传输校验码验证
- `POST /api/phoenix/transfer_start_type` - 传输起始类型
- `POST /api/phoenix/tan_lobby_login` - Tan Lobby 登录
- `POST /api/phoenix/tan_lobby_create` - 创建 Tan Lobby
- `POST /api/phoenix/tan_lobby_transfer_server` - Tan Lobby 传输服务器

### Open API（租赁服/用户）
- `POST /api/open/g79/user_detail` - 用户详情
- `POST /api/open/g79/rental_search` - 租赁服搜索
- `POST /api/open/g79/lobby_room` - 大厅房间
- `POST /api/open/g79/rental_available` - 可用租赁服列表
- `POST /api/open/g79/rental_details` - 租赁服详情
- `POST /api/open/g79/user_settings` - 用户设置
- `POST /api/open/g79/user_search` - 用户搜索
- `POST /api/open/g79/download_info` - 下载信息
- `POST /api/open/g79/rental_enter_world` - 进入租赁服

### 其他接口
- `GET /api/new` - 生成 UUID 或验证 Cookie

## ⚠️ 注意事项

1. **网络连接**：确保服务器可以访问网易我的世界官方服务器
2. **Cookie 有效性**：使用 Cookie 认证时，确保 Cookie 未过期
3. **端口占用**：启动服务器前确认端口未被占用
4. **防火墙**：如需外网访问，请配置防火墙规则
5. **版本兼容**：本工具支持网易 MC 3.9.0 SDK 版本

## 🐛 故障排除

### 服务器无法启动
```bash
# 检查端口占用
netstat -tlnp | grep 8080

# 查看日志
./funauth-server --verbose
```

### 无法连接租赁服
1. 检查 Cookie 是否有效
2. 确认租赁服 ID 正确
3. 查看服务器日志获取详细错误信息

### API 返回认证失败
- 检查 Cookie 或 Token 是否正确
- 确认认证头格式正确
- 尝试重新登录获取新的 Token

## 📚 技术支持

如遇问题，请提交 Issue 到 GitHub 仓库，并提供：
- 错误日志
- 复现步骤
- 环境信息（操作系统、Go 版本等）

## 📄 许可证

本项目遵循开源许可证，详见 LICENSE 文件。

---

**版本信息：**
- FunAuth Version: 1.0.0
- Supported SDK: 3.9.0
- Build Date: 2024
- Platform: Linux amd64
