# Cogent Core Architecture Diagrams

This document provides visual representations of the Cogent Core framework architecture using text-based diagrams suitable for documentation and code review.

---

## 1. Layered Architecture Overview

```
+============================================================================+
|                          APPLICATION LAYER                                  |
|  +----------------------------------------------------------------------+  |
|  |  User Applications  |  Custom Widgets  |  App-Specific Logic         |  |
|  +----------------------------------------------------------------------+  |
+============================================================================+
                                    |
                                    v
+============================================================================+
|                           CORE PACKAGE                                      |
|  +------------------+  +------------------+  +------------------+          |
|  |     Widgets      |  |     Scenes       |  |     Stages       |          |
|  | Button, Slider,  |  | Widget tree root |  | Window, Dialog,  |          |
|  | TextField, List, |  | Painter owner    |  | Menu, Tooltip    |          |
|  | Frame, Tree...   |  | Event manager    |  | Lifecycle mgmt   |          |
|  +------------------+  +------------------+  +------------------+          |
|  +------------------+  +------------------+  +------------------+          |
|  |     Layout       |  |     Events       |  |     Render       |          |
|  | Size negotiation |  | Focus, Click,    |  | Scene rendering  |          |
|  | Position, Scroll |  | Drag, Keyboard   |  | Widget painting  |          |
|  +------------------+  +------------------+  +------------------+          |
+============================================================================+
                                    |
                                    v
+============================================================================+
|                        FOUNDATION PACKAGES                                  |
|  +------------+  +------------+  +------------+  +------------+            |
|  |   tree     |  |   styles   |  |   events   |  |   colors   |            |
|  | Node base  |  | Style      |  | Event defs |  | Schemes    |            |
|  | Plan mgmt  |  | Units      |  | Listeners  |  | Gradients  |            |
|  | Traversal  |  | States     |  | Types      |  | Uniform    |            |
|  +------------+  +------------+  +------------+  +------------+            |
|  +------------+  +------------+  +------------+  +------------+            |
|  |   paint    |  |   types    |  |   math32   |  |   base     |            |
|  | Painter    |  | Registry   |  | Vectors    |  | Errors     |            |
|  | Paths      |  | Code gen   |  | Matrices   |  | Utilities  |            |
|  | Images     |  | Reflection |  | Geometry   |  | I/O        |            |
|  +------------+  +------------+  +------------+  +------------+            |
+============================================================================+
                                    |
                                    v
+============================================================================+
|                       SYSTEM ABSTRACTION                                    |
|  +----------------------------------------------------------------------+  |
|  |                        system.App                                     |  |
|  |  Platform()  |  NewWindow()  |  Clipboard()  |  Cursor()  |  Quit()  |  |
|  +----------------------------------------------------------------------+  |
|  +----------------------------------------------------------------------+  |
|  |                       system.Window                                   |  |
|  |  Size()  |  Events()  |  Composer()  |  Close()  |  SetTitle()       |  |
|  +----------------------------------------------------------------------+  |
|  +----------------------------------------------------------------------+  |
|  |                      system/composer                                  |  |
|  |              Composition of rendered surfaces                         |  |
|  +----------------------------------------------------------------------+  |
+============================================================================+
                                    |
                                    v
+============================================================================+
|                        GPU / RENDERING                                      |
|  +----------------------------------------------------------------------+  |
|  |                           gpu (WebGPU)                                |  |
|  |  GPU  |  Device  |  Surface  |  Pipeline  |  Texture  |  Buffer     |  |
|  +----------------------------------------------------------------------+  |
|  +----------------------------------------------------------------------+  |
|  |                        paint/render                                   |  |
|  |  Renderer interface  |  Path ops  |  Image ops  |  Text ops          |  |
|  +----------------------------------------------------------------------+  |
|  +----------------------------------------------------------------------+  |
|  |                       gpu/gpudraw                                     |  |
|  |              GPU-accelerated 2D drawing primitives                    |  |
|  +----------------------------------------------------------------------+  |
+============================================================================+
                                    |
                                    v
+============================================================================+
|                       PLATFORM DRIVERS                                      |
|  +------------------+  +------------------+  +------------------+          |
|  |     Desktop      |  |       Web        |  |     Mobile       |          |
|  | GLFW + Native    |  | WASM + JS APIs   |  | iOS / Android    |          |
|  | macOS/Win/Linux  |  | Browser Canvas   |  | Native bindings  |          |
|  +------------------+  +------------------+  +------------------+          |
|  +------------------+                                                      |
|  |    Offscreen     |                                                      |
|  | Testing/Headless |                                                      |
|  +------------------+                                                      |
+============================================================================+
```

