### Study Plan

- Day 1 — Go 并发模型 & Goroutine 心智模型

核心概念

	•	Go concurrency ≠ parallelism
	•	M:N 调度模型（G / M / P）
	•	goroutine 创建成本 & 生命周期
	•	什么时候用并发，什么时候不用

- Day 2 — Channel 基础（Blocking 是核心）

核心概念

	•	channel 是 typed, blocking queue
	•	send / receive 的阻塞语义
	•	unbuffered vs buffered

Project: 用channel实现producer consumer pattern

- Day 3 — Channel Patterns（非常重要）

核心并发模式

	•	fan-out：一个 input → 多个 worker
	•	fan-in：多个 worker → 一个 output
	•	pipeline：stage by stage

- Day 4 — Select & 超时 & 非阻塞

核心概念

```
select {
case v := <-ch:
case <-time.After(1 * time.Second):
default:
}
```

	•	select 是 并发控制器
	•	timeout / try-receive / try-send
	•	fairness（select 随机选 ready case）

1. 写一个worker 如果500ms内没活干就退出

2. 实现一个non blocking send 

- Day 5 — Context（生产级必会）

核心概念

- context用途

- timeout

- request scoped value

- 谁创建 谁cancel 

- 不要存大对象 

练习: 改造day 3 的pipeline 支持整体cancel 以及cancle之后不泄露go routine

- Day 6 — sync 包（Mutex / WaitGroup / Once）

核心概念

- mutex

- conditional var 

练习: implementing blocking queue

- Day 7 — 常见 Bug & Anti-Patterns（救命）

- Day 8 — Worker Pool & Backpressure（工程核心）

- Day 9 — 并发设计思维（不是语法）

- Day 10 — 总复盘 & 面试准备

