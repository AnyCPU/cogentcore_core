# ISSUE_005: Reflection Performance in Hot Paths

## Issue Summary

| Field | Value |
|-------|-------|
| **ID** | ISSUE_005 |
| **Title** | Uncached Reflection Operations in Performance-Critical Paths |
| **Severity** | Medium-High |
| **Category** | Performance |
| **File(s)** | `/home/user/cogentcore_core/base/reflectx/structs.go` |
| **Lines** | 30-64, 103-138, 152-192 |
| **Status** | Open |

---

## Description

The `base/reflectx` package provides reflection utilities that are used throughout the codebase for struct field walking, default value handling, and type introspection. These operations are called frequently (during styling, layout, and serialization) but do not cache the reflection metadata, leading to repeated expensive operations.

---

## Affected Code

### 1. WalkFields - Uncached Type Iteration

**File:** `/home/user/cogentcore_core/base/reflectx/structs.go`

```go
// Lines 30-51
func WalkFields(parent reflect.Value,
    should func(parent reflect.Value, field reflect.StructField, value reflect.Value) bool,
    walk func(parent reflect.Value, parentField *reflect.StructField, field reflect.StructField, value reflect.Value)) {
    walkFields(parent, nil, should, walk)
}

func walkFields(parent reflect.Value, parentField *reflect.StructField,
    should func(...) bool,
    walk func(...)) {
    typ := parent.Type()  // Called every time
    for i := 0; i < typ.NumField(); i++ {  // Iterates fields every call
        field := typ.Field(i)  // Allocates StructField
        if !field.IsExported() {
            continue
        }
        // ... rest of logic
    }
}
```

### 2. SetFromDefaultTags - Repeated Type Reflection

```go
// Lines 103-129
func SetFromDefaultTags(v any) error {
    ov := reflect.ValueOf(v)
    if IsNil(ov) {
        return nil
    }
    val := NonPointerValue(ov)
    typ := val.Type()  // Type() called every time
    for i := 0; i < typ.NumField(); i++ {  // NumField() every call
        f := typ.Field(i)  // Allocates StructField each iteration
        // ... process field
    }
    return nil
}
```

### 3. NonDefaultFields - Full Type Introspection

```go
// Lines 152-192
func NonDefaultFields(v any) map[string]any {
    res := map[string]any{}
    rv := Underlying(reflect.ValueOf(v))
    rt := rv.Type()  // Type lookup
    nf := rt.NumField()  // Field count
    for i := 0; i < nf; i++ {
        fv := rv.Field(i)
        ft := rt.Field(i)  // StructField allocation
        // ... process
    }
    return res
}
```

---

## Performance Analysis

### Benchmark Results (Estimated)

| Operation | Without Cache | With Cache | Improvement |
|-----------|--------------|------------|-------------|
| WalkFields (10 fields) | ~500ns | ~50ns | 10x |
| SetFromDefaultTags | ~800ns | ~100ns | 8x |
| NonDefaultFields | ~1200ns | ~200ns | 6x |
| Full style apply | ~5000ns | ~1000ns | 5x |

### Call Frequency Analysis

These functions are called in hot paths:

1. **Style Application**: Called for every widget during style updates
2. **Layout**: Called during size calculation
3. **Serialization**: Called during save/load operations
4. **Widget Creation**: Called during initialization

For a scene with 100 widgets, each style update could involve:
- ~500 calls to WalkFields
- ~200 calls to type introspection
- Accumulated overhead: significant

---

## Why Reflection is Expensive

### 1. Type Information Lookup

```go
typ := val.Type()  // Involves interface dispatch and type lookup
```

### 2. StructField Allocation

```go
field := typ.Field(i)  // Returns StructField by value (copy)
```

`reflect.StructField` is a relatively large struct:

```go
type StructField struct {
    Name      string
    PkgPath   string
    Type      Type
    Tag       StructTag  // string
    Offset    uintptr
    Index     []int      // slice allocation
    Anonymous bool
}
```

### 3. Repeated Computation

The same type information is computed repeatedly for the same types, even though it never changes.

---

## Recommended Fix

### Solution: Type Metadata Cache

Create a cache for type metadata to avoid repeated reflection:

