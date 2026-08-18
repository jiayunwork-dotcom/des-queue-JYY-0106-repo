# 示例

与 README 一致的典型用法：

```sh
# M/M/1：lambda=0.5, mu=1.0, 8000 顾客, seed=42
./des-queue -lambda 0.5 -mu 1.0 -servers 1 -customers 8000 -seed 42

# M/M/c：3 个服务台
./des-queue -lambda 0.5 -mu 1.0 -servers 3 -customers 8000 -seed 42

# 不稳态输入（lambda >= mu，M/M/1）→ exit 1
./des-queue -lambda 2 -mu 1 -servers 1
```

参数示例文件 `args.txt` 列出了推荐参数组合，可用 `sh -c "./des-queue $(cat example/args.txt)"` 方式执行。
