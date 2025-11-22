# ISSUE_002: Goroutine Leak Risk in AsyncLock

## Issue Summary

| Field | Value |
|-------|-------|
| **ID** | ISSUE_002 |
| **Title** | Infinite Blocking Patterns in AsyncLock Can Cause Goroutine Leaks |
| **Severity** | High |
| **Category** | Concurrency / Resource Management |
| **File(s)** | `/home/user/cogentcore_core/core/render.go` |
| **Lines** | 38-70 |
| **Status** | Open |

---

## Description

The `AsyncLock` method in `WidgetBase` uses infinite `select {}` statements to block goroutines when certain error conditions are detected (deleted widget, unavailable render context). While this prevents crashes, it creates permanent goroutine leaks that cannot be cleaned up, leading to memory leaks and resource exhaustion over time.

---

## Affected Code

### Location: `/home/user/cogentcore_core/core/render.go`

```go
// Lines 38-70
func (wb *WidgetBase) AsyncLock() {
    rc := wb.Scene.renderContext()
    if rc == nil {
        if wb.Scene.hasFlag(sceneHasShown) {
            // If the scene has been shown but there is no render context,
            // we are probably being deleted, so we just block forever.
            if DebugSettings.UpdateTrace {
                fmt.Println("AsyncLock: scene shown but no render context; blocking forever:", wb)
            }
            select {}  // INFINITE BLOCK - Goroutine leak!
        }
        // Otherwise, if we haven't been shown yet, we just wait until we are
        // and then try again.
        if DebugSettings.UpdateTrace {
            fmt.Println("AsyncLock: waiting for scene to be shown:", wb)
        }
        onShow := make(chan struct{})
        wb.OnShow(func(e events.Event) {
            onShow <- struct{}{}
        })
        <-onShow
        wb.AsyncLock() // try again
        return
    }
    rc.Lock()
    if wb.This == nil {
        rc.Unlock()
        if DebugSettings.UpdateTrace {
            fmt.Println("AsyncLock: widget deleted; blocking forever:", wb)
        }
        select {}  // INFINITE BLOCK - Goroutine leak!
    }
    wb.Scene.setFlag(true, sceneUpdating)
}
```

---

## Problem Analysis

### 1. Permanent Goroutine Leak

When `select {}` is used without any case statements, the goroutine blocks forever and can never be garbage collected:

```go
// This goroutine will NEVER terminate
go func() {
    widget.AsyncLock()  // If widget is deleted, blocks forever
    // Code after this never executes
    widget.Update()
    widget.AsyncUnlock()
}()
```

### 2. No Cancellation Mechanism

There's no way for:
- The caller to cancel the operation
- The application to clean up during shutdown
- The framework to detect and recover from this state

### 3. Memory Accumulation

Each leaked goroutine holds references to:
- The widget being accessed
- The scene
- The closure's captured variables
- Stack memory (~2KB minimum per goroutine)

### 4. Debug-Only Visibility

The leak is only visible when `DebugSettings.UpdateTrace` is enabled:

```go
if DebugSettings.UpdateTrace {
    fmt.Println("AsyncLock: widget deleted; blocking forever:", wb)
}
select {}  // Silent leak in production
```

---

## Impact Assessment

| Impact Type | Severity | Description |
|-------------|----------|-------------|
| **Memory Leak** | High | Unbounded goroutine accumulation |
| **Resource Exhaustion** | High | Can eventually exhaust system resources |
| **Debugging Difficulty** | Medium | Silent failure, hard to diagnose |
| **Application Stability** | High | Long-running apps will degrade |

### Scenario Analysis

| Scenario | Frequency | Impact |
|----------|-----------|--------|
| Widget deleted during async op | Medium | Goroutine leak |
| Scene destroyed during async op | Low | Goroutine leak |
| Rapid widget creation/deletion | High | Multiple leaks |
| Application shutdown | Always | Delayed/incomplete |

---

## Reproduction Steps

```go
func LeakDemo() {
    // Create a scene with a widget
    scene := core.NewScene()
    widget := core.NewButton(scene)

    // Start async operation
    go func() {
        widget.AsyncLock()
        defer widget.AsyncUnlock()
        // Do async work
    }()

    // Delete widget before async operation completes
    widget.Delete()

    // The goroutine is now leaked forever
    // runtime.NumGoroutine() will increase permanently
}
```

---

## Recommended Fix

### Option 1: Context-Based Cancellation (Preferred)

