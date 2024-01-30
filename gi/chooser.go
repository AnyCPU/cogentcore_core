// Copyright (c) 2018, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

import (
	"fmt"
	"log/slog"
	"strconv"

	"cogentcore.org/core/abilities"
	"cogentcore.org/core/colors"
	"cogentcore.org/core/cursors"
	"cogentcore.org/core/enums"
	"cogentcore.org/core/events"
	"cogentcore.org/core/glop/sentence"
	"cogentcore.org/core/gti"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/keyfun"
	"cogentcore.org/core/ki"
	"cogentcore.org/core/pi/complete"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/units"
)

// Chooser is for selecting items from a dropdown list, with an optional
// edit TextField for typing directly.
// The items can be of any type, including enum values -- they are converted
// to strings for the display.  If the items are of type [icons.Icon], then they
// are displayed using icons instead.
type Chooser struct {
	Box

	// Type is the styling type of the chooser.
	Type ChooserTypes

	// Items are the chooser items available for selection.
	Items []ChooserItem

	// Icon is an optional icon displayed on the left side of the chooser.
	Icon icons.Icon `view:"show-name"`

	// Indicator is the icon to use for the indicator displayed on the
	// right side of the chooser.
	Indicator icons.Icon `view:"show-name"`

	// Editable is whether provide a text field for editing the value,
	// or just a button for selecting items.
	Editable bool

	// TODO(kai): implement AllowNew button

	// AllowNew is whether to allow the user to add new items to the
	// chooser through the editable textfield (if Editable is set to
	// true) and a button at the end of the chooser menu.
	AllowNew bool

	// Placeholder, if Editable is set to true, is the text that is
	// displayed in the text field when it is empty. It must be set
	// using [Chooser.SetPlaceholder].
	Placeholder string `set:"-"`

	// ItemsFuncs is a slice of functions to call before showing the items
	// of the chooser, which is typically used to configure them
	// (eg: if they are based on dynamic data). The functions are called
	// in ascending order such that the items added in the first function
	// will appear before those added in the last function. Use
	// [Chooser.AddItemsFunc] to add a new items function. If at least
	// one ItemsFunc is specified, the items of the chooser will be
	// cleared before calling the functions.
	ItemsFuncs []func() `copier:"-" set:"-"`

	// CurrentItem is the currently selected item.
	CurrentItem ChooserItem `json:"-" xml:"-" set:"-"`

	// CurrentIndex is the index of the currently selected item
	// in [Chooser.Items].
	CurrentIndex int `json:"-" xml:"-" set:"-"`
}

// ChooserItem is an item that can be used in a [Chooser].
type ChooserItem struct {

	// Value is the underlying value of the chooser item.
	Value any

	// Func, if non-nil, is a function to call whenever this
	// item is selected as the current value of the chooser.
	Func func()

	// Label is the label displayed to the user for this item.
	// If it is empty, then [ToLabel] of [ChooserItem.Value] is
	// used instead.
	Label string

	// Icon is the icon displayed to the user for this item.
	Icon icons.Icon

	// Tooltip is the tooltip displayed to the user for this item.
	Tooltip string
}

// GetLabel returns the effective label of this chooser item.
// If [ChooserItem.Label] is set, it returns that. Otherwise,
// it uses [ToLabel] of [ChooserItem.Value]
func (ci *ChooserItem) GetLabel() string {
	if ci.Label != "" {
		return ci.Label
	}
	if ci.Value == nil {
		return ""
	}
	return ToLabel(ci.Value)
}

// ChooserTypes is an enum containing the
// different possible types of combo boxes
type ChooserTypes int32 //enums:enum -trim-prefix Chooser

const (
	// ChooserFilled represents a filled
	// Chooser with a background color
	// and a bottom border
	ChooserFilled ChooserTypes = iota
	// ChooserOutlined represents an outlined
	// Chooser with a border on all sides
	// and no background color
	ChooserOutlined
)

func (ch *Chooser) OnInit() {
	ch.Box.OnInit()
	ch.HandleEvents()
	ch.SetStyles()
}

