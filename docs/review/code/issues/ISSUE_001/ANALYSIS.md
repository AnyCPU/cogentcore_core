# ISSUE_001: Timer Race Condition in Event Handling

## Issue Summary

| Field | Value |
|-------|-------|
| **ID** | ISSUE_001 |
| **Title** | Potential Race Condition and Deadlock Risk in Timer Management |
| **Severity** | Critical |
| **Category** | Concurrency |
| **File(s)** | `/home/user/cogentcore_core/core/events.go` |
| **Lines** | 610-706 |
| **Status** | Open |

---

## Description

The `Events.handleLong` function manages long hover and long press timers using a combination of `timerMu` mutex and render context locks. The current implementation has potential lock ordering issues that could lead to deadlocks, and there are race conditions between timer callbacks and the main event handling code.

---

## Affected Code

### Location: `/home/user/cogentcore_core/core/events.go`

```go
// Lines 610-706
func (em *Events) handleLong(e events.Event, deep Widget, w *Widget, pos *image.Point,
    t **time.Timer, styp, etyp events.Types, stime time.Duration, sdist int) {

    em.timerMu.Lock()      // LOCK #1: timerMu acquired first
    defer em.timerMu.Unlock()

    // ... logic ...

    *t = time.AfterFunc(stime, func() {
        win := em.RenderWindow()
        if win == nil {
            return
        }
        rc := win.renderContext()
        rc.Lock()              // LOCK #2: rc acquired second
        defer rc.Unlock()

        em.timerMu.Lock()      // LOCK #3: timerMu acquired third (re-acquisition)
        defer em.timerMu.Unlock()

        if tree.IsNil(*w) {
            return
        }
        (*w).AsWidget().Send(styp, e)
        *t = nil
    })
}
```

---

## Problem Analysis

### 1. Lock Ordering Inconsistency

The code exhibits inconsistent lock ordering between the main path and the timer callback:

**Main Event Path (implicit from callers):**
```
rc.Lock() -> timerMu.Lock()
```

**Timer Callback:**
```
rc.Lock() -> timerMu.Lock()
```

While this specific ordering appears consistent, the problem arises because:
- The `handleLong` function acquires `timerMu` first
- But the timer callback acquires `rc` first, then `timerMu`

If another code path holds `timerMu` and waits for `rc`, while the timer callback holds `rc` and waits for `timerMu`, a deadlock occurs.

### 2. Race Condition on Widget Reference

The `*w` pointer is accessed both in the main goroutine and the timer callback without proper synchronization during the window between when the timer fires and when the lock is acquired:

```go
// Timer callback
rc.Lock()
defer rc.Unlock()
em.timerMu.Lock()          // Widget could change here!
defer em.timerMu.Unlock()
if tree.IsNil(*w) {        // *w may have been modified
    return
}
```

### 3. Stale Event Reference

The event `e` is captured by the timer closure but may become stale or invalid by the time the timer fires:

```go
*t = time.AfterFunc(stime, func() {
    // 'e' was captured at timer creation time
    // but may be from a pool that's been recycled
    (*w).AsWidget().Send(styp, e)  // Using potentially stale event
})
```

---

## Impact Assessment

| Impact Type | Severity | Description |
|-------------|----------|-------------|
| **Deadlock** | Critical | Application can freeze completely |
| **Race Condition** | High | Unpredictable behavior, crashes |
| **Memory Safety** | Medium | Stale pointer access possible |
| **User Experience** | High | Long hover/press features may fail |

---

## Reproduction Steps

1. Rapidly move mouse over hoverable elements
2. Trigger multiple long hover/press candidates in quick succession
3. Under load, the deadlock may manifest as application freeze
4. Race conditions may cause intermittent crashes or missed events

---

## Recommended Fix

### Option 1: Consistent Lock Ordering (Preferred)

Establish a strict lock hierarchy: always acquire `rc` before `timerMu`.

```go
func (em *Events) handleLong(e events.Event, deep Widget, w *Widget, pos *image.Point,
    t **time.Timer, styp, etyp events.Types, stime time.Duration, sdist int) {

    // Get render context first
    win := em.RenderWindow()
    if win == nil {
        return
    }
    rc := win.renderContext()
    if rc == nil {
        return
    }

    rc.Lock()
    defer rc.Unlock()
    em.timerMu.Lock()
    defer em.timerMu.Unlock()

    // ... rest of logic with consistent ordering ...
}
```

### Option 2: Copy Required State

Capture all necessary state before timer creation to avoid lock acquisition in callback:

```go
*t = time.AfterFunc(stime, func() {
    // Capture widget reference at creation time
    widget := *w
    if tree.IsNil(widget) {
        return
    }

    win := em.RenderWindow()
    if win == nil {
        return
    }
    rc := win.renderContext()
    rc.Lock()
    defer rc.Unlock()

    // Now safe to proceed
    widget.AsWidget().Send(styp, events.NewEvent(styp))  // Create fresh event
})
```

### Option 3: Use Channel-Based Coordination

Replace timer callbacks with channel-based event delivery:

```go
type timerEvent struct {
    widget Widget
    typ    events.Types
}

// In event loop
select {
case te := <-em.timerChan:
    te.widget.AsWidget().Send(te.typ, nil)
case <-otherChannels:
    // ...
}
```

---

## Testing Requirements

1. **Race Detector Test:**
```bash
go test -race ./core/... -count=100
```

2. **Stress Test:**
```go
func TestHandleLongRace(t *testing.T) {
    // Create multiple rapid hover/unhover cycles
    // Verify no deadlocks or races
}
```

3. **Deadlock Detection:**
```go
func TestNoDeadlock(t *testing.T) {
    done := make(chan bool)
    go func() {
        // Trigger long hover scenarios
        done <- true
    }()
    select {
    case <-done:
        // Success
    case <-time.After(5 * time.Second):
        t.Fatal("Potential deadlock detected")
    }
}
```

---

## Related Issues

- Lock ordering in `AsyncLock` (ISSUE_002)
- Timer management patterns across codebase
- Event pooling considerations

---

## References

- [Go Data Race Detector](https://go.dev/doc/articles/race_detector)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Memory Model](https://go.dev/ref/mem)

---

*Created: 2025-11-22*
*Last Updated: 2025-11-22*
