# CLI Template

![Go Version](https://img.shields.io/badge/Go-1.20-blue)

## 介绍

cli-tpl是Go语言编写的一个简单的用于快速开发命令行工具的模板，基于[cobra](https://github.com/spf13/cobra)
、[viper](https://github.com/spf13/viper)、[zap](https://github.com/uber-go/zap)  

## 要求
* Go：1.20+

## 目录

* [说明](#说明)
    * [克隆代码](#克隆代码)
    * [全局选项](#全局选项)
    * [输出默认配置文件](#输出默认配置文件)
    * [动态修改配置参数](#动态修改配置参数)
    * [配置文件格式支持](#配置文件格式支持)
* [原则](#原则)
    * [代码规范](#代码规范)
    * [标准输出和退出码](#标准输出和退出码)

## 说明

### 克隆代码

```bash
$ git clone https://github.com/vvfock3r/cli-tpl.git
$ cd cli-tpl/
$ go mod tidy
```

### 全局选项

```bash
$ go run main.go -h
Simple Command-Line Interface Template
For details, please refer to https://github.com/vvfock3r/cli-tpl

Usage:
  cli-tpl [flags]
  cli-tpl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Configuration operation

Flags:
  -c, --config-file string   config file
  -h, --help                 help message
      --log-format string    log format (default "console")
      --log-level string     log level (default "info")
      --log-output string    log output (default "stdout")
  -v, --version              version message

Use "cli-tpl [command] --help" for more information about a command.
```

### 输出默认配置文件

```bash
# 输出默认的配置
$ go run main.go config default
global:
  log:
    # level:  日志级别，支持debug,info,warn,error等
    # format: 日志格式，支持console和json
    # output: 输出位置，支持stdout,stderr或任意文件名，多个输出使用逗号分割
    level: info
    format: console
    output: stdout

# 以上配置是根据项目etc/default.yaml内容输出而来,该文件是被编译到二进制命令中的

# 以下两种启动方式是等价的
$ go run main.go # 方式一
$ go run main.go config default > default.yaml
$ go run main.go -c default.yaml # 方式二
```

### 动态修改配置参数

```bash
# 1.要求必须使用配置文件
$ go run main.go config default > default.yaml
$ go run main.go -c default.yaml

# 2.以修改日志格式举例,修改default.yaml中的 global.log.format 为 json,对比修改前后日志输出的不同    
2023-02-08 14:09:32     INFO    cmd/root.go:72  root command run
2023-02-08 14:09:33     INFO    cmd/root.go:72  root command run
2023-02-08 14:09:34     INFO    cmd/root.go:72  root command run
2023-02-08 14:09:35     WARN    viper/util.go:34        config update trigger   {"operation": "create", "filename": "/root/cli-tpl/default.yaml"}
{"level":"warn","time":"2023-02-08 14:09:35","caller":"viper/util.go:43","message":"config reload success","name":"global.log","detail":"success"}
{"level":"warn","time":"2023-02-08 14:09:35","caller":"viper/util.go:34","message":"config update trigger","operation":"write","filename":"/root/cli-tpl/default.yaml"}
{"level":"warn","time":"2023-02-08 14:09:35","caller":"viper/util.go:43","message":"config reload success","name":"global.log","detail":"success"}
{"level":"info","time":"2023-02-08 14:09:35","caller":"cmd/root.go:72","message":"root command run"}
{"level":"info","time":"2023-02-08 14:09:36","caller":"cmd/root.go:72","message":"root command run"}
{"level":"info","time":"2023-02-08 14:09:37","caller":"cmd/root.go:72","message":"root command run"}

# 3.若要修改的参数同时也使用命令行选项指定值，则该修改不生效。优先级由高到低规则如下:   
#   1.viper.Set()设置的值
#   2.命令行中读取的值
#   3.环境变量中读取的值
#   4.配置文件中读取的值
#   5.远程存储读取的值
#   6.viper.SetDefault()设置的值

# 4.对于自定义配置，需要调用以下方法将函数注册到viper
#   viperutil.RegisterWatchFunc(name, func)
#   name只用于显示日志使用，func是真正执行的函数
```

### 配置文件格式支持

```go
// viper默认支持以下格式的配置文件
//  json 
//  toml
//  yaml 
//  yml
//  properties 
//  props 
//  prop 
//  hcl 
//  tfvars 
//  dotenv 
//  env 
//  ini
//
// 如果只想支持一种或多种，请修改 cmd/root.go
//  viper.SupportedExts = []string{"yaml"}
```

## 原则

### 代码规范

```bash
# 设置Git Hooks
$ git config core.hooksPath .githooks

# 在每次提交前会执行.githooks目录下的钩子脚本，比如
$ git add * && git commit -m "test: git hooks" 
pre-commit
    RUN go mod tidy
    RUN gofmt -w -r "interface{} -> any" .
    RUN go vet .
[main 931a3e8] update
 1 file changed, 39 insertions(+), 232 deletions(-)
 rewrite README.md (94%)
```

### 标准输出和退出码

* `-h/--help`和`-v/--version`输出到**stdout**，退出码为 **0**
* 错误类信息输出到 **stderr**，退出码为 **1**
* 默认所有的日志(即使是error类型)输出到 **stdout**