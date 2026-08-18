# des-queue

M/M/1 与 M/M/c 排队系统离散事件仿真工具。

## 功能

- M/M/1 和 M/M/c 排队系统仿真
- 理论值对照（稳态公式）
- 优先级队列、限长队列多种排队策略
- 排队网络（串联/并联）仿真
- 多种到达模式（泊松、批量、时变）
- 瞬时指标记录与分析
- HTTP API + Web UI

## 构建与运行

```bash
go build -o des-queue .
./des-queue -lambda 0.5 -mu 1.0 -servers 1 -customers 8000
./des-queue serve -addr :8080
```