func (ch *Chooser) SetStyles() {
	ch.Icon = icons.None
	ch.Indicator = icons.KeyboardArrowDown
	ch.Style(func(s *styles.Style) {
		s.SetAbilities(true, abilities.Activatable, abilities.Hoverable, abilities.LongHoverable)
		if !ch.Editable {
			s.SetAbilities(true, abilities.Focusable)
		}
		s.Cursor = cursors.Pointer
		s.Text.Align = styles.Center
		s.Border.Radius = styles.BorderRadiusSmall
		s.Padding.Set(units.Dp(8), units.Dp(16))
		switch ch.Type {
		case ChooserFilled:
			s.Background = colors.C(colors.Scheme.Secondary.Container)
			s.Color = colors.Scheme.Secondary.OnContainer
		case ChooserOutlined:
			s.Border.Style.Set(styles.BorderSolid)
			s.Border.Width.Set(units.Dp(1))
			s.Border.Color.Set(colors.Scheme.OnSurfaceVariant)
		}
		// textfield handles everything
		if ch.Editable {
			s.Border = styles.Border{}
			s.MaxBorder = s.Border
			s.Background = nil
			s.StateLayer = 0
			s.Padding.Zero()
			if ch.Type == ChooserFilled {
				s.Border.Radius = styles.BorderRadiusExtraSmallTop
			} else {
				s.Border.Radius = styles.BorderRadiusExtraSmall
			}
		}
	})
	ch.OnWidgetAdded(func(w Widget) {
		switch w.PathFrom(ch) {
		case "parts":
			w.Style(func(s *styles.Style) {
				s.Align.Content = styles.Center
				s.Align.Items = styles.Center
			})
		case "parts/icon":
			w.Style(func(s *styles.Style) {
				s.Margin.Zero()
				s.Padding.Zero()
				s.Margin.Right.Ch(1)
			})
		case "parts/label":
			w.Style(func(s *styles.Style) {
				s.SetNonSelectable()
				s.SetTextWrap(false)
				s.Margin.Zero()
				s.Padding.Zero()
			})
		case "parts/text":
			text := w.(*TextField)
			text.Placeholder = ch.Placeholder
			if ch.Type == ChooserFilled {
				text.Type = TextFieldFilled
			} else {
				text.Type = TextFieldOutlined
			}
			ch.HandleChooserTextFieldEvents(text)
			text.Style(func(s *styles.Style) {
				s.Grow = ch.Styles.Grow // we grow like our parent
			})
		case "parts/indicator":
			w.Style(func(s *styles.Style) {
				s.Font.Size.Dp(16)
				s.Min.X.Em(1)
				s.Min.Y.Em(1)
				s.Justify.Self = styles.End
				s.Align.Self = styles.Center
			})
		}
	})
}

func (ch *Chooser) ConfigWidget() {
	config := ki.Config{}

	ici := -1
	var lbi, txi, indi int
	// editable handles through textfield
	if ch.Icon.IsSet() && !ch.Editable {
		config.Add(IconType, "icon")
		ici = 0
	}
	if ch.Editable {
		lbi = -1
		txi = len(config)
		config.Add(TextFieldType, "text")
	} else {
		txi = -1
		lbi = len(config)
		config.Add(LabelType, "label")
	}
	if !ch.Indicator.IsSet() {
		ch.Indicator = icons.KeyboardArrowRight
	}
	// editable handles through textfield
	if !ch.Editable {
		indi = len(config)
		config.Add(IconType, "indicator")
	}

	ch.ConfigParts(config, func() {
		if ici >= 0 {
			ic := ch.Parts.Child(ici).(*Icon)
			ic.SetIcon(ch.Icon)
		}
		if ch.Editable {
			tx := ch.Parts.Child(txi).(*TextField)
			tx.SetText(ch.CurrentItem.GetLabel())
			tx.SetLeadingIcon(ch.Icon)
			tx.SetTrailingIcon(ch.Indicator, func(e events.Event) {
				ch.OpenMenu(e)
			})
			tx.Config() // this is essential
			tx.SetCompleter(tx, ch.CompleteMatch, ch.CompleteEdit)
		} else {
			lbl := ch.Parts.Child(lbi).(*Label)
			lbl.SetText(ch.CurrentItem.GetLabel())
			lbl.Config() // this is essential

			ic := ch.Parts.Child(indi).(*Icon)
			ic.SetIcon(ch.Indicator)
		}
	})
}

// LabelWidget returns the label widget if present
func (ch *Chooser) LabelWidget() *Label {
	if ch.Parts == nil {
		return nil
	}
	lbi := ch.Parts.ChildByName("label")
	if lbi == nil {
		return nil
	}
	return lbi.(*Label)
}