```go
func (wb *WidgetBase) AsyncLockWithContext(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        rc := wb.Scene.renderContext()
        if rc == nil {
            if wb.Scene.hasFlag(sceneHasShown) {
                // Scene shown but no render context - being deleted
                return errors.New("widget is being deleted")
            }

            // Wait for scene to be shown with timeout
            onShow := make(chan struct{}, 1)
            wb.OnShow(func(e events.Event) {
                select {
                case onShow <- struct{}{}:
                default:
                }
            })

            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-onShow:
                continue // Retry
            case <-time.After(30 * time.Second):
                return errors.New("timeout waiting for scene")
            }
        }

        rc.Lock()
        if wb.This == nil {
            rc.Unlock()
            return errors.New("widget deleted")
        }
        wb.Scene.setFlag(true, sceneUpdating)
        return nil
    }
}

// Backward-compatible wrapper
func (wb *WidgetBase) AsyncLock() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := wb.AsyncLockWithContext(ctx); err != nil {
        if DebugSettings.UpdateTrace {
            fmt.Println("AsyncLock failed:", err, wb)
        }
        // Caller should handle error, but for compatibility we could panic
        // or use a global error handler
    }
}
```

### Option 2: Return Error Instead of Blocking

```go
func (wb *WidgetBase) TryAsyncLock() (bool, error) {
    rc := wb.Scene.renderContext()
    if rc == nil {
        if wb.Scene.hasFlag(sceneHasShown) {
            return false, errors.New("render context unavailable, scene being deleted")
        }
        return false, errors.New("scene not yet shown")
    }

    rc.Lock()
    if wb.This == nil {
        rc.Unlock()
        return false, errors.New("widget deleted")
    }
    wb.Scene.setFlag(true, sceneUpdating)
    return true, nil
}
```

### Option 3: Timeout-Based Blocking

```go
func (wb *WidgetBase) AsyncLockWithTimeout(timeout time.Duration) error {
    deadline := time.Now().Add(timeout)

    for time.Now().Before(deadline) {
        rc := wb.Scene.renderContext()
        if rc == nil {
            if wb.Scene.hasFlag(sceneHasShown) {
                return errors.New("widget being deleted")
            }
            time.Sleep(10 * time.Millisecond)
            continue
        }

        rc.Lock()
        if wb.This == nil {
            rc.Unlock()
            return errors.New("widget deleted")
        }
        wb.Scene.setFlag(true, sceneUpdating)
        return nil
    }
    return errors.New("timeout acquiring async lock")
}
```

---

## Migration Strategy

1. **Phase 1:** Add new context-aware method alongside existing
2. **Phase 2:** Update internal callers to use new method
3. **Phase 3:** Mark old method as deprecated
4. **Phase 4:** Update documentation with best practices

### Example Usage After Fix:

```go
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := widget.AsyncLockWithContext(ctx); err != nil {
        log.Printf("Failed to acquire lock: %v", err)
        return
    }
    defer widget.AsyncUnlock()

    // Safe to update widget
    widget.Update()
}()
```

---

## Testing Requirements

### 1. Leak Detection Test

```go
func TestAsyncLockNoLeak(t *testing.T) {
    initialGoroutines := runtime.NumGoroutine()

    // Create and destroy many widgets with async operations
    for i := 0; i < 100; i++ {
        scene := core.NewScene()
        widget := core.NewButton(scene)

        done := make(chan bool)
        go func() {
            ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
            defer cancel()
            widget.AsyncLockWithContext(ctx)
            done <- true
        }()

        widget.Delete()
        <-done
    }

    // Force GC
    runtime.GC()
    time.Sleep(100 * time.Millisecond)

    finalGoroutines := runtime.NumGoroutine()
    leaked := finalGoroutines - initialGoroutines

    if leaked > 5 { // Allow small variance
        t.Errorf("Goroutine leak detected: started=%d, ended=%d, leaked=%d",
            initialGoroutines, finalGoroutines, leaked)
    }
}
```

### 2. Cancellation Test

```go
func TestAsyncLockCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    scene := core.NewScene()
    widget := core.NewButton(scene)
    // Don't show scene - renderContext will be nil

    errCh := make(chan error)
    go func() {
        errCh <- widget.AsyncLockWithContext(ctx)
    }()

    // Cancel after short delay
    time.Sleep(50 * time.Millisecond)
    cancel()

    select {
    case err := <-errCh:
        if !errors.Is(err, context.Canceled) {
            t.Errorf("Expected context.Canceled, got: %v", err)
        }
    case <-time.After(1 * time.Second):
        t.Error("Cancellation did not work - goroutine may be leaked")
    }
}
```

---

## Related Issues

- ISSUE_001: Timer race conditions also use complex locking
- General goroutine lifecycle patterns in the codebase
- Shutdown/cleanup procedures

---

## References

- [Go Blog: Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Go Blog: Context](https://go.dev/blog/context)
- [Goroutine Leaks - The Forgotten Sender](https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html)

---

*Created: 2025-11-22*
*Last Updated: 2025-11-22*
