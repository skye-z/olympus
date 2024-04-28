# Olympus - 私有制品代理加速仓库

Olympus 是一个支持多种包管理器, 允许私有化部署的制品代理加速仓库, 提供代理加速与缓存服务.

## 开发计划

* [ ] 移除制品
* [ ] OAuth2 登陆
* [ ] 包管理器层面鉴权

## 适用场景

使用 Olympus 可以免于单独配置各个设备的网络代理, 仅需将源地址指向 Olympus 即可;

并且在统一的仓库管理中, 相同的包不再会重复拉取, 有助于降低旧包拉取时间于流量消耗.

## 支持列表

* Maven Registry
* NPM Registry
* Go Registry

## 安装

请复制下方命令到服务器终端中执行, 脚本提供了 Olympus 的安装、卸载与开启自启设置服务

```shell
bash -c "$(curl -fsSL https://betax.dev/sc/olympus.sh)"
```

## 控制

```shell
# 启动 Harbor
systemctl start olympus
# 停止 Harbor
systemctl stop olympus
# 查看 Harbor 状态与日志
systemctl status olympus
```

## 构建

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o olympus -ldflags '-s -w'
```