```go
package reflectx

import (
    "reflect"
    "sync"
)

// TypeInfo caches reflection information for a struct type.
type TypeInfo struct {
    Type           reflect.Type
    Fields         []FieldInfo
    ExportedFields []FieldInfo  // Pre-filtered
    NumFields      int
}

// FieldInfo caches information about a struct field.
type FieldInfo struct {
    Index       int
    Name        string
    Type        reflect.Type
    Tag         reflect.StructTag
    Default     string            // Parsed default tag
    Anonymous   bool
    Offset      uintptr
}

// Global type cache
var (
    typeCache   sync.Map  // map[reflect.Type]*TypeInfo
    typeCacheMu sync.Mutex
)

// GetTypeInfo returns cached type information for the given type.
func GetTypeInfo(t reflect.Type) *TypeInfo {
    // Fast path: check cache
    if info, ok := typeCache.Load(t); ok {
        return info.(*TypeInfo)
    }

    // Slow path: build and cache
    typeCacheMu.Lock()
    defer typeCacheMu.Unlock()

    // Double-check after acquiring lock
    if info, ok := typeCache.Load(t); ok {
        return info.(*TypeInfo)
    }

    info := buildTypeInfo(t)
    typeCache.Store(t, info)
    return info
}

func buildTypeInfo(t reflect.Type) *TypeInfo {
    // Unwrap pointer types
    for t.Kind() == reflect.Ptr {
        t = t.Elem()
    }

    if t.Kind() != reflect.Struct {
        return &TypeInfo{Type: t}
    }

    nf := t.NumField()
    info := &TypeInfo{
        Type:      t,
        NumFields: nf,
        Fields:    make([]FieldInfo, nf),
    }

    for i := 0; i < nf; i++ {
        sf := t.Field(i)
        fi := FieldInfo{
            Index:     i,
            Name:      sf.Name,
            Type:      sf.Type,
            Tag:       sf.Tag,
            Default:   sf.Tag.Get("default"),
            Anonymous: sf.Anonymous,
            Offset:    sf.Offset,
        }
        info.Fields[i] = fi

        if sf.IsExported() {
            info.ExportedFields = append(info.ExportedFields, fi)
        }
    }

    return info
}
```

### Updated WalkFields

```go
// WalkFieldsCached uses cached type information for better performance.
func WalkFieldsCached(parent reflect.Value,
    should func(info *TypeInfo, fi *FieldInfo, value reflect.Value) bool,
    walk func(info *TypeInfo, fi *FieldInfo, value reflect.Value)) {

    info := GetTypeInfo(parent.Type())
    walkFieldsCached(parent, info, should, walk)
}

func walkFieldsCached(parent reflect.Value, info *TypeInfo,
    should func(*TypeInfo, *FieldInfo, reflect.Value) bool,
    walk func(*TypeInfo, *FieldInfo, reflect.Value)) {

    for i := range info.ExportedFields {
        fi := &info.ExportedFields[i]
        value := parent.Field(fi.Index)

        if !should(info, fi, value) {
            continue
        }

        if fi.Type.Kind() == reflect.Struct && fi.Anonymous {
            subInfo := GetTypeInfo(fi.Type)
            walkFieldsCached(value, subInfo, should, walk)
        } else {
            walk(info, fi, value)
        }
    }
}
```

### Updated SetFromDefaultTags

```go
func SetFromDefaultTagsCached(v any) error {
    ov := reflect.ValueOf(v)
    if IsNil(ov) {
        return nil
    }
    val := NonPointerValue(ov)
    info := GetTypeInfo(val.Type())

    for i := range info.ExportedFields {
        fi := &info.ExportedFields[i]
        fv := val.Field(fi.Index)

        if NonPointerType(fi.Type).Kind() == reflect.Struct && fi.Default == "" {
            SetFromDefaultTagsCached(PointerValue(fv).Interface())
            continue
        }

        if fi.Default != "" {
            err := SetFromDefaultTag(fv, fi.Default)
            if err != nil {
                return fmt.Errorf("error setting field %q: %w", fi.Name, err)
            }
        }
    }
    return nil
}
```

---

## Additional Optimizations

### 1. Pre-compute Default Values

```go
type FieldInfo struct {
    // ... existing fields ...

    // Pre-parsed default value
    HasDefault    bool
    DefaultParsed any  // nil if no default
}

func buildTypeInfo(t reflect.Type) *TypeInfo {
    // ... existing code ...

    for i := 0; i < nf; i++ {
        sf := t.Field(i)
        fi := FieldInfo{
            // ... existing fields ...
        }

        // Pre-parse default if present
        if def := sf.Tag.Get("default"); def != "" {
            fi.HasDefault = true
            // Could pre-parse common types here
        }

        info.Fields[i] = fi
    }
}
```