// IconWidget returns the icon widget if present
func (ch *Chooser) IconWidget() *Icon {
	if ch.Parts == nil {
		return nil
	}
	ici := ch.Parts.ChildByName("icon")
	if ici == nil {
		return nil
	}
	return ici.(*Icon)
}

// SetIconUpdate sets the icon and drives an update, for the already-displayed case.
func (ch *Chooser) SetIconUpdate(ic icons.Icon) *Chooser {
	updt := ch.UpdateStart()
	defer ch.UpdateEndRender(updt)

	ch.Icon = ic
	if ch.Editable {
		tf := ch.TextField()
		if tf != nil {
			tf.SetLeadingIconUpdate(ic)
		}
	} else {
		iw := ch.IconWidget()
		if iw != nil {
			iw.SetIconUpdate(ic)
		}
	}
	return ch
}

// TextField returns the text field of an editable Chooser
// if present
func (ch *Chooser) TextField() *TextField {
	if ch.Parts == nil {
		return nil
	}
	tf := ch.Parts.ChildByName("text", 2)
	if tf == nil {
		return nil
	}
	return tf.(*TextField)
}

// AddItemsFunc adds the given function to [Chooser.ItemsFuncs].
// These functions are called before showing the items of the chooser,
// and they are typically used to configure them (eg: if they are based
// on dynamic data). The functions are called in ascending order such
// that the items added in the first function will appear before those
// added in the last function. If at least one ItemsFunc is specified,
// the items, labels, icons, and tooltips of the chooser will be cleared
// before calling the functions.
func (ch *Chooser) AddItemsFunc(f func()) *Chooser {
	ch.ItemsFuncs = append(ch.ItemsFuncs, f)
	return ch
}

// CallItemsFuncs calls [Chooser.ItemsFuncs].
func (ch *Chooser) CallItemsFuncs() {
	if len(ch.ItemsFuncs) == 0 {
		return
	}
	ch.Items = nil
	for _, f := range ch.ItemsFuncs {
		f()
	}
}

// SetTypes sets the [Chooser.Items] from the given types.
// If setFirst is true, it sets the current item to the first item
// in the list.
func (ch *Chooser) SetTypes(ts []*gti.Type, setFirst bool) *Chooser {
	ch.Items = make([]ChooserItem, len(ts))
	for i, typ := range ts {
		ch.Items[i] = ChooserItem{Value: typ}
	}
	if setFirst {
		ch.SetCurIndex(0)
	}
	return ch
}

// SetStrings sets the [Chooser.Items] from the given strings.
// If setFirst is true, it sets the current item to the first item
// in the list.
func (ch *Chooser) SetStrings(ss []string, setFirst bool) *Chooser {
	ch.Items = make([]ChooserItem, len(ss))
	for i, s := range ss {
		ch.Items[i] = ChooserItem{Value: s}
	}
	if setFirst {
		ch.SetCurIndex(0)
	}
	return ch
}

// SetIconItems sets the [Chooser.Items] from the given icons.
// If setFirst is true, it sets the current item to the first item
// in the list.
func (ch *Chooser) SetIconItems(is []icons.Icon, setFirst bool) *Chooser {
	ch.Items = make([]ChooserItem, len(is))
	for i, ic := range is {
		ch.Items[i] = ChooserItem{Value: ic, Label: sentence.Case(string(ic)), Icon: ic}
	}
	if setFirst {
		ch.SetCurIndex(0)
	}
	return ch
}

// SetEnums sets the [Chooser.Items] from the given enums.
// If setFirst is true, it sets the current item to the first item
// in the list.
func (ch *Chooser) SetEnums(es []enums.Enum, setFirst bool) *Chooser {
	ch.Items = make([]ChooserItem, len(es))
	for i, enum := range es {
		str := enum.String()
		lbl := sentence.Case(str)
		// TODO(kai): this desc is not always correct because we
		// don't have the name of the enum value pre-generator-transformation
		// (same as with Switches) (#774)
		tip := sentence.Doc(enum.Desc(), str, lbl)
		ch.Items[i] = ChooserItem{Value: enum, Label: lbl, Tooltip: tip}
	}
	if setFirst {
		ch.SetCurIndex(0)
	}
	return ch
}

