# ISSUE_004: Missing Context Propagation in Long-Running Operations

## Issue Summary

| Field | Value |
|-------|-------|
| **ID** | ISSUE_004 |
| **Title** | Missing context.Context Propagation for Cancellation and Timeout Support |
| **Severity** | High |
| **Category** | API Design / Concurrency |
| **File(s)** | Multiple async paths throughout codebase |
| **Status** | Open |

---

## Description

Many long-running and asynchronous operations in the Cogent Core framework do not accept or propagate `context.Context`. This prevents:
1. Graceful cancellation of operations
2. Timeout enforcement
3. Proper resource cleanup during shutdown
4. Integration with Go's standard cancellation patterns

---

## Affected Areas

### 1. Async Widget Operations

**File:** `/home/user/cogentcore_core/core/render.go`

```go
// Current API - no context support
func (wb *WidgetBase) AsyncLock()
func (wb *WidgetBase) AsyncUnlock()
```

### 2. File Tree Operations

**File:** `/home/user/cogentcore_core/filetree/tree.go`

```go
// Background directory reading - no cancellation
go func() {
    // Long-running file system traversal
    // Cannot be cancelled
}()
```

### 3. Network Operations

**File:** `/home/user/cogentcore_core/htmlcore/handler.go`

```go
// HTTP fetching without context
go func() {
    // Network request
    // No timeout or cancellation
}()
```

### 4. Video Processing

**File:** `/home/user/cogentcore_core/video/video.go`

```go
// Video decoding loops - no cancellation support
go func() {
    for {
        // Decode frames
        // No way to stop cleanly
    }
}()
```

### 5. GPU Compute Operations

**File:** `/home/user/cogentcore_core/gpu/compute.go`

```go
// Compute dispatching without context
go func() {
    // GPU operations
    // Cannot be cancelled or timed out
}()
```

---

## Problem Analysis

### 1. No Graceful Shutdown

When the application exits or a window closes, background goroutines cannot be cleanly terminated:

```go
// Current pattern
func StartBackgroundTask() {
    go func() {
        for {
            // Process forever
            // No way to stop
        }
    }()
}

// When application closes:
// - Goroutine continues running
// - Resources not released
// - Potential crashes on cleanup
```

### 2. No Timeout Support

Long operations can hang indefinitely:

```go
// Current pattern
func FetchData() {
    go func() {
        resp, _ := http.Get(url)  // No timeout
        // Could block forever on slow network
    }()
}
```

### 3. No Deadline Propagation

When a parent operation has a deadline, child operations are unaware:

```go
// Parent has 5 second timeout
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// But child operations ignore it
widget.AsyncLock()  // No ctx parameter
doExpensiveWork()   // Ignores parent timeout
```

---

## Impact Assessment

| Impact Type | Severity | Description |
|-------------|----------|-------------|
| **Resource Leaks** | High | Goroutines and resources not cleaned up |
| **Unresponsive Shutdown** | High | App hangs on close |
| **Poor UX** | Medium | Cannot cancel slow operations |
| **Testing Difficulty** | Medium | Hard to test timeout behavior |
| **Composability** | High | Cannot integrate with context-aware code |

---

## Go Standard Patterns

The standard Go pattern for long-running operations:

```go
// Function signature includes context
func DoWork(ctx context.Context, args ...any) error {
    // Check for cancellation at key points
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Do work, periodically checking context
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case result := <-workChan:
            process(result)
        }
    }
}
```

---

## Recommended Fix

### Phase 1: Core Async Operations

**1. AsyncLock with Context:**

```go
// New API
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
                return ErrWidgetDeleted
            }
            // Wait with context
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-wb.Scene.showChan:
                continue
            }
        }

        if rc.TryLock() {
            if wb.This == nil {
                rc.Unlock()
                return ErrWidgetDeleted
            }
            wb.Scene.setFlag(true, sceneUpdating)
            return nil
        }

        // Brief sleep before retry, with context check
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(time.Millisecond):
        }
    }
}

// Backward compatibility
func (wb *WidgetBase) AsyncLock() {
    ctx := context.Background()
    if err := wb.AsyncLockWithContext(ctx); err != nil {
        // Handle error - log or panic depending on policy
    }
}
```

### Phase 2: File Operations