### 2. Interface Type Assertions

For commonly accessed interface types, pre-compute type assertions:

```go
type TypeInfo struct {
    // ... existing fields ...

    ImplementsShouldSaver bool
    ImplementsStringer    bool
}

func buildTypeInfo(t reflect.Type) *TypeInfo {
    // ... existing code ...

    shouldSaverType := reflect.TypeOf((*ShouldSaver)(nil)).Elem()
    info.ImplementsShouldSaver = t.Implements(shouldSaverType) ||
        reflect.PtrTo(t).Implements(shouldSaverType)
}
```

### 3. Field Index Lookup

```go
type TypeInfo struct {
    // ... existing fields ...

    FieldsByName map[string]int  // name -> index
}

func buildTypeInfo(t reflect.Type) *TypeInfo {
    // ... existing code ...

    info.FieldsByName = make(map[string]int, nf)
    for i, fi := range info.Fields {
        info.FieldsByName[fi.Name] = i
    }
}

// Fast field lookup
func (ti *TypeInfo) FieldByName(name string) (*FieldInfo, bool) {
    if idx, ok := ti.FieldsByName[name]; ok {
        return &ti.Fields[idx], true
    }
    return nil, false
}
```

---

## Benchmark Validation

Add benchmarks to verify improvement:

```go
func BenchmarkWalkFieldsUncached(b *testing.B) {
    type TestStruct struct {
        Field1 string `default:"test"`
        Field2 int    `default:"42"`
        Field3 bool   `default:"true"`
        Field4 float64
        Field5 string
        // ... more fields
    }

    v := reflect.ValueOf(TestStruct{})
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        WalkFields(v,
            func(parent reflect.Value, field reflect.StructField, value reflect.Value) bool {
                return true
            },
            func(parent reflect.Value, parentField *reflect.StructField, field reflect.StructField, value reflect.Value) {
                _ = field.Name
            })
    }
}

func BenchmarkWalkFieldsCached(b *testing.B) {
    type TestStruct struct {
        // Same as above
    }

    v := reflect.ValueOf(TestStruct{})
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        WalkFieldsCached(v,
            func(info *TypeInfo, fi *FieldInfo, value reflect.Value) bool {
                return true
            },
            func(info *TypeInfo, fi *FieldInfo, value reflect.Value) {
                _ = fi.Name
            })
    }
}
```

---

## Migration Strategy

### Phase 1: Add Cached Variants

```go
// Keep existing
func WalkFields(...)

// Add cached version
func WalkFieldsCached(...)
```

### Phase 2: Internal Migration

Update internal callers to use cached versions:

```go
// In styles package
func (s *Style) ApplyDefaults() {
    // Before
    reflectx.SetFromDefaultTags(s)

    // After
    reflectx.SetFromDefaultTagsCached(s)
}
```

### Phase 3: Deprecation

```go
// Deprecated: Use WalkFieldsCached for better performance.
func WalkFields(...)
```

---

## Testing Requirements

### 1. Correctness Test

```go
func TestCachedMatchesUncached(t *testing.T) {
    type TestStruct struct {
        Field1 string `default:"test"`
        Field2 int
    }

    v := reflect.ValueOf(TestStruct{})

    // Collect results from both versions
    var uncached, cached []string

    WalkFields(v, ...)  // collect to uncached
    WalkFieldsCached(v, ...)  // collect to cached

    assert.Equal(t, uncached, cached)
}
```

### 2. Benchmark Test

```go
func TestPerformanceImprovement(t *testing.T) {
    // Run benchmarks and verify improvement
    // At least 3x faster for cached version
}
```

### 3. Concurrent Access Test

```go
func TestCacheConcurrentAccess(t *testing.T) {
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            GetTypeInfo(reflect.TypeOf(TestStruct{}))
        }()
    }
    wg.Wait()
}
```

---

## Related Issues

- General performance optimization patterns
- Style application performance
- Layout calculation performance

---

## References

- [Go Blog: The Laws of Reflection](https://go.dev/blog/laws-of-reflection)
- [Reflection Performance](https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field)
- [sync.Map Documentation](https://pkg.go.dev/sync#Map)
- [Caching Reflection Results](https://medium.com/a-journey-with-go/go-how-to-take-advantage-of-the-caches-the-go-runtime-provides-e6f36f38c6b1)

---

*Created: 2025-11-22*
*Last Updated: 2025-11-22*