// SetEnum sets the [Chooser.Items] from [enums.Enum.Values] of the given enum.
// If setFirst is true, it sets the current item to the first item
// in the list.
func (ch *Chooser) SetEnum(enum enums.Enum, setFirst bool) *Chooser {
	return ch.SetEnums(enum.Values(), setFirst)
}

// FindItem finds the given item value on the list of items and returns its index
func (ch *Chooser) FindItem(it any) int {
	for i, v := range ch.Items {
		if it == v.Value {
			return i
		}
	}
	return -1
}

// SetPlaceholder sets the given placeholder text and
// indicates that nothing has been selected.
func (ch *Chooser) SetPlaceholder(text string) *Chooser {
	ch.Placeholder = text
	if !ch.Editable {
		ch.CurrentItem.Label = text
		ch.ShowCurVal()
	}
	ch.CurrentIndex = -1
	return ch
}

// SetCurVal sets the current item and index to those associated with the given value.
// If the given item is not found, it adds it to the items list. It also sets the text
// of the chooser to the label of the item.
func (ch *Chooser) SetCurVal(it any) *Chooser {
	ch.CurrentIndex = ch.FindItem(it)
	if ch.CurrentIndex < 0 { // add to list if not found
		ch.CurrentIndex = len(ch.Items)
		ch.Items = append(ch.Items, ChooserItem{Value: it})
	}
	ch.CurrentItem = ch.Items[ch.CurrentIndex]
	ch.ShowCurVal()
	return ch
}

// SetCurIndex sets the current index and the item associated with it.
func (ch *Chooser) SetCurIndex(idx int) *Chooser {
	ch.CurrentIndex = idx
	ch.CurrentItem = ch.Items[idx]
	ch.ShowCurVal()
	return ch
}

// SetCurText sets the current index and item based on the given text string.
// It can only be used for editable choosers.
func (ch *Chooser) SetCurText(text string) *Chooser {
	for i, item := range ch.Items {
		if text == item.GetLabel() {
			ch.SetCurIndex(i)
			return ch
		}
	}
	if !ch.AllowNew {
		// TODO: use validation
		slog.Error("invalid Chooser value", "value", text)
		return ch
	}
	ch.Items = append(ch.Items, ChooserItem{Value: text})
	ch.SetCurIndex(len(ch.Items) - 1)
	return ch
}

// ShowCurVal updates the display to present the current item.
func (ch *Chooser) ShowCurVal() *Chooser {
	updt := ch.UpdateStart()
	defer ch.UpdateEndRender(updt)

	if ch.Editable {
		tf := ch.TextField()
		if tf != nil {
			tf.SetTextUpdate(ch.CurrentItem.GetLabel())
		}
	} else {
		lbl := ch.LabelWidget()
		if lbl != nil {
			lbl.SetTextUpdate(ch.CurrentItem.GetLabel())
		}
	}
	if ch.CurrentItem.Icon.IsSet() {
		picon := ch.Icon
		ch.SetIcon(ch.CurrentItem.Icon)
		if ch.Icon != picon {
			ch.Update()
			ch.SetNeedsLayout(true)
		}
	}
	if ch.CurrentItem.Tooltip != "" {
		ch.SetTooltip(ch.CurrentItem.Tooltip)
	}
	return ch
}

// SelectItem selects the item at the given index and updates the chooser to display it.
func (ch *Chooser) SelectItem(idx int) *Chooser {
	if ch.This() == nil {
		return ch
	}
	updt := ch.UpdateStart()
	ch.SetCurIndex(idx)
	ch.UpdateEndLayout(updt)
	return ch
}

// SelectItemAction selects the item at the given index and updates the chooser to display it.
// It also sends an [events.Change] event to indicate that the value has changed.
func (ch *Chooser) SelectItemAction(idx int) *Chooser {
	if ch.This() == nil {
		return ch
	}
	ch.SelectItem(idx)
	ch.SendChange()
	return ch
}

// MakeItemsMenu constructs a menu of all the items. It is used when the chooser is clicked.
func (ch *Chooser) MakeItemsMenu(m *Scene) {
	ch.CallItemsFuncs()
	for i, it := range ch.Items {
		i := i
		nm := "item-" + strconv.Itoa(i)
		bt := NewButton(m, nm)
		bt.SetText(it.GetLabel()).SetIcon(it.Icon).SetTooltip(it.Tooltip)
		bt.SetSelected(i == ch.CurrentIndex)
		bt.OnClick(func(e events.Event) {
			ch.SelectItemAction(i)
		})
	}
}

