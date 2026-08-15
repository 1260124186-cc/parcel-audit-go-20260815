# 修复前基线复现

以下内容针对修复前的 `base_bug_002` 状态。

## 可观察症状

运行完整测试时，标签规范化会改变传入标签内容，且内存计划存储会观察到调用方后续对原计划的修改。

## 触发命令

```sh
docker run --rm parcel-audit-002-base:arm64 go test ./...
```

## 实际完整错误输出

```text
ok  	github.com/1260124186-cc/parcel-audit-go-20260815/cmd/parcel-audit	0.001s
--- FAIL: TestNormalizeLabelsDoesNotMutateInput (0.00s)
    labels_test.go:16: input labels changed = []string{"cold", "fragile", "fragile"}, want []string{" Cold ", "fragile", "cold"}
FAIL
FAIL	github.com/1260124186-cc/parcel-audit-go-20260815/internal/domain	0.001s
ok  	github.com/1260124186-cc/parcel-audit-go-20260815/internal/service	0.001s
--- FAIL: TestMemoryStoreOwnsPlanSnapshot (0.00s)
    memory_test.go:24: stored labels = []string{"changed"}, want []string{"cold"}
FAIL
FAIL	github.com/1260124186-cc/parcel-audit-go-20260815/internal/store	0.001s
?   	github.com/1260124186-cc/parcel-audit-go-20260815/internal/transport	[no test files]
FAIL
```

## 期望行为

标签规范化后应返回规范化结果，但传入标签内容保持不变。计划交给内存存储后，调用方对原计划的后续修改不应影响后续读取和审计。