```go
// Current
func (ft *Tree) ReadDir(path string)

// Proposed
func (ft *Tree) ReadDirWithContext(ctx context.Context, path string) error {
    entries, err := os.ReadDir(path)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        // Process entry
    }
    return nil
}
```

### Phase 3: Network Operations

```go
// Current
func fetchURL(url string) ([]byte, error)

// Proposed
func fetchURL(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

### Phase 4: Scene/Stage Lifecycle

```go
type Scene struct {
    // Add context for lifecycle management
    ctx    context.Context
    cancel context.CancelFunc
}

func NewScene() *Scene {
    ctx, cancel := context.WithCancel(context.Background())
    return &Scene{
        ctx:    ctx,
        cancel: cancel,
    }
}

func (sc *Scene) Close() {
    sc.cancel()  // Signals all background work to stop
}
```

---

## Implementation Guidelines

### 1. Context Placement Convention

Following Go conventions, context should be first parameter:

```go
// Good
func DoWork(ctx context.Context, arg1 string, arg2 int) error

// Bad
func DoWork(arg1 string, ctx context.Context, arg2 int) error
```

### 2. Context Checking Points

Check context at:
- Start of function
- Before expensive operations
- In loops (at least once per iteration)
- Before I/O operations

```go
func ProcessItems(ctx context.Context, items []Item) error {
    // Check at start
    if err := ctx.Err(); err != nil {
        return err
    }

    for i, item := range items {
        // Check in loop
        select {
        case <-ctx.Done():
            return fmt.Errorf("cancelled after %d items: %w", i, ctx.Err())
        default:
        }

        // Expensive operation
        if err := process(ctx, item); err != nil {
            return err
        }
    }
    return nil
}
```

### 3. Background Goroutine Pattern

```go
func (s *Service) Start(ctx context.Context) {
    go func() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                // Clean shutdown
                return
            case <-ticker.C:
                s.doWork()
            }
        }
    }()
}
```

---

## Migration Strategy

### Phase 1: Add Context-Aware Variants

Keep existing APIs, add new versions with context:

```go
// Existing (keep for compatibility)
func (wb *WidgetBase) AsyncLock()

// New
func (wb *WidgetBase) AsyncLockWithContext(ctx context.Context) error
```

### Phase 2: Internal Migration

Update internal code to use context-aware versions:

```go
// Before
go func() {
    widget.AsyncLock()
    defer widget.AsyncUnlock()
    // work
}()

// After
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := widget.AsyncLockWithContext(ctx); err != nil {
        log.Printf("Failed to lock: %v", err)
        return
    }
    defer widget.AsyncUnlock()
    // work
}()
```

### Phase 3: Deprecation

Mark old APIs as deprecated with clear migration path:

```go
// Deprecated: Use AsyncLockWithContext instead for cancellation support.
func (wb *WidgetBase) AsyncLock()
```

### Phase 4: Documentation

Update documentation with context usage examples and best practices.

---

## Testing Requirements

### 1. Cancellation Test

```go
func TestContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    started := make(chan struct{})
    finished := make(chan error)

    go func() {
        close(started)
        finished <- doLongOperation(ctx)
    }()

    <-started
    cancel()

    select {
    case err := <-finished:
        if !errors.Is(err, context.Canceled) {
            t.Errorf("Expected context.Canceled, got: %v", err)
        }
    case <-time.After(time.Second):
        t.Error("Operation did not cancel in time")
    }
}
```

### 2. Timeout Test

```go
func TestContextTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    start := time.Now()
    err := doLongOperation(ctx)
    elapsed := time.Since(start)

    if !errors.Is(err, context.DeadlineExceeded) {
        t.Errorf("Expected deadline exceeded, got: %v", err)
    }
    if elapsed > 200*time.Millisecond {
        t.Errorf("Took too long to cancel: %v", elapsed)
    }
}
```

---

## Related Issues

- ISSUE_002: AsyncLock infinite blocking
- ISSUE_001: Timer management
- General goroutine lifecycle patterns

---

## References

- [Go Blog: Context](https://go.dev/blog/context)
- [Go Package: context](https://pkg.go.dev/context)
- [Context Best Practices](https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go)
- [Context and Structs](https://go.dev/blog/context-and-structs)

---

*Created: 2025-11-22*
*Last Updated: 2025-11-22*