---

## 2. Tree and Widget Hierarchy

```
                          tree.Node (interface)
                               |
                               | AsTree() *NodeBase
                               | Init()
                               | OnAdd()
                               | Destroy()
                               | CopyFieldsFrom()
                               v
                       +---------------+
                       |  tree.NodeBase |
                       +---------------+
                       | Name           |
                       | This           |
                       | Parent         |
                       | Children       |
                       | Properties     |
                       | Updaters       |
                       | Makers         |
                       +---------------+
                               |
               +---------------+----------------+
               |               |                |
               v               v                v
        core.Widget      xyz.Node3D        Other tree types
        (interface)      (3D graphics)
               |
               | AsWidget() *WidgetBase
               | Style()
               | SizeUp/Down/Final()
               | Position()
               | Render()
               v
       +------------------+
       |  core.WidgetBase  |
       +------------------+
       | Tooltip           |
       | Parts             |
       | Geom              |
       | Styles            |
       | Stylers           |
       | Listeners         |
       | ContextMenus      |
       | Scene             |
       +------------------+
               |
    +----------+----------+----------+----------+
    |          |          |          |          |
    v          v          v          v          v
 Frame      Button     Slider    TextField   Text
    |
    +-- Frame embeds WidgetBase
    |   and implements Layouter
    |
    +-- Children laid out according
        to Style.Direction, Wrap, etc.
```

---

## 3. Scene and Stage Relationships

```
+------------------------------------------------------------------+
|                        renderWindow                               |
|  +------------------------------------------------------------+  |
|  |                      stages (Mains)                         |  |
|  |  +------------------------------------------------------+  |  |
|  |  |              Stage (WindowStage)                      |  |  |
|  |  |  +------------------------------------------------+  |  |  |
|  |  |  |                   Scene                         |  |  |  |
|  |  |  |  +---------+  +---------+  +---------+         |  |  |  |
|  |  |  |  | Frame   |  | Painter |  | Events  |         |  |  |  |
|  |  |  |  | (root)  |  | (render)|  | (input) |         |  |  |  |
|  |  |  |  +---------+  +---------+  +---------+         |  |  |  |
|  |  |  +------------------------------------------------+  |  |  |
|  |  |  +------------------------------------------------+  |  |  |
|  |  |  |              stages (Popups)                    |  |  |  |
|  |  |  |  +----------+  +----------+  +----------+      |  |  |  |
|  |  |  |  | Menu     |  | Tooltip  |  | Dialog   |      |  |  |  |
|  |  |  |  | Stage    |  | Stage    |  | Stage    |      |  |  |  |
|  |  |  |  | +Scene   |  | +Scene   |  | +Scene   |      |  |  |  |
|  |  |  |  +----------+  +----------+  +----------+      |  |  |  |
|  |  |  +------------------------------------------------+  |  |  |
|  |  +------------------------------------------------------+  |  |
|  |  +------------------------------------------------------+  |  |
|  |  |              Stage (DialogStage)                      |  |  |
|  |  |  ... (additional main stages stack)                   |  |  |
|  |  +------------------------------------------------------+  |  |
|  +------------------------------------------------------------+  |
+------------------------------------------------------------------+

Stage Types:
  MainStages:                      PopupStages:
  - WindowStage (full window)      - MenuStage
  - DialogStage (dialog/modal)     - TooltipStage
                                   - SnackbarStage
                                   - CompleterStage
```

---

## 4. Event Flow

