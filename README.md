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

- Day 5 — Context（生产级必会）

- Day 6 — sync 包（Mutex / WaitGroup / Once）

- Day 7 — 常见 Bug & Anti-Patterns（救命）

- Day 8 — Worker Pool & Backpressure（工程核心）

- Day 9 — 并发设计思维（不是语法）

- Day 10 — 总复盘 & 面试准备

