# Olympus - 私有制品中转仓

Olympus 是一个支持多种包管理器的私有制品中转仓库, 提供中央仓库代理加速与私有制品分发服务.

## 支持列表

* Maven Registry
* NPM Registry
* Go Registry

## 功能

1. 制品: 查询列表和详情并针对某一制品进行拉取和删除
2. 仓库: 代理用户拉取请求, 下载制品并定时清理
3. 安全: 检查制品漏洞并予以提示
4. 统计: 记录制品拉取频率
5. 用户: 账号密码登陆、OAuth2登陆

## 构建

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o olympus -ldflags '-s -w'
```
