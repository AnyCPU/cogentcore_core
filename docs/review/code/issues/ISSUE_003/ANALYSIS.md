# ISSUE_003: Inconsistent Mutex Embedding Patterns

## Issue Summary

| Field | Value |
|-------|-------|
| **ID** | ISSUE_003 |
| **Title** | Inconsistent Mutex Usage Patterns Across Packages |
| **Severity** | High |
| **Category** | Code Quality / Concurrency |
| **File(s)** | Multiple (50+ mutex instances) |
| **Status** | Open |

---

## Description

The codebase uses over 50 instances of `sync.Mutex` and `sync.RWMutex` with inconsistent patterns:
- Some structs embed the mutex directly
- Some use named fields with varying names (`mu`, `Mu`, `timerMu`, etc.)
- Some export the mutex, others don't
- Documentation of protected fields varies

This inconsistency makes it harder to:
1. Understand what's protected by each mutex
2. Maintain consistent locking patterns
3. Audit for race conditions
4. Onboard new developers

---

## Pattern Analysis

### Pattern 1: Direct Embedding (Anonymous)

**Used in:** `gpu/surface.go`, `core/stages.go`, `core/sprite.go`

```go
// gpu/surface.go
type Surface struct {
    sync.Mutex  // Embedded anonymously
    Format      TextureFormat
    // ...
}

// Usage
s.Lock()
defer s.Unlock()
```

**Pros:**
- Clean, idiomatic Go
- Methods directly available on struct

**Cons:**
- Lock() is exported even if struct is internal
- Unclear what fields are protected

### Pattern 2: Named Exported Field

**Used in:** `undo/undo.go`, `events/deque.go`, `text/lines/undo.go`

```go
// undo/undo.go
type Stack[T any] struct {
    Mu    sync.Mutex  // Exported
    // ...
}

// events/deque.go
type Deque[T any] struct {
    Mu   sync.Mutex
    // ...
}
```

**Pros:**
- Explicit naming
- Can be accessed externally if needed

**Cons:**
- Inconsistent naming (Mu vs mu)
- Exported mutex can be misused

### Pattern 3: Named Unexported Field

**Used in:** `core/events.go`, `core/windowgeometry.go`, `core/textfield.go`

```go
// core/events.go
type Events struct {
    timerMu sync.Mutex  // Specific purpose in name
    // ...
}

// core/windowgeometry.go
type windowGeometrySaver struct {
    mu sync.RWMutex
    // ...
}

// core/textfield.go
type TextField struct {
    cursorMu sync.Mutex
    // ...
}
```

**Pros:**
- Purpose-specific naming
- Internal implementation detail

**Cons:**
- Naming inconsistency (timerMu, cursorMu, mu)

### Pattern 4: Multiple Mutexes in One Struct

**Used in:** `text/parse/filestates.go`, `core/completer.go`

```go
// text/parse/filestates.go
type FileStates struct {
    SwitchMu sync.Mutex  // One mutex
    ProcMu   sync.Mutex  // Another mutex
    // ...
}

// core/completer.go
type completer struct {
    delayMu    sync.Mutex
    showMu     sync.Mutex
    // ...
}
```

---

## Inventory of Current Usage

| Package | File | Field Name | Type | Exported |
|---------|------|------------|------|----------|
| undo | undo.go | Mu | sync.Mutex | Yes |
| gpu/gpudraw | drawer.go | (embedded) | sync.Mutex | N/A |
| core | events.go | timerMu | sync.Mutex | No |
| core | tabs.go | mu | sync.Mutex | No |
| events | deque.go | Mu | sync.Mutex | Yes |
| core | windowgeometry.go | mu | sync.RWMutex | No |
| gpu | surface.go | (embedded) | sync.Mutex | N/A |
| core | stages.go | (embedded) | sync.Mutex | N/A |
| core | sprite.go | (embedded) | sync.Mutex | N/A |
| core | renderwindow.go | (global) | sync.Mutex | No |
| core | textfield.go | cursorMu | sync.Mutex | No |
| paint/renderers/rasterx | glyphcache.go | (embedded) | sync.Mutex | N/A |
| icons | icons.go | usedMu | sync.Mutex | No |
| filetree | dir.go | (embedded) | sync.Mutex | N/A |
| svg | svg.go | shaperMu | sync.Mutex | No |
| text/lines | lines.go | markupDelayMu | sync.Mutex | No |
| text/lines | lines.go | (embedded) | sync.Mutex | N/A |
| text/spell | model.go | (embedded) | sync.RWMutex | N/A |
| text/spell | spell.go | mu | sync.RWMutex | No |

---

## Problems Identified

### 1. No Clear Guideline

There's no documented standard for choosing between patterns.

### 2. Documentation Gap

Most mutexes lack comments explaining:
- What fields they protect
- Lock ordering requirements
- When locking is required