```
+------------------+
|  System/Driver   |
| (glfw, web, etc) |
+--------+---------+
         |
         | Raw events (mouse, keyboard, touch)
         v
+------------------+
|  events.Source   |
| (in Window)      |
+--------+---------+
         |
         | events.Event
         v
+------------------+
|   renderWindow   |
| handleEvents()   |
+--------+---------+
         |
         | Route to appropriate Stage
         v
+------------------+
|     Stage        |
| (Main or Popup)  |
+--------+---------+
         |
         v
+-------------------+
|   Scene.Events    |
| handleEvent()     |
+--------+----------+
         |
    +----+----+
    |         |
    v         v
Position   Focus
Events     Events
    |         |
    v         v
+---------+  +---------+
|mouseIn  |  | focus   |
|BBox     |  | widget  |
|stack    |  |         |
+---------+  +---------+
    |            |
    v            v
+--------------------+
| Widget.HandleEvent |
|  Listeners.Call()  |
+--------------------+
         |
         | Event marked Handled
         | or bubbles up
         v
+--------------------+
|  Event Processing  |
|     Complete       |
+--------------------+

Key Events Flow:
  MouseDown -> press widget set
  MouseUp   -> Click event if same widget
  MouseMove -> Hover events, drag detection
  KeyDown   -> Chord generation
  KeyChord  -> Shortcut handling, then focus
```

---

## 5. Rendering Pipeline

```
+-------------------------------------------------------------------+
|                     Widget.Render()                                |
|  +-------------------------------------------------------------+  |
|  |  1. Compute styles (Style())                                |  |
|  |  2. Apply transforms                                        |  |
|  |  3. Draw background, border, shadows                        |  |
|  |  4. Draw content (text, images, children)                   |  |
|  +-------------------------------------------------------------+  |
+-------------------------------------------------------------------+
                                |
                                v
+-------------------------------------------------------------------+
|                     paint.Painter                                  |
|  +-------------------------------------------------------------+  |
|  |  Accumulate draw operations:                                 |  |
|  |  - MoveTo, LineTo, CubeTo (path building)                   |  |
|  |  - Rectangle, Circle, Ellipse (shapes)                      |  |
|  |  - Draw() -> add Path to render list                        |  |
|  |  - DrawText() -> add Text to render list                    |  |
|  |  - DrawImage() -> add Image op to render list               |  |
|  +-------------------------------------------------------------+  |
|  |  render.Render (accumulated operations list)                 |  |
|  +-------------------------------------------------------------+  |
+-------------------------------------------------------------------+
                                |
                                v
+-------------------------------------------------------------------+
|                   render.Renderer                                  |
|  +-------------------------------------------------------------+  |
|  |  Execute render operations to target:                        |  |
|  |  - Rasterize paths to pixels                                |  |
|  |  - Composite images                                         |  |
|  |  - Shape and render text                                    |  |
|  |  - Apply effects (blur, shadows)                            |  |
|  +-------------------------------------------------------------+  |
+-------------------------------------------------------------------+
                                |
                                v
+-------------------------------------------------------------------+
|                       Scene Image                                  |
|  +-------------------------------------------------------------+  |
|  |  RGBA image containing rendered scene                        |  |
|  +-------------------------------------------------------------+  |
+-------------------------------------------------------------------+
                                |
                                v
+-------------------------------------------------------------------+
|               system/composer.Composer                             |
|  +-------------------------------------------------------------+  |
|  |  Compose scenes, sprites, overlays                          |  |
|  |  Handle direct render widgets (video, 3D)                   |  |
|  +-------------------------------------------------------------+  |
+-------------------------------------------------------------------+
                                |
                                v
+-------------------------------------------------------------------+
|                  GPU Surface (WebGPU)                              |
|  +-------------------------------------------------------------+  |
|  |  Upload to GPU texture                                       |  |
|  |  Present to screen                                          |  |
|  +-------------------------------------------------------------+  |
+-------------------------------------------------------------------+
```

---

## 6. Widget Lifecycle