func (ch *Chooser) HandleEvents() {
	ch.HandleSelectToggle()

	ch.OnClick(func(e events.Event) {
		if ch.OpenMenu(e) {
			e.SetHandled()
		}
	})
	ch.OnFinal(events.KeyChord, func(e events.Event) {
		if DebugSettings.KeyEventTrace {
			fmt.Printf("Chooser KeyChordEvent: %v\n", ch.Path())
		}
		kf := keyfun.Of(e.KeyChord())
		switch {
		case kf == keyfun.MoveUp:
			e.SetHandled()
			if len(ch.Items) > 0 {
				idx := ch.CurrentIndex - 1
				if idx < 0 {
					idx += len(ch.Items)
				}
				ch.SelectItemAction(idx)
			}
		case kf == keyfun.MoveDown:
			e.SetHandled()
			if len(ch.Items) > 0 {
				idx := ch.CurrentIndex + 1
				if idx >= len(ch.Items) {
					idx -= len(ch.Items)
				}
				ch.SelectItemAction(idx)
			}
		case kf == keyfun.PageUp:
			e.SetHandled()
			if len(ch.Items) > 10 {
				idx := ch.CurrentIndex - 10
				for idx < 0 {
					idx += len(ch.Items)
				}
				ch.SelectItemAction(idx)
			}
		case kf == keyfun.PageDown:
			e.SetHandled()
			if len(ch.Items) > 10 {
				idx := ch.CurrentIndex + 10
				for idx >= len(ch.Items) {
					idx -= len(ch.Items)
				}
				ch.SelectItemAction(idx)
			}
		case kf == keyfun.Enter || (!ch.Editable && e.KeyRune() == ' '):
			// if !(kt.Rune == ' ' && chb.Sc.Type == ScCompleter) {
			e.SetHandled()
			ch.Send(events.Click, e)
		// }
		default:
			tf := ch.TextField()
			if tf == nil {
				break
			}
			// if we don't have anything special to do,
			// we just give our key event to our textfield
			tf.HandleEvent(e)
		}
	})
}

// OpenMenu opens the chooser menu that displays all of the items.
// It returns false if there are no items.
func (ch *Chooser) OpenMenu(e events.Event) bool {
	pos := ch.ContextMenuPos(e)
	if ch.Parts != nil {
		if indic := ch.Parts.ChildByName("indicator", 3); indic != nil {
			pos = indic.(Widget).ContextMenuPos(nil) // use the pos
		}
	}
	m := NewMenu(ch.MakeItemsMenu, ch.This().(Widget), pos)
	if m == nil {
		return false
	}
	m.Run()
	return true
}

func (ch *Chooser) HandleChooserTextFieldEvents(tf *TextField) {
	tf.OnChange(func(e events.Event) {
		ch.SetCurText(tf.Text())
		ch.SendChange(e)
	})
	tf.OnFocus(func(e events.Event) {
		if ch.IsReadOnly() {
			return
		}
		ch.CallItemsFuncs()
	})
	tf.OnClick(func(e events.Event) {
		tf.OfferComplete(dontForce)
	})
}

// CompleteMatch is the [complete.MatchFunc] used for the
// editable text field part of the Chooser (if it exists).
func (ch *Chooser) CompleteMatch(data any, text string, posLn, posCh int) (md complete.Matches) {
	md.Seed = text
	comps := make(complete.Completions, len(ch.Items))
	for i, item := range ch.Items {
		comps[i] = complete.Completion{
			Text: item.GetLabel(),
			Desc: item.Tooltip,
			Icon: string(item.Icon),
		}
	}
	md.Matches = complete.MatchSeedCompletion(comps, md.Seed)
	return md
}

// CompleteEdit is the [complete.EditFunc] used for the
// editable textfield part of the Chooser (if it exists).
func (ch *Chooser) CompleteEdit(data any, text string, cursorPos int, completion complete.Completion, seed string) (ed complete.Edit) {
	return complete.Edit{
		NewText:       completion.Text,
		ForwardDelete: len([]rune(text)),
	}
}
