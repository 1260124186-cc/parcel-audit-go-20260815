# 修复前 Base 状态复现

以下内容针对修复前的 base 状态。

## 可观察症状

调用方在审计开始前取消请求后，审计没有返回取消错误。计划读取仍会继续，
并可能产生普通审计结果。

## 触发命令

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t parcel-audit-004-base:arm64 .
docker run --rm --entrypoint go parcel-audit-004-base:arm64 test ./internal/store ./internal/service
```

## 实际完整输出

```text
--- FAIL: TestMemoryStoreLoadRespectsCancelledContext (0.00s)
    memory_test.go:35: Load() error = <nil>, want context cancellation
FAIL
FAIL    github.com/1260124186-cc/parcel-audit-go-20260815/internal/store  0.001s
--- FAIL: TestAuditRespectsCancelledContext (0.00s)
    audit_test.go:55: Audit() error = nil, want cancelled context error
FAIL
FAIL    github.com/1260124186-cc/parcel-audit-go-20260815/internal/service  0.001s
FAIL
```

## 期望行为

当调用方已取消请求时，审计应尽早停止并返回与取消状态一致的错误；不应继续
处理计划或返回普通审计结果。