```
                    +-------------------+
                    |   tree.New[T]()   |
                    +-------------------+
                              |
                              v
                    +-------------------+
                    |    InitNode()     |
                    | - Set This ptr    |
                    | - Call Init()     |
                    +-------------------+
                              |
                              v
                    +-------------------+
                    |  Widget.Init()    |
                    | - Setup Stylers   |
                    | - Setup Listeners |
                    | - Setup Updater   |
                    +-------------------+
                              |
                              v
                    +-------------------+
                    |  AddChild/Insert  |
                    | - Set Parent      |
                    | - Call OnAdd()    |
                    +-------------------+
                              |
                              v
                    +-------------------+
                    |  Widget.OnAdd()   |
                    | - Set Scene ref   |
                    | - Apply WidgetInit|
                    +-------------------+
                              |
                              v
            +----------------------------------+
            |        UPDATE CYCLE              |
            |  +----------------------------+  |
            |  | 1. RunUpdaters()           |  |
            |  |    - Value bindings        |  |
            |  |    - UpdateFromMake()      |  |
            |  +----------------------------+  |
            |  | 2. Plan.Update()           |  |
            |  |    - Reconcile children    |  |
            |  +----------------------------+  |
            |  | 3. Style()                 |  |
            |  |    - Apply Stylers         |  |
            |  |    - ToDots()              |  |
            |  +----------------------------+  |
            +----------------------------------+
                              |
                              v
            +----------------------------------+
            |         LAYOUT CYCLE             |
            |  +----------------------------+  |
            |  | 1. SizeUp()                |  |
            |  |    - Bottom-up min sizes   |  |
            |  +----------------------------+  |
            |  | 2. SizeDown() [iterative]  |  |
            |  |    - Top-down allocation   |  |
            |  |    - Grow distribution     |  |
            |  +----------------------------+  |
            |  | 3. SizeFinal()             |  |
            |  |    - Final sizes           |  |
            |  +----------------------------+  |
            |  | 4. Position()              |  |
            |  |    - Set relative pos      |  |
            |  +----------------------------+  |
            |  | 5. ApplyScenePos()         |  |
            |  |    - Set absolute pos      |  |
            |  |    - Compute BBox          |  |
            |  +----------------------------+  |
            +----------------------------------+
                              |
                              v
            +----------------------------------+
            |         RENDER CYCLE             |
            |  +----------------------------+  |
            |  | RenderWidget()             |  |
            |  |   - Check visibility       |  |
            |  |   - Call Render()          |  |
            |  |   - Render Parts           |  |
            |  |   - Render Children        |  |
            |  +----------------------------+  |
            +----------------------------------+
                              |
                              v
                    +-------------------+
                    |    Destroy()      |
                    | - Delete children |
                    | - Set This = nil  |
                    +-------------------+
```

---

## 7. Styling System

```
+--------------------------------------------------------------------+
|                        Style Application                            |
+--------------------------------------------------------------------+

Widget.Style() called during update:

    +------------------+
    |  Styles.Defaults |
    |  (reset to base) |
    +--------+---------+
             |
             v
    +------------------+
    |  InheritFields   |
    |  (from parent)   |
    |  - Color         |
    |  - Opacity       |
    |  - Font settings |
    +--------+---------+
             |
             v
    +------------------+     +------------------+
    | Stylers.First    |---->| scene/app level  |
    +------------------+     | global styles    |
             |               +------------------+
             v
    +------------------+     +------------------+
    | Stylers.Normal   |---->| widget-specific  |
    +------------------+     | base styles      |
             |               +------------------+
             v
    +------------------+     +------------------+
    | Stylers.Final    |---->| override styles  |
    +------------------+     | state-dependent  |
             |               +------------------+
             v
    +------------------+
    |    ToDots()      |
    | Convert units to |
    | device pixels    |
    +--------+---------+
             |
             v
    +------------------+
    | ComputeActual-   |
    |   Background     |
    | Apply opacity,   |
    | state layers     |
    +------------------+


State-Dependent Styling:

    widget.Styler(func(s *styles.Style) {
        s.Background = colors.Scheme.Surface

        if s.Is(states.Hovered) {
            s.StateLayer = 0.08
        }
        if s.Is(states.Focused) {
            s.Border = s.MaxBorder
        }
        if s.Is(states.Disabled) {
            s.Opacity = 0.38
        }
    })
```

---

## 8. Platform Driver Architecture

