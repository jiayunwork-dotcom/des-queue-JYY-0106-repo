# des-queue

离散事件仿真（DES）与排队论分析 CLI：M/M/1 与 M/M/c 排队系统。

事件驱动引擎使用基于最小堆的未来事件表（FEL），到达/服务间隔为固定 seed 的指数分布随机变量。仿真输出利用率 ρ、平均队长 L、平均逗留时间 W、平均等待 Wq 与吞吐量，并与 M/M/1 理论公式对照相对误差：

- ρ = λ/μ
- L = ρ/(1-ρ)
- W = 1/(μ-λ)
- Wq = λ/(μ(μ-λ))

## 构建

```sh
go build -o des-queue .
```

## 用法

```sh
./des-queue -lambda 0.5 -mu 1.0 -servers 1 -customers 8000 -seed 42
```

参数：

| 参数 | 默认 | 说明 |
|------|------|------|
| `-lambda` | 0.5 | 到达率 λ |
| `-mu` | 1.0 | 服务率 μ（每服务台） |
| `-servers` | 1 | 服务台数量（1 为 M/M/1，>1 为 M/M/c） |
| `-customers` | 8000 | 仿真顾客数 |
| `-seed` | 42 | 随机种子（结果可复现） |

非法参数（λ/μ ≤ 0、servers < 1、customers < 1）或 M/M/1 不稳态（λ ≥ μ）时退出码 1 并输出明确错误信息；正常输出报告并退出码 0。

## 包结构

- `internal/event`：事件类型与基于 container/heap 的 FEL。
- `internal/sim`：配置校验、指数间隔序列器、M/M/1 与 M/M/c 仿真引擎。
- `internal/stats`：M/M/1 理论公式、误差对照与文本报告。

## 示例

见 `example/`。
