# Project-1: 归并排序

### 二路归并排序

#### Bench 数据对比

cmd: `go test -bench .`

- BenchmarkNormalSort: 3683251390 ns/op
- BenchmarkMergeSort:  3957939179 ns/op

### 内存和 CPU 分析

cmd: `go test -bench . -version v1`

#### CPU

| sort method | Flat  |  Flat% |  Sum%  |  Cum  |  Cum%  |      Name        |
| ----------- | ----- | ------ | ------ | ----- | ------ | ---------------- |
|    merge    | 1.54s | 44.13% | 44.13% | 2.52s | 72.21% | mergesort.merge  |

- **分析**：
    1. merge 判断逻辑过于复杂
    ```go
    28   310ms      310ms   for i, j, k = 0, 0, 0; (j < lb) || (k < lc); { 
	29   230ms      230ms   if (j < lb) && (!(k < lc) || (b[j] <= c[k])) {
	34   370ms      370ms   if (k < lc) && (!(j < lb) || c[k] < b[j]) {
    ```
    `这个地方虽然形式上看起来工整简单，但实际上的逻辑判断很复杂。`

    2. 临时变量 b 的初始化和分配调用耗时
    ```go
    21       .      940ms    b := make([]int64, lb)
    ```
    `这个地方每次调用都会创建和分配临时变量。`

- **优化**：
    1. 根据二路归并合并的几种情况，简化逻辑为：
    2. 提前分配公用的空间，分配操作次数降为O(1)。

#### Mem
- **分析**：
    1. 临时变量 b 的大量分配
    ```go
    21    1.56GB     1.56GB  b := make([]int64, lb) 
    ```
- **优化**：
    1. 一次性分配足够大的内存空间。