```
+------------------------------------------------------------------+
|                      system.App (interface)                       |
+-------------------+----------------------------------------------+
| Platform()        | Returns current platform enum                 |
| NewWindow()       | Creates new system window                     |
| Clipboard()       | Returns clipboard interface                   |
| Cursor()          | Returns cursor interface                      |
| RunOnMain()       | Execute on main thread                        |
| MainLoop()        | Run app event loop                            |
+-------------------+----------------------------------------------+
                              |
          +-------------------+-------------------+
          |                   |                   |
          v                   v                   v
+------------------+  +------------------+  +------------------+
|  Desktop Driver  |  |    Web Driver    |  |  Mobile Driver   |
+------------------+  +------------------+  +------------------+
| GLFW window      |  | WASM/JS interop  |  | Native bridges   |
| Native menus     |  | Canvas element   |  | iOS: ObjectiveC  |
| File dialogs     |  | IndexedDB        |  | Android: JNI     |
| System clipboard |  | localStorage     |  | Touch handling   |
+------------------+  +------------------+  +------------------+
         |                    |                    |
         v                    v                    v
+------------------+  +------------------+  +------------------+
| macOS: Cocoa     |  | Browser APIs     |  | UIKit / Activity |
| Windows: Win32   |  | WebGPU/WebGL     |  | Platform GPU     |
| Linux: X11/Wayland| | DOM events       |  | Native events    |
+------------------+  +------------------+  +------------------+


Platform-Specific Files:

    system/driver/
    +-- base/           # Shared base implementations
    |   +-- app.go
    |   +-- window.go
    +-- desktop/        # Desktop platforms (GLFW)
    |   +-- app.go
    |   +-- window.go
    |   +-- desktop_darwin.go
    |   +-- desktop_linux.go
    |   +-- desktop_windows.go
    +-- web/            # WebAssembly
    |   +-- app.go
    |   +-- window.go
    |   +-- jsfs/       # Virtual filesystem
    +-- ios/            # iOS/iPadOS
    +-- android/        # Android
    +-- offscreen/      # Testing/headless
```

---

## 9. Package Dependencies (Simplified)

```
                            core
                              |
        +---------------------+---------------------+
        |                     |                     |
        v                     v                     v
      tree                 styles               events
        |                     |                     |
        +----------+----------+                     |
                   |                                |
                   v                                |
                 types  <---------------------------+
                   |
        +----------+----------+
        |          |          |
        v          v          v
      base      math32     colors
        |
        +------ errors
        +------ reflectx
        +------ fsx
        +------ ...

                              gpu
                               |
                               v
                        cogentcore/webgpu
                               |
                               v
                        wgpu (native)


Core imports these foundation packages:
  - tree (widget tree structure)
  - styles (styling system)
  - events (event definitions)
  - colors (color management)
  - math32 (geometry, vectors)
  - paint (rendering)
  - system (platform abstraction)
  - types (type registry)
  - base/* (utilities)
```

---

## 10. Concurrency Model

```
+------------------------------------------------------------------+
|                    Main Goroutine                                 |
|  +------------------------------------------------------------+  |
|  |  system.TheApp.MainLoop()                                   |  |
|  |    - Process system events                                  |  |
|  |    - Execute RunOnMain functions                            |  |
|  +------------------------------------------------------------+  |
+------------------------------------------------------------------+
         ^                              |
         |                              | Events
         | RunOnMain(f)                 v
         |                   +--------------------+
         |                   |  renderWindow      |
         |                   |  Event Loop        |
         |                   |  (per window)      |
         |                   +--------------------+
         |                              |
         |                              v
+------------------------------------------------------------------+
|                   Worker Goroutines                               |
|  +------------------------------------------------------------+  |
|  |  Async operations (file I/O, network, etc.)                 |  |
|  |                                                              |  |
|  |  MUST use AsyncLock/AsyncUnlock for widget updates:         |  |
|  |                                                              |  |
|  |    go func() {                                              |  |
|  |        result := doSlowWork()                               |  |
|  |        widget.AsyncLock()                                   |  |
|  |        widget.SetText(result)                               |  |
|  |        widget.NeedsRender()                                 |  |
|  |        widget.AsyncUnlock()                                 |  |
|  |    }()                                                      |  |
|  +------------------------------------------------------------+  |
+------------------------------------------------------------------+

Synchronization Points:

  +------------------+
  | renderContext    |
  | (per window)     |
  |  - Mutex         |
  |  - rebuild flag  |
  |  - textShaper    |
  +------------------+
          |
          | Lock()
          v
  +------------------+
  | Scene updates    |
  | - Style changes  |
  | - Layout changes |
  | - Render changes |
  +------------------+

Thread Safety Rules:
  1. All widget mutations must be on main thread OR use AsyncLock
  2. Event handlers run on main thread
  3. Stylers, Makers, Updaters run on main thread
  4. Rendering runs on main thread
  5. Worker goroutines must use AsyncLock for any widget access
```

---

## Notes on Diagrams

These ASCII diagrams are designed to be:
- Version control friendly (text-based)
- Renderable in any text viewer
- Easy to update as architecture evolves
- Useful for code review discussions

For more detailed visual diagrams, consider:
- Mermaid.js for GitHub rendering
- PlantUML for detailed UML
- draw.io for interactive editing
