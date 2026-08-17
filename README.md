# NetShieldQ

NetShieldQ 是一个包含 Vue/Electron 客户端和 Go 服务端的安全通信项目。

## 目录

- `netshieldq/`：Vue 2 + Electron 客户端
- `pnetshieldq/`：Go 服务端、Swagger 文档及数据库初始化脚本

## 本地配置

复制 `pnetshieldq/config/config.example.json` 为
`pnetshieldq/config/config.json`，并填写本机数据库、MQTT 和 OSS 配置。
真实配置文件已被 Git 忽略，请勿提交密码、Token 或云服务密钥。

启动后端前还需要设置 `JWT_SIGNING_KEY` 环境变量，使用足够长的随机值。

## 前端

```sh
cd netshieldq
npm install
npm run serve
```

## 后端

```sh
cd pnetshieldq
go run .
```
