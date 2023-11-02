# quic基本性能测试
## 测试环境
- server
```
CPU(s):                8
Model name:            Intel(R) Xeon(R) CPU           X5650  @ 2.67GHz
```

- client
```
the same as server
```
- 测试目标
  1. 协程资源消耗 
  2. cpu消耗
  3. 性能分析
  4. 异常分析

- 测试代码
https://github.com/someview/quic-benchmark
```
go run ./server/main.go
go run ./client/main.go
```
mode: multi, 多协程频繁调度
mode: single, 多协程，某个协程频繁调度
mode: silent，多协程，不调度

## 测试结果
### 正常情况

```

mode:multi, clientNum:1000, routineNum: 3000, paylaod: 13字节 总效率(send+recv)/2: 3.5e6/s  ws: 2.0e5  cpu: 200%
mode: multi, clientNum:1e4, routinueNum: 30000,    payload: 13字节 总效率(send+recv)/2:ws1-2倍    cpu:  %%200-300%
mode: multi，clientNum:1e4(一半空闲), routinueNum: 30000, payload: 13字节,总效率(send+recv)/2:ws5倍 cpu: 200%
mode: 
```
## 基本分析
```
cpu: 基本是是ws的两倍, 传输效率: ws的10倍不到,协程是否稳定: 稳定
单条stream有限流，并且流控很快耗尽
```
## 可行性分析
- 架构
  契合当前im系统结构，只需要少量改动就能满足需求  
- 功能 
  不稳定,go的webtransport实现仍有bug,协程数暴涨，连接不能大量建立
  转发效率远高于websocket，相同cpu消耗下高了一个数量级
- 安全性
  只支持tls,不支持http，内网使用
- 环境
  需要负载均衡器支持quic、http3
- 兼容性问题
  当前草案阶段


