# 敏感词项目说明
![Go](https://github.com/zeromicro/go-zero/workflows/Go/badge.svg?branch=master)
![Go Report Card](https://goreportcard.com/badge/github.com/zeromicro/go-zero)
![codecov](https://codecov.io/gh/zeromicro/go-zero/branch/master/graph/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/zeromicro/go-zero.svg)](https://pkg.go.dev/github.com/zeromicro/go-zero)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
## 产品要求
### 完全匹配
例如：侵权词为：nike 文案中出现nike侵权
### 不区分大小写
侵权词为：NIKE 产品文案中出现nike侵权
### 商标品牌词大于等于两个单词单词中间加符号、空格算侵权
侵权词：speed bridge 产品文案中出现 speed &&& bridge侵权 （特殊符号包含： '~', '`', '@', '#', '$', '%', '^', '&', '*', '-', '_', '+', '=', ':', '\\', '|', '/', ' ', '\r', '\n', '·', '.'）
### 单个英语单词分开不算侵权
ni ke
ni、ke
### 单个英语单词前后加字母、数字不算侵权，前后加符号算侵权
nike123不侵权
anikeb不侵权
&nike？算侵权
### 中文日语完整词中间加符号不算侵权
侵权词：阿迪达斯  产品文案中出现 阿迪&达斯 不算侵权
## 启动命令
.\程序包 -f etc/test.yaml
## sensitive-word-service/etc
配置文件所在目录：
test.yaml 测试环境配置文件
pre.yaml  预发环境配置文件
prod.yaml 生产环境配置文件

配置文件配置的信息主要配置的是nacos的信息，实际程序配置全部在对应的nacos

## sensitive-word-service\internal\config
配置结构体目录

## sensitive-word-service\internal\core
敏感词匹配功能核心目录
extend.go 主要是根据业务要求定制的一些匹配扩展逻辑
trie.go   AC自动机算法实现

## sensitive-word-service\internal\infra\repo
基础设施目录下的仓库目录
sensitive_word.go rpc调用敏感词库

## sensitive-word-service\internal\jsonrpc
jsonrpc服务目录，主要存放jsonrpc调用入口类

## sensitive-word-service\internal\listener
word_change_listener.go rabbitMq监听器，主要是监听消息，触发重新加载敏感词

## sensitive-word-service\internal\logic
sensitive_word.go 敏感词业务逻辑类