**Bad Example:**
```go
type Lines struct {
    sync.Mutex  // Protects... everything? Something?
    // 15 fields follow
}
```

**Good Example (rare in codebase):**
```go
type Events struct {
    // mutex that protects timer variable updates (e.g., hover AfterFunc's).
    timerMu sync.Mutex
    // ...
}
```

### 3. Naming Inconsistency

| Name Used | Count | Context |
|-----------|-------|---------|
| `Mu` | 8 | Various |
| `mu` | 12 | Various |
| `*Mu` (suffix) | 15 | Purpose-specific |
| (embedded) | 15 | Anonymous |

### 4. Export Inconsistency

Some exported structs have unexported mutexes, some have exported. The rationale isn't clear.

---

## Recommended Standards

### Standard 1: Choose Based on Scope

```go
// For internal implementation - use unexported field
type internal struct {
    mu sync.Mutex
    // ...
}

// For types where external synchronization may be needed
type Exported struct {
    // Mu protects concurrent access to Value.
    // Callers must hold Mu before reading or writing Value.
    Mu    sync.Mutex
    Value int
}
```

### Standard 2: Document Protected Fields

```go
type SafeBuffer struct {
    // mu protects buf and offset. All methods that read or write
    // these fields must hold mu.
    mu     sync.Mutex
    buf    []byte   // protected by mu
    offset int      // protected by mu
    name   string   // immutable, no protection needed
}
```

### Standard 3: Purpose-Specific Names for Multiple Mutexes

```go
type ComplexType struct {
    // dataMu protects data and dataVersion.
    dataMu      sync.Mutex
    data        []byte
    dataVersion int

    // stateMu protects state and lastUpdate.
    // Lock ordering: always acquire dataMu before stateMu if both needed.
    stateMu    sync.Mutex
    state      State
    lastUpdate time.Time
}
```

### Standard 4: Avoid Exported Anonymous Embedding

```go
// BAD: Lock() is exported, unclear what's protected
type Widget struct {
    sync.Mutex
    // ...
}

// GOOD: Internal detail, clear documentation
type Widget struct {
    // mu protects all mutable fields.
    mu sync.Mutex
    // ...
}
```

---

## Recommended Fix

### Phase 1: Document Existing Mutexes

Add comments to all existing mutex fields explaining what they protect:

```go
// Before
type Lines struct {
    sync.Mutex
    // ...fields...
}

// After
type Lines struct {
    // Lines embeds a Mutex that protects all mutable fields.
    // All public methods acquire this lock; callers should not
    // need to lock externally.
    sync.Mutex
    // ...fields...
}
```

### Phase 2: Standardize Naming

Choose a convention and apply consistently:
- Embedded for single-mutex types with full struct protection
- Named `mu` for internal single mutex
- Named `{purpose}Mu` for specific-purpose mutexes

### Phase 3: Add Lock Ordering Documentation

For packages with multiple related mutexes:

```go
// LOCK ORDERING for core package:
// When acquiring multiple locks, always use this order:
// 1. renderContext.Lock()
// 2. scene.Mutex
// 3. events.timerMu
// 4. widget-specific locks
//
// Never hold a lower-numbered lock while acquiring a higher-numbered one.
```

---

## Migration Example

### Before:

```go
// text/lines/lines.go
type Lines struct {
    markupDelayMu sync.Mutex
    sync.Mutex
    // 20+ fields, unclear what each mutex protects
}
```

### After:

```go
// text/lines/lines.go
type Lines struct {
    // Lines embeds sync.Mutex for primary synchronization.
    // This mutex protects all line content, markup, and view state.
    // All public methods acquire this lock automatically.
    // See package documentation for full concurrency contract.
    sync.Mutex

    // markupDelayMu protects only the markup delay timer state.
    // This is separate from the main mutex to allow markup delay
    // changes without blocking line operations.
    // Lock ordering: acquire Lines.Mutex before markupDelayMu.
    markupDelayMu sync.Mutex
    markupDelay   *time.Timer // protected by markupDelayMu

    // Below fields protected by embedded Mutex
    lines  []Line
    markup []Markup
    // ...
}
```

---

## Testing Requirements

1. **Race Detector Coverage:**
```bash
go test -race ./... -count=10
```

2. **Lock Order Verification:**
```go
// Consider using golang.org/x/sync/semaphore for testing
// or custom lock order checking in debug builds
```

---

## Related Issues

- ISSUE_001: Timer race conditions involve timerMu
- ISSUE_002: AsyncLock uses render context locking
- General concurrency patterns in the codebase

---

## References

- [Go Wiki: Mutex or Channel](https://github.com/golang/go/wiki/MutexOrChannel)
- [Effective Go: Sharing by Communicating](https://go.dev/doc/effective_go#sharing)
- [Go Concurrency Patterns](https://go.dev/talks/2012/concurrency.slide)

---

*Created: 2025-11-22*
*Last Updated: 2025-11-22*
