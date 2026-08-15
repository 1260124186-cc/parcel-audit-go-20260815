# 修复前复现

以下内容针对修复前的 base 状态。

## 可观察症状

当发货计划缺少必填的货件标识时，命令行没有将其提示为输入校验错误，
而是返回笼统的审计失败信息。

## 触发命令

```sh
docker run --rm --platform linux/arm64 --entrypoint go parcel-audit-003-base test ./cmd/parcel-audit
```

## 实际完整输出

```text
--- FAIL: TestRunReportsValidationError (0.00s)
    main_test.go:18: run() error = audit failed: validate plan: shipments: shipment id is required, want validation message
FAIL
FAIL	github.com/1260124186-cc/parcel-audit-go-20260815/cmd/parcel-audit	0.001s
FAIL
```

## 期望行为

无效的发货计划应向用户明确提示其未满足的输入校验规则，例如提示缺少货件标识，
而不是只显示通用的审计失败文本。
