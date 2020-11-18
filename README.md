# PHOENIX

[![Go Report Card](https://goreportcard.com/badge/github.com/lyonsdpy/phoenix?style=flat-square)](https://goreportcard.com/r/report/github.com/lyonsdpy/phoenix)
[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green)](https://github.com/lyonsdpy/phoenix/LICENSE)

![phoenix Logo](logo/phoenix_tran.png)

PHOENIX 是一个运维工具，具备如下功能：

* 使用多种方式从目标服务器或网络设备上进行信息采集、配置下发、文件上传下载(目前支持:`icmp`,`snmp`,`ssh`,`netconf`,`scp`,`sftp`)
* 自身支持计划任务，能够定时或周期性的执行相应的采集动作，能够作为运维系统、监控系统的采集器；如执行每日配置备份、每秒ping连通性检查等；且可以灵活的查询和管理已经下发的计划任务
* 支持restful-API的同步、异步交互；目前异步模式仅支持定义http回调的方式，上层服务端还需要实现相应的接收方法

**Note**:

* 当多个请求对同一台设备进行snmp交互时存在锁机制，获取到锁时才能对设备进行操作，避免设备无法处理并发的snmp请求(尤其是交换机、路由器等网络设备)
* 当多个请求对同一台设备进行ssh、netconf、scp、sftp交互时，存在一个全局可配置的信号量，避免并发的请求超出设备的vty限制数量(默认为4)
* 对于`netconf`与`ssh_shell`方式的请求，系统提供了一个可配置超时时长的连接池，在超时时间内会复用已经建立的ssh连接，不会发起新的ssh协商，适用于短时间内多次频繁的请求

**TODO LIST**:

* **增加任务结果消息队列**:增加本地缓存的异步方式，任务执行后将结果放在本次缓存中，由北向的服务端在后续消费；该缓存支持设置最大占用内存和超时时长；类似消息中间件的生产-消费者
* **prometheus-exporter支持**:用户可以指定该任务为一个监控任务，并将采集的数据按prometheus保存在上面所述的缓存中
* **open-falcon支持**:支持将监控的指标以open-falcon的数据格式上报到指定目标
* **snmp结果缓存**:考虑到北向的不同系统或逻辑可能会反复的采集同一设备的相同非监控指标(如接口名称、索引等配置信息)，后续会支持snmp请求时增加一个缓存开关(借鉴CDN原理)指明该指标是否可缓存，在缓存时间内，对相同设备相同指标的snmp请求，由缓存直接应答，减轻设备负担；同样用户也可强制请求直接回原获取并更新缓存
* 支持snmp-trap上报
* 支持syslog上报
* 支持telemetry上报

## 安装

### 源码安装

```sh
git clone github.com/lyonsdpy/phoenix
cd phoenix
make win # windows环境
make linux # linux环境
```

>源码安装要求具备go环境，版本>=1.14

### 获取可执行文件

访问 `github.com/lyonsdpy/phoenix/release` 获取已编译的可执行文件

## 使用

直接运行 `./phoenix` ,即可开启使用默认参数的进程

### 查看帮助

使用 `./phoenix -help` 或 `./phoenix -h` 查看参数设置帮助

### API文档

本工具使用go swag管理restful-API文档，在程序运行后访问如下地址查看API文档的使用方式:

`http://127.0.0.1:[port]/swagger/index.html`
