# RTNS

[![GoDoc](https://godoc.org/github.com/RTradeLtd/rtns?status.svg)](https://godoc.org/github.com/RTradeLtd/rtns) [![codecov](https://codecov.io/gh/RTradeLtd/rtns/branch/master/graph/badge.svg)](https://codecov.io/gh/RTradeLtd/rtns) [![Build Status](https://travis-ci.com/RTradeLtd/rtns.svg?branch=master)](https://travis-ci.com/RTradeLtd/rtns) ![GitHub release](https://img.shields.io/github/release/RTradeLtd/rtns.svg?style=flat-square)

RTNS (RTrade 命名系统) 是一个独立的IPNS记录管理服务，通过与[kaas](https://github.com/RTradeLtd/kaas)搭配使用来增强IPNS记录发布的安全性..

它本质上是 [go-ipfs/namesys](https://github.com/ipfs/go-ipfs/tree/master/namesys) 的优化版。

## 多语言

[![](https://img.shields.io/badge/Lang-English-blue.svg)](README.md)  [![jaywcjlove/sb](https://jaywcjlove.github.io/sb/lang/chinese.svg)](README-zh.md)

## 用例

导入依赖包:

```Golang
import "github.com/RTradeLtd/rtns"
```
你需要初始化一个libp2p服务，并提供一个分布式哈希表作为构造函数的参数传入。该依赖包含一个辅助工具，可以将KaaS客户端转化为有效的密钥库接口类型。

## 开发

### 使用`$GOPATH`

可参考: https://splice.com/blog/contributing-open-source-git-repositories-go/

1. 拉取仓库新分支
2. 克隆仓库： `git clone git@github.com:RTradeLtd/rtns.git $GOPATH/src/github.com/RTradeLtd/rtns`
3. 运行：`cd $GOPATH/src/github.com/RTradeLtd/rtns`
3. 启动远端服务：

```bash
git remote rename origin upstream
git remote add origin git@github.com:<your-github-username>/rtns.git
```
4. 添加全局环境变量：`export GO111MODULE=on` to `.bashrc` or `.bash_profile` (if you're on a Mac) or `.zshrc` (if you're using [zsh](https://github.com/robbyrussell/oh-my-zsh))。运行 `source <rc-file>`来确保
5. 运行 `go mod download` 来安装依赖
6. 运行测试脚本： `go test ./...`

### 外部 $GOPATH

1. 在你机器的任意位置拉取并克隆新分支
2. 运行： `cd rtns`
3. 新建一个远程流仓库

```bash
git remote add upstream git@github.com:RTradeLtd/rtns.git
```
4. 运行`go mod download` 来下载依赖
5. 运行 `go test ./...`进行测试

## 限制

* 在Temporal中使用时，因为KaaS系统故障所产生的任意密钥都不适合自动重新发布

## 未来发展

* DNSLink支持
* 以TNS (Temporal Name Server) 网关的形式表现。
* 启用KaaS后端监听服务支持
  * 在遇到上述KaaS故障和错误时，将反复轮训检索所有可用KaaS主机的私钥信息。
* 为IPNS pubsub启用自动主题订阅
  * 这将涉及使用rtfs来调用IPFS节点，为主题建立订阅信息。