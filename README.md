# Olympus - 私有制品代理加速仓库

Olympus 是一个支持多种包管理器, 允许私有化部署的制品代理加速仓库, 提供代理加速与缓存服务.

## 支持列表

* Maven Registry
* NPM Registry
* Go Registry

## 构建

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o olympus -ldflags '-s -w'
```
