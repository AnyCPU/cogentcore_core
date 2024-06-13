// Code generated by "core generate"; DO NOT EDIT.

package main

import (
	"fmt"
	"image"
	"image/draw"
	"maps"
	"strings"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/colors"
	"cogentcore.org/core/colors/gradient"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/keymap"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/pages"
	"cogentcore.org/core/paint"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/texteditor"
	"cogentcore.org/core/tree"
	"cogentcore.org/core/views"
)

func init() {
	maps.Copy(pages.Examples, PagesExamples)
}

// PagesExamples are the compiled pages examples for this app.
var PagesExamples = map[string]func(parent core.Widget){
	"getting-started/hello-world-0": func(parent core.Widget) {
		b := parent
		core.NewButton(b).SetText("Hello, World!")
	},
	"basics/widgets-0": func(parent core.Widget) {
		core.NewButton(parent).SetText("Click me!").SetIcon(icons.Add)
	},
	"basics/events-0": func(parent core.Widget) {
		core.NewButton(parent).SetText("Click me!").OnClick(func(e events.Event) {
			core.MessageSnackbar(parent, "Button clicked")
		})
	},
	"basics/events-1": func(parent core.Widget) {
		core.NewButton(parent).SetText("Click me!").OnClick(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprint("Button clicked at ", e.Pos()))
		})
	},
	"basics/styling-0": func(parent core.Widget) {
		core.NewText(parent).SetText("Bold text").Styler(func(s *styles.Style) {
			s.Font.Weight = styles.WeightBold
		})
	},
	"basics/styling-1": func(parent core.Widget) {
		core.NewButton(parent).SetText("Success button").Styler(func(s *styles.Style) {
			s.Background = colors.C(colors.Scheme.Success.Base)
			s.Color = colors.C(colors.Scheme.Success.On)
		})
	},
	"basics/styling-2": func(parent core.Widget) {
		core.NewFrame(parent).Styler(func(s *styles.Style) {
			s.Min.Set(units.Dp(50))
			s.Background = colors.C(colors.Scheme.Primary.Base)
		})
	},
	"widgets/buttons-0": func(parent core.Widget) {
		core.NewButton(parent).SetText("Download")
	},
	"widgets/buttons-1": func(parent core.Widget) {
		core.NewButton(parent).SetIcon(icons.Download)
	},
	"widgets/buttons-2": func(parent core.Widget) {
		core.NewButton(parent).SetText("Download").SetIcon(icons.Download)
	},
	"widgets/buttons-3": func(parent core.Widget) {
		core.NewButton(parent).SetText("Send").SetIcon(icons.Send).OnClick(func(e events.Event) {
			core.MessageSnackbar(parent, "Message sent")
		})
	},
	"widgets/buttons-4": func(parent core.Widget) {
		core.NewButton(parent).SetText("Share").SetIcon(icons.Share).SetMenu(func(m *core.Scene) {
			core.NewButton(m).SetText("Copy link")
			core.NewButton(m).SetText("Send message")
		})
	},
	"widgets/buttons-5": func(parent core.Widget) {
		core.NewButton(parent).SetText("Save").SetShortcut("Command+S").OnClick(func(e events.Event) {
			core.MessageSnackbar(parent, "File saved")
		})
	},
	"widgets/buttons-6": func(parent core.Widget) {
		core.NewButton(parent).SetText("Open").SetKey(keymap.Open).OnClick(func(e events.Event) {
			core.MessageSnackbar(parent, "File opened")
		})
	},
	"widgets/buttons-7": func(parent core.Widget) {
		core.NewButton(parent).SetType(core.ButtonFilled).SetText("Filled")
	},
	"widgets/buttons-8": func(parent core.Widget) {
		core.NewButton(parent).SetType(core.ButtonTonal).SetText("Tonal")
	},
	"widgets/buttons-9": func(parent core.Widget) {
		core.NewButton(parent).SetType(core.ButtonElevated).SetText("Elevated")
	},
	"widgets/buttons-10": func(parent core.Widget) {
		core.NewButton(parent).SetType(core.ButtonOutlined).SetText("Outlined")
	},
	"widgets/buttons-11": func(parent core.Widget) {
		core.NewButton(parent).SetType(core.ButtonText).SetText("Text")
	},
	"widgets/buttons-12": func(parent core.Widget) {
		core.NewButton(parent).SetType(core.ButtonAction).SetText("Action")
	},
	"widgets/canvases-0": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.FillBox(math32.Vector2{}, math32.Vec2(1, 1), colors.C(colors.Scheme.Primary.Base))
		})
	},
	"widgets/canvases-1": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.MoveTo(0, 0)
			pc.LineTo(1, 1)
			pc.StrokeStyle.Color = colors.C(colors.Scheme.Error.Base)
			pc.Stroke()
		})
	},
	"widgets/canvases-2": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.MoveTo(0, 0)
			pc.LineTo(1, 1)
			pc.StrokeStyle.Color = colors.C(colors.Scheme.Error.Base)
			pc.StrokeStyle.Width.Dp(8)
			pc.ToDots()
			pc.Stroke()
		})
	},
	"widgets/canvases-3": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.DrawCircle(0.5, 0.5, 0.5)
			pc.FillStyle.Color = colors.C(colors.Scheme.Success.Base)
			pc.Fill()
		})
	},
	"widgets/canvases-4": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.DrawEllipse(0.5, 0.5, 0.5, 0.25)
			pc.FillStyle.Color = colors.C(colors.Scheme.Success.Base)
			pc.Fill()
		})
	},
	"widgets/canvases-5": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.DrawEllipticalArc(0.5, 0.5, 0.5, 0.25, math32.Pi, 2*math32.Pi)
			pc.FillStyle.Color = colors.C(colors.Scheme.Success.Base)
			pc.Fill()
		})
	},
	"widgets/canvases-6": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.DrawRegularPolygon(6, 0.5, 0.5, 0.5, math32.Pi)
			pc.FillStyle.Color = colors.C(colors.Scheme.Success.Base)
			pc.Fill()
		})
	},
	"widgets/canvases-7": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.MoveTo(0, 0)
			pc.QuadraticTo(0.5, 0.25, 1, 1)
			pc.StrokeStyle.Color = colors.C(colors.Scheme.Error.Base)
			pc.Stroke()
		})
	},
	"widgets/canvases-8": func(parent core.Widget) {
		core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.MoveTo(0, 0)
			pc.CubicTo(0.5, 0.25, 0.25, 0.5, 1, 1)
			pc.StrokeStyle.Color = colors.C(colors.Scheme.Error.Base)
			pc.Stroke()
		})
	},
	"widgets/canvases-9": func(parent core.Widget) {
		c := core.NewCanvas(parent).SetDraw(func(pc *paint.Context) {
			pc.FillBox(math32.Vector2{}, math32.Vec2(1, 1), colors.C(colors.Scheme.Warn.Base))
		})
		c.Styler(func(s *styles.Style) {
			s.Min.Set(units.Dp(128))
		})
	},
	"widgets/choosers-0": func(parent core.Widget) {
		core.NewChooser(parent).SetStrings("macOS", "Windows", "Linux")
	},
	"widgets/choosers-1": func(parent core.Widget) {
		core.NewChooser(parent).SetItems(
			core.ChooserItem{Value: "Computer", Icon: icons.Computer, Tooltip: "Use a computer"},
			core.ChooserItem{Value: "Phone", Icon: icons.Smartphone, Tooltip: "Use a phone"},
		)
	},
	"widgets/choosers-2": func(parent core.Widget) {
		core.NewChooser(parent).SetPlaceholder("Choose a platform").SetStrings("macOS", "Windows", "Linux")
	},
	"widgets/choosers-3": func(parent core.Widget) {
		core.NewChooser(parent).SetStrings("Apple", "Orange", "Strawberry").SetCurrentValue("Orange")
	},
	"widgets/choosers-4": func(parent core.Widget) {
		core.NewChooser(parent).SetType(core.ChooserOutlined).SetStrings("Apple", "Orange", "Strawberry")
	},
	"widgets/choosers-5": func(parent core.Widget) {
		core.NewChooser(parent).SetIcon(icons.Sort).SetStrings("Newest", "Oldest", "Popular")
	},
	"widgets/choosers-6": func(parent core.Widget) {
		core.NewChooser(parent).SetEditable(true).SetStrings("Newest", "Oldest", "Popular")
	},
	"widgets/choosers-7": func(parent core.Widget) {
		core.NewChooser(parent).SetAllowNew(true).SetStrings("Newest", "Oldest", "Popular")
	},
	"widgets/choosers-8": func(parent core.Widget) {
		core.NewChooser(parent).SetEditable(true).SetAllowNew(true).SetStrings("Newest", "Oldest", "Popular")
	},
	"widgets/choosers-9": func(parent core.Widget) {
		ch := core.NewChooser(parent).SetStrings("Newest", "Oldest", "Popular")
		ch.OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("Sorting by %v", ch.CurrentItem.Value))
		})
	},
	"widgets/dialogs-0": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("Message")
		bt.OnClick(func(e events.Event) {
			core.MessageDialog(bt, "Something happened", "Message")
		})
	},
	"widgets/dialogs-1": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("Error")
		bt.OnClick(func(e events.Event) {
			core.ErrorDialog(bt, errors.New("invalid encoding format"), "Error loading file")
		})
	},
	"widgets/dialogs-2": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("Confirm")
		bt.OnClick(func(e events.Event) {
			d := core.NewBody().AddTitle("Confirm").AddText("Send message?")
			d.AddBottomBar(func(parent core.Widget) {
				d.AddCancel(parent).OnClick(func(e events.Event) {
					core.MessageSnackbar(bt, "Dialog canceled")
				})
				d.AddOK(parent).OnClick(func(e events.Event) {
					core.MessageSnackbar(bt, "Dialog accepted")
				})
			})
			d.RunDialog(bt)
		})
	},
	"widgets/dialogs-3": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("Input")
		bt.OnClick(func(e events.Event) {
			d := core.NewBody().AddTitle("Input").AddText("What is your name?")
			tf := core.NewTextField(d)
			d.AddBottomBar(func(parent core.Widget) {
				d.AddCancel(parent)
				d.AddOK(parent).OnClick(func(e events.Event) {
					core.MessageSnackbar(bt, "Your name is "+tf.Text())
				})
			})
			d.RunDialog(bt)
		})
	},
	"widgets/dialogs-4": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("Full window")
		bt.OnClick(func(e events.Event) {
			d := core.NewBody().AddTitle("Full window dialog")
			d.RunFullDialog(bt)
		})
	},
	"widgets/dialogs-5": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("New window")
		bt.OnClick(func(e events.Event) {
			d := core.NewBody().AddTitle("New window dialog")
			d.RunWindowDialog(bt)
		})
	},
	"widgets/frames-0": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-1": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-2": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Gap.Set(units.Em(2))
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-3": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Max.X.Em(10)
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-4": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Overflow.X = styles.OverflowAuto
			s.Max.X.Em(10)
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-5": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Wrap = true
			s.Max.X.Em(10)
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-6": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Display = styles.Grid
			s.Columns = 2
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
		core.NewButton(fr).SetText("Fourth")
	},
	"widgets/frames-7": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Background = colors.C(colors.Scheme.Warn.Container)
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-8": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Background = gradient.NewLinear().AddStop(colors.Yellow, 0).AddStop(colors.Orange, 0.5).AddStop(colors.Red, 1)
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-9": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Border.Width.Set(units.Dp(4))
			s.Border.Color.Set(colors.C(colors.Scheme.Outline))
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-10": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Border.Radius = styles.BorderRadiusLarge
			s.Border.Width.Set(units.Dp(4))
			s.Border.Color.Set(colors.C(colors.Scheme.Outline))
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/frames-11": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 1)
			s.Border.Width.Set(units.Dp(4))
			s.Border.Color.Set(colors.C(colors.Scheme.Outline))
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
	"widgets/icons-0": func(parent core.Widget) {
		core.NewButton(parent).SetIcon(icons.Send)
	},
	"widgets/icons-1": func(parent core.Widget) {
		core.NewIcon(parent).SetIcon(icons.Home)
	},
	"widgets/icons-2": func(parent core.Widget) {
		core.NewButton(parent).SetIcon(icons.Home.Fill())
	},
	"widgets/images-0": func(parent core.Widget) {
		errors.Log(core.NewImage(parent).OpenFS(myImage, "image.png"))
	},
	"widgets/images-1": func(parent core.Widget) {
		img := core.NewImage(parent)
		errors.Log(img.OpenFS(myImage, "image.png"))
		img.Styler(func(s *styles.Style) {
			s.Min.Set(units.Dp(256))
		})
	},
	"widgets/images-2": func(parent core.Widget) {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		draw.Draw(img, image.Rect(10, 5, 100, 90), colors.C(colors.Scheme.Warn.Container), image.Point{}, draw.Src)
		draw.Draw(img, image.Rect(20, 20, 60, 50), colors.C(colors.Scheme.Success.Base), image.Point{}, draw.Src)
		draw.Draw(img, image.Rect(60, 70, 80, 100), colors.C(colors.Scheme.Error.Base), image.Point{}, draw.Src)
		core.NewImage(parent).SetImage(img)
	},
	"widgets/meters-0": func(parent core.Widget) {
		core.NewMeter(parent)
	},
	"widgets/meters-1": func(parent core.Widget) {
		core.NewMeter(parent).SetValue(0.7)
	},
	"widgets/meters-2": func(parent core.Widget) {
		core.NewMeter(parent).SetMin(5.7).SetMax(18).SetValue(10.2)
	},
	"widgets/meters-3": func(parent core.Widget) {
		core.NewMeter(parent).Styler(func(s *styles.Style) {
			s.Direction = styles.Column
		})
	},
	"widgets/meters-4": func(parent core.Widget) {
		core.NewMeter(parent).SetType(core.MeterCircle)
	},
	"widgets/meters-5": func(parent core.Widget) {
		core.NewMeter(parent).SetType(core.MeterSemicircle)
	},
	"widgets/meters-6": func(parent core.Widget) {
		core.NewMeter(parent).SetType(core.MeterCircle).SetText("50%")
	},
	"widgets/meters-7": func(parent core.Widget) {
		core.NewMeter(parent).SetType(core.MeterSemicircle).SetText("50%")
	},
	"widgets/sliders-0": func(parent core.Widget) {
		core.NewSlider(parent)
	},
	"widgets/sliders-1": func(parent core.Widget) {
		core.NewSlider(parent).SetValue(0.7)
	},
	"widgets/sliders-2": func(parent core.Widget) {
		core.NewSlider(parent).SetMin(5.7).SetMax(18).SetValue(10.2)
	},
	"widgets/sliders-3": func(parent core.Widget) {
		core.NewSlider(parent).SetStep(0.2)
	},
	"widgets/sliders-4": func(parent core.Widget) {
		core.NewSlider(parent).SetStep(0.2).SetEnforceStep(true)
	},
	"widgets/sliders-5": func(parent core.Widget) {
		core.NewSlider(parent).SetIcon(icons.DeployedCode.Fill())
	},
	"widgets/sliders-6": func(parent core.Widget) {
		sr := core.NewSlider(parent)
		sr.OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("OnChange: %v", sr.Value))
		})
	},
	"widgets/sliders-7": func(parent core.Widget) {
		sr := core.NewSlider(parent)
		sr.OnInput(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("OnInput: %v", sr.Value))
		})
	},
	"widgets/snackbars-0": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("Message")
		bt.OnClick(func(e events.Event) {
			core.MessageSnackbar(bt, "New messages loaded")
		})
	},
	"widgets/snackbars-1": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("Error")
		bt.OnClick(func(e events.Event) {
			core.ErrorSnackbar(bt, errors.New("file not found"), "Error loading page")
		})
	},
	"widgets/snackbars-2": func(parent core.Widget) {
		bt := core.NewButton(parent).SetText("Custom")
		bt.OnClick(func(e events.Event) {
			core.NewBody().AddSnackbarText("Files updated").
				AddSnackbarButton("Refresh", func(e events.Event) {
					core.MessageSnackbar(bt, "Refreshed files")
				}).AddSnackbarIcon(icons.Close).NewSnackbar(bt).Run()
		})
	},
	"widgets/spinners-0": func(parent core.Widget) {
		core.NewSpinner(parent)
	},
	"widgets/spinners-1": func(parent core.Widget) {
		core.NewSpinner(parent).SetValue(12.7)
	},
	"widgets/spinners-2": func(parent core.Widget) {
		core.NewSpinner(parent).SetMin(-0.5).SetMax(2.7)
	},
	"widgets/spinners-3": func(parent core.Widget) {
		core.NewSpinner(parent).SetStep(6)
	},
	"widgets/spinners-4": func(parent core.Widget) {
		core.NewSpinner(parent).SetStep(4).SetEnforceStep(true)
	},
	"widgets/spinners-5": func(parent core.Widget) {
		core.NewSpinner(parent).SetType(core.TextFieldOutlined)
	},
	"widgets/spinners-6": func(parent core.Widget) {
		core.NewSpinner(parent).SetFormat("%X").SetStep(1).SetValue(44)
	},
	"widgets/spinners-7": func(parent core.Widget) {
		sp := core.NewSpinner(parent)
		sp.OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("Value changed to %g", sp.Value))
		})
	},
	"widgets/splits-0": func(parent core.Widget) {
		sp := core.NewSplits(parent)
		core.NewText(sp).SetText("First")
		core.NewText(sp).SetText("Second")
	},
	"widgets/splits-1": func(parent core.Widget) {
		sp := core.NewSplits(parent)
		core.NewText(sp).SetText("First")
		core.NewText(sp).SetText("Second")
		core.NewText(sp).SetText("Third")
		core.NewText(sp).SetText("Fourth")
	},
	"widgets/splits-2": func(parent core.Widget) {
		sp := core.NewSplits(parent).SetSplits(0.2, 0.8)
		core.NewText(sp).SetText("First")
		core.NewText(sp).SetText("Second")
	},
	"widgets/splits-3": func(parent core.Widget) {
		sp := core.NewSplits(parent)
		sp.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
		})
		core.NewText(sp).SetText("First")
		core.NewText(sp).SetText("Second")
	},
	"widgets/splits-4": func(parent core.Widget) {
		sp := core.NewSplits(parent)
		sp.Styler(func(s *styles.Style) {
			s.Direction = styles.Row
		})
		core.NewText(sp).SetText("First")
		core.NewText(sp).SetText("Second")
	},
	"widgets/svgs-0": func(parent core.Widget) {
		errors.Log(core.NewSVG(parent).OpenFS(mySVG, "icon.svg"))
	},
	"widgets/svgs-1": func(parent core.Widget) {
		svg := core.NewSVG(parent)
		errors.Log(svg.OpenFS(mySVG, "icon.svg"))
		svg.Styler(func(s *styles.Style) {
			s.Min.Set(units.Dp(128))
		})
	},
	"widgets/svgs-2": func(parent core.Widget) {
		svg := core.NewSVG(parent)
		svg.SetReadOnly(false)
		errors.Log(svg.OpenFS(mySVG, "icon.svg"))
	},
	"widgets/svgs-3": func(parent core.Widget) {
		errors.Log(core.NewSVG(parent).ReadString(`<rect width="100" height="100" fill="red"/>`))
	},
	"widgets/switches-0": func(parent core.Widget) {
		core.NewSwitch(parent)
	},
	"widgets/switches-1": func(parent core.Widget) {
		core.NewSwitch(parent).SetText("Remember me")
	},
	"widgets/switches-2": func(parent core.Widget) {
		core.NewSwitch(parent).SetType(core.SwitchCheckbox).SetText("Remember me")
	},
	"widgets/switches-3": func(parent core.Widget) {
		core.NewSwitch(parent).SetType(core.SwitchRadioButton).SetText("Remember me")
	},
	"widgets/switches-4": func(parent core.Widget) {
		sw := core.NewSwitch(parent).SetText("Remember me")
		sw.OnChange(func(e events.Event) {
			core.MessageSnackbar(sw, fmt.Sprintf("Switch is %v", sw.IsChecked()))
		})
	},
	"widgets/switches-5": func(parent core.Widget) {
		core.NewSwitches(parent).SetStrings("Go", "Python", "C++")
	},
	"widgets/switches-6": func(parent core.Widget) {
		core.NewSwitches(parent).SetItems(
			core.SwitchItem{Value: "Go", Tooltip: "Elegant, fast, and easy-to-use"},
			core.SwitchItem{Value: "Python", Tooltip: "Slow and duck-typed"},
			core.SwitchItem{Value: "C++", Tooltip: "Hard to use and slow to compile"},
		)
	},
	"widgets/switches-7": func(parent core.Widget) {
		core.NewSwitches(parent).SetMutex(true).SetStrings("Go", "Python", "C++")
	},
	"widgets/switches-8": func(parent core.Widget) {
		core.NewSwitches(parent).SetType(core.SwitchChip).SetStrings("Go", "Python", "C++")
	},
	"widgets/switches-9": func(parent core.Widget) {
		core.NewSwitches(parent).SetType(core.SwitchCheckbox).SetStrings("Go", "Python", "C++")
	},
	"widgets/switches-10": func(parent core.Widget) {
		core.NewSwitches(parent).SetType(core.SwitchRadioButton).SetStrings("Go", "Python", "C++")
	},
	"widgets/switches-11": func(parent core.Widget) {
		core.NewSwitches(parent).SetType(core.SwitchSegmentedButton).SetStrings("Go", "Python", "C++")
	},
	"widgets/switches-12": func(parent core.Widget) {
		core.NewSwitches(parent).SetStrings("Go", "Python", "C++").Styler(func(s *styles.Style) {
			s.Direction = styles.Column
		})
	},
	"widgets/switches-13": func(parent core.Widget) {
		sw := core.NewSwitches(parent).SetStrings("Go", "Python", "C++")
		sw.OnChange(func(e events.Event) {
			core.MessageSnackbar(sw, fmt.Sprintf("Currently selected: %v", sw.SelectedItems()))
		})
	},
	"widgets/tabs-0": func(parent core.Widget) {
		ts := core.NewTabs(parent)
		ts.NewTab("First")
		ts.NewTab("Second")
	},
	"widgets/tabs-1": func(parent core.Widget) {
		ts := core.NewTabs(parent)
		first := ts.NewTab("First")
		core.NewText(first).SetText("I am first!")
		second := ts.NewTab("Second")
		core.NewText(second).SetText("I am second!")
	},
	"widgets/tabs-2": func(parent core.Widget) {
		ts := core.NewTabs(parent)
		ts.NewTab("First")
		ts.NewTab("Second")
		ts.NewTab("Third")
		ts.NewTab("Fourth")
	},
	"widgets/tabs-3": func(parent core.Widget) {
		ts := core.NewTabs(parent)
		ts.NewTab("First", icons.Home)
		ts.NewTab("Second", icons.Explore)
	},
	"widgets/tabs-4": func(parent core.Widget) {
		ts := core.NewTabs(parent).SetType(core.FunctionalTabs)
		ts.NewTab("First")
		ts.NewTab("Second")
		ts.NewTab("Third")
	},
	"widgets/tabs-5": func(parent core.Widget) {
		ts := core.NewTabs(parent).SetType(core.NavigationAuto)
		ts.NewTab("First", icons.Home)
		ts.NewTab("Second", icons.Explore)
		ts.NewTab("Third", icons.History)
	},
	"widgets/tabs-6": func(parent core.Widget) {
		ts := core.NewTabs(parent).SetNewTabButton(true)
		ts.NewTab("First")
		ts.NewTab("Second")
	},
	"widgets/text-fields-0": func(parent core.Widget) {
		core.NewTextField(parent)
	},
	"widgets/text-fields-1": func(parent core.Widget) {
		core.NewText(parent).SetText("Name:")
		core.NewTextField(parent).SetPlaceholder("Jane Doe")
	},
	"widgets/text-fields-2": func(parent core.Widget) {
		core.NewTextField(parent).SetText("Hello, world!")
	},
	"widgets/text-fields-3": func(parent core.Widget) {
		core.NewTextField(parent).SetText("This is a long sentence that demonstrates how text field content can overflow onto multiple lines")
	},
	"widgets/text-fields-4": func(parent core.Widget) {
		core.NewTextField(parent).SetType(core.TextFieldOutlined)
	},
	"widgets/text-fields-5": func(parent core.Widget) {
		core.NewTextField(parent).SetTypePassword()
	},
	"widgets/text-fields-6": func(parent core.Widget) {
		core.NewTextField(parent).AddClearButton()
	},
	"widgets/text-fields-7": func(parent core.Widget) {
		core.NewTextField(parent).SetLeadingIcon(icons.Euro).SetTrailingIcon(icons.OpenInNew, func(e events.Event) {
			core.MessageSnackbar(parent, "Opening shopping cart")
		})
	},
	"widgets/text-fields-8": func(parent core.Widget) {
		tf := core.NewTextField(parent)
		tf.SetValidator(func() error {
			if !strings.Contains(tf.Text(), "Go") {
				return errors.New("Must contain Go")
			}
			return nil
		})
	},
	"widgets/text-fields-9": func(parent core.Widget) {
		tf := core.NewTextField(parent)
		tf.OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, "OnChange: "+tf.Text())
		})
	},
	"widgets/text-fields-10": func(parent core.Widget) {
		tf := core.NewTextField(parent)
		tf.OnInput(func(e events.Event) {
			core.MessageSnackbar(parent, "OnInput: "+tf.Text())
		})
	},
	"widgets/text-0": func(parent core.Widget) {
		core.NewText(parent).SetText("Hello, world!")
	},
	"widgets/text-1": func(parent core.Widget) {
		core.NewText(parent).SetText("This is a very long sentence that demonstrates how text content will overflow onto multiple lines when the size of the text exceeds the size of its surrounding container; text widgets are customizable widget that Cogent Core provides, allowing you to display many kinds of text")
	},
	"widgets/text-2": func(parent core.Widget) {
		core.NewText(parent).SetText(`<b>You</b> can use <i>HTML</i> <u>formatting</u> inside of <b><i><u>Cogent Core</u></i></b> text, including <span style="color:red;background-color:yellow">custom styling</span> and <a href="https://example.com">links</a>`)
	},
	"widgets/text-3": func(parent core.Widget) {
		core.NewText(parent).SetType(core.TextHeadlineMedium).SetText("Hello, world!")
	},
	"widgets/text-4": func(parent core.Widget) {
		core.NewText(parent).SetText("Hello,\n\tworld!").Styler(func(s *styles.Style) {
			s.Font.Size.Dp(21)
			s.Font.Style = styles.Italic
			s.Text.WhiteSpace = styles.WhiteSpacePre
			s.Color = colors.C(colors.Scheme.Success.Base)
			s.SetMono(true)
		})
	},
	"widgets/tooltips-0": func(parent core.Widget) {
		core.NewButton(parent).SetIcon(icons.Add).SetTooltip("Add a new item to the list")
	},
	"widgets/tooltips-1": func(parent core.Widget) {
		core.NewSlider(parent)
	},
	"views/values-0": func(parent core.Widget) {
		// views.NewValue(parent, colors.Orange)
	},
	"views/values-1": func(parent core.Widget) {
		// t := time.Now()
		//
		//	views.NewValue(parent, &t).OnChange(func(e events.Event) {
		//	    core.MessageSnackbar(parent, "The time is "+t.Format(time.DateTime))
		//	})
	},
	"views/values-2": func(parent core.Widget) {
		// views.NewValue(parent, 70, `view:"slider"`)
	},
	"views/keyed-lists-0": func(parent core.Widget) {
		views.NewKeyedList(parent).SetMap(&map[string]int{"Go": 1, "C++": 3, "Python": 5})
	},
	"views/keyed-lists-1": func(parent core.Widget) {
		views.NewKeyedList(parent).SetInline(true).SetMap(&map[string]int{"Go": 1, "C++": 3})
	},
	"views/keyed-lists-2": func(parent core.Widget) {
		m := map[string]int{"Go": 1, "C++": 3, "Python": 5}
		views.NewKeyedList(parent).SetMap(&m).OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("Map: %v", m))
		})
	},
	"views/keyed-lists-3": func(parent core.Widget) {
		views.NewKeyedList(parent).SetMap(&map[string]int{"Go": 1, "C++": 3, "Python": 5}).SetReadOnly(true)
	},
	"views/keyed-lists-4": func(parent core.Widget) {
		views.NewKeyedList(parent).SetMap(&map[string]any{"Go": 1, "C++": "C-like", "Python": true})
	},
	"views/keyed-lists-5": func(parent core.Widget) {
		// views.NewValue(parent, &map[string]int{"Go": 1, "C++": 3})
	},
	"views/keyed-lists-6": func(parent core.Widget) {
		// views.NewValue(parent, &map[string]int{"Go": 1, "C++": 3, "Python": 5})
	},
	"views/lists-0": func(parent core.Widget) {
		views.NewList(parent).SetSlice(&[]int{1, 3, 5})
	},
	"views/lists-1": func(parent core.Widget) {
		views.NewInlineList(parent).SetSlice(&[]int{1, 3, 5})
	},
	"views/lists-2": func(parent core.Widget) {
		sl := []int{1, 3, 5}
		views.NewList(parent).SetSlice(&sl).OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("Slice: %v", sl))
		})
	},
	"views/lists-3": func(parent core.Widget) {
		views.NewList(parent).SetSlice(&[]int{1, 3, 5}).SetReadOnly(true)
	},
	"views/lists-4": func(parent core.Widget) {
		// views.NewValue(parent, &[]int{1, 3, 5})
	},
	"views/lists-5": func(parent core.Widget) {
		// views.NewValue(parent, &[]int{1, 3, 5, 7, 9})
	},
	"views/forms-0": func(parent core.Widget) {
		type person struct {
			Name string
			Age  int
		}
		views.NewForm(parent).SetStruct(&person{Name: "Go", Age: 35})
	},
	"views/forms-1": func(parent core.Widget) {
		type person struct {
			Name string
			Age  int
		}
		views.NewForm(parent).SetInline(true).SetStruct(&person{Name: "Go", Age: 35})
	},
	"views/forms-2": func(parent core.Widget) {
		type person struct {
			Name string
			Age  int
		}
		p := person{Name: "Go", Age: 35}
		views.NewForm(parent).SetStruct(&p).OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("You are %v", p))
		})
	},
	"views/forms-3": func(parent core.Widget) {
		type person struct {
			Name string `immediate:"+"`
			Age  int
		}
		p := person{Name: "Go", Age: 35}
		views.NewForm(parent).SetStruct(&p).OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("You are %v", p))
		})
	},
	"views/forms-4": func(parent core.Widget) {
		type person struct {
			Name string
			Age  int `view:"-"`
		}
		views.NewForm(parent).SetStruct(&person{Name: "Go", Age: 35})
	},
	"views/forms-5": func(parent core.Widget) {
		type person struct {
			Name string `edit:"-"`
			Age  int
		}
		views.NewForm(parent).SetStruct(&person{Name: "Go", Age: 35})
	},
	"views/forms-6": func(parent core.Widget) {
		type person struct {
			Name string
			Age  int
		}
		views.NewForm(parent).SetStruct(&person{Name: "Go", Age: 35}).SetReadOnly(true)
	},
	"views/forms-7": func(parent core.Widget) {
		type Person struct {
			Name string
			Age  int
		}
		type employee struct {
			Person
			Role string
		}
		views.NewForm(parent).SetStruct(&employee{Person{Name: "Go", Age: 35}, "Programmer"})
	},
	"views/forms-8": func(parent core.Widget) {
		type person struct {
			Name string
			Age  int
		}
		type employee struct {
			Role    string
			Manager person `view:"add-fields"`
		}
		views.NewForm(parent).SetStruct(&employee{"Programmer", person{Name: "Go", Age: 35}})
	},
	"views/forms-9": func(parent core.Widget) {
		type person struct {
			Name      string `default:"Gopher"`
			Age       int    `default:"20:30"`
			Precision int    `default:"64,32"`
		}
		views.NewForm(parent).SetStruct(&person{Name: "Go", Age: 35, Precision: 50})
	},
	"views/forms-10": func(parent core.Widget) {
		type person struct {
			Name string
			Age  int
		}
		// core.NewValue(&person{Name: "Go", Age: 35}, "", parent)
	},
	"views/forms-11": func(parent core.Widget) {
		type person struct {
			Name        string
			Age         int
			Job         string
			LikesGo     bool
			LikesPython bool
		}
		// core.NewValue(&person{Name: "Go", Age: 35, Job: "Programmer", LikesGo: true}, "", parent)
	},
	"views/tables-0": func(parent core.Widget) {
		type language struct {
			Name   string
			Rating int
		}
		views.NewTable(parent).SetSlice(&[]language{{"Go", 10}, {"Python", 5}})
	},
	"views/tables-1": func(parent core.Widget) {
		type language struct {
			Name   string
			Rating int
		}
		sl := []language{{"Go", 10}, {"Python", 5}}
		views.NewTable(parent).SetSlice(&sl).OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, fmt.Sprintf("Languages: %v", sl))
		})
	},
	"views/tables-2": func(parent core.Widget) {
		type language struct {
			Name   string
			Rating int `view:"-"`
		}
		views.NewTable(parent).SetSlice(&[]language{{"Go", 10}, {"Python", 5}})
	},
	"views/tables-3": func(parent core.Widget) {
		type language struct {
			Name   string
			Rating int `view:"-" table:"+"`
		}
		views.NewTable(parent).SetSlice(&[]language{{"Go", 10}, {"Python", 5}})
	},
	"views/tables-4": func(parent core.Widget) {
		type language struct {
			Name   string `edit:"-"`
			Rating int
		}
		views.NewTable(parent).SetSlice(&[]language{{"Go", 10}, {"Python", 5}})
	},
	"views/tables-5": func(parent core.Widget) {
		type language struct {
			Name   string
			Rating int
		}
		views.NewTable(parent).SetSlice(&[]language{{"Go", 10}, {"Python", 5}}).SetReadOnly(true)
	},
	"views/tables-6": func(parent core.Widget) {
		type language struct {
			Name   string
			Rating int
		}
		// views.NewValue(parent, &[]language{{"Go", 10}, {"Python", 5}})
	},
	"views/text-editors-0": func(parent core.Widget) {
		texteditor.NewSoloEditor(parent)
	},
	"views/text-editors-1": func(parent core.Widget) {
		texteditor.NewSoloEditor(parent).Buffer.SetTextString("Hello, world!")
	},
	"views/text-editors-2": func(parent core.Widget) {
		texteditor.NewSoloEditor(parent).Buffer.SetLang("go").SetTextString(`package main

func main() {
    fmt.Println("Hello, world!")
}
`)
	},
	"views/text-editors-3": func(parent core.Widget) {
		errors.Log(texteditor.NewSoloEditor(parent).Buffer.OpenFS(myFile, "file.go"))
	},
	"views/text-editors-4": func(parent core.Widget) {
		tb := texteditor.NewBuffer().SetTextString("Hello, world!")
		texteditor.NewEditor(parent).SetBuffer(tb)
		texteditor.NewEditor(parent).SetBuffer(tb)
	},
	"views/text-editors-5": func(parent core.Widget) {
		te := texteditor.NewSoloEditor(parent)
		te.OnInput(func(e events.Event) {
			core.MessageSnackbar(parent, "OnInput: "+te.Buffer.String())
		})
	},
	"views/trees-0": func(parent core.Widget) {
		tv := views.NewTree(parent).SetText("Root")
		views.NewTree(tv)
		c2 := views.NewTree(tv)
		views.NewTree(c2)
	},
	"views/trees-1": func(parent core.Widget) {
		n := tree.NewNodeBase()
		tree.NewNodeBase(n)
		c2 := tree.NewNodeBase(n)
		tree.NewNodeBase(c2)
		views.NewTree(parent).SyncTree(n)
	},
	"views/trees-2": func(parent core.Widget) {
		n := tree.NewNodeBase()
		tree.NewNodeBase(n)
		c2 := tree.NewNodeBase(n)
		tree.NewNodeBase(c2)
		views.NewTree(parent).SyncTree(n).OnChange(func(e events.Event) {
			core.MessageSnackbar(parent, "Tree changed")
		})
	},
	"views/trees-3": func(parent core.Widget) {
		n := tree.NewNodeBase()
		tree.NewNodeBase(n)
		c2 := tree.NewNodeBase(n)
		tree.NewNodeBase(c2)
		views.NewTree(parent).SyncTree(n).SetReadOnly(true)
	},
	"views/trees-4": func(parent core.Widget) {
		n := tree.NewNodeBase()
		tree.NewNodeBase(n)
		c2 := tree.NewNodeBase(n)
		tree.NewNodeBase(c2)
		// views.NewValue(parent, n)
	},
	"advanced/styling-0": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.OnWidgetAdded(func(w core.Widget) { // TODO(config)
			w.AsWidget().Styler(func(s *styles.Style) {
				s.Color = colors.C(colors.Scheme.Error.Base)
			})
		})
		core.NewText(fr).SetText("Label")
		core.NewSwitch(fr).SetText("Switch")
		core.NewTextField(fr).SetText("Text field")
	},
	"advanced/styling-1": func(parent core.Widget) {
		fr := core.NewFrame(parent)
		fr.OnWidgetAdded(func(w core.Widget) {
			switch w := w.(type) {
			case *core.Button:
				w.Styler(func(s *styles.Style) {
					s.Border.Radius = styles.BorderRadiusSmall
				})
			}
		})
		core.NewButton(fr).SetText("First")
		core.NewButton(fr).SetText("Second")
		core.NewButton(fr).SetText("Third")
	},
}
