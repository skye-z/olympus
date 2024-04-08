# Olympus - 私有制品代理加速仓库

Olympus 是一个支持多种包管理器, 允许私有化部署的制品代理加速仓库, 提供代理加速与缓存服务.

## 适用场景

使用 Olympus 可以免于单独配置各个设备的网络代理, 仅需将源地址指向 Olympus 即可;

并且在统一的仓库管理中, 相同的包不再会重复拉取, 有助于降低旧包拉取时间于流量消耗.

## 支持列表

* Maven Registry
* NPM Registry
* Go Registry

## 构建

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o olympus -ldflags '-s -w'
```
