// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giv

import (
	"log/slog"
	"strconv"
	"time"

	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/girl/styles"
	"goki.dev/girl/units"
	"goki.dev/goosi/events"
	"goki.dev/gti"
	"goki.dev/icons"
	"goki.dev/laser"
	"goki.dev/pi/v2/filecat"
)

// TimeView is a view for selecting a time
type TimeView struct {
	gi.Frame

	// the time that we are viewing
	Time time.Time `set:"-"`

	// value view that needs to have SaveTmp called on it whenever a change is made to one of the underlying values -- pass this down to any sub-views created from a parent
	TmpSave Value `json:"-" xml:"-"`

	// a record of parent View names that have led up to this view -- displayed as extra contextual information in view dialog windows
	ViewPath string

	// the raw input hour
	Hour int `set:"-"`

	// whether we are in PM mode (so we have to add 12h to everything)
	PM bool `set:"-"`
}

// SetTime sets the source time and updates the view
func (tv *TimeView) SetTime(tim time.Time) *TimeView {
	updt := tv.UpdateStart()
	tv.Time = tim
	tv.UpdateEndRender(updt)
	tv.SendChange()
	return tv
}

func (tv *TimeView) ConfigWidget(sc *gi.Scene) {
	if tv.HasChildren() {
		return
	}
	updt := tv.UpdateStart()

	tv.SetLayout(gi.LayoutHoriz)

	hour := gi.NewTextField(tv, "hour")
	if gi.Prefs.Clock24 {
		tv.Hour = tv.Time.Hour()
		hour.SetText(strconv.Itoa(tv.Hour))
	} else {
		tv.Hour = tv.Time.Hour() % 12
		hour.SetText(strconv.Itoa(tv.Hour))
	}
	hour.Style(func(s *styles.Style) {
		s.Font.Size.Dp(57)
		s.SetFixedWidth(units.Dp(96))
	})
	hour.OnChange(func(e events.Event) {
		hr, err := strconv.Atoi(hour.Text())
		// TODO(kai/snack)
		if err != nil {
			slog.Error(err.Error())
		}
		tv.Hour = hr
		if tv.PM {
			// only add to local
			hr += 12
		}
		// we set our hour and keep everything else
		tt := tv.Time
		tv.Time = time.Date(tt.Year(), tt.Month(), tt.Day(), hr, tt.Minute(), tt.Second(), tt.Nanosecond(), tt.Location())
		if tv.TmpSave != nil {
			tv.TmpSave.SetValue(tv.Time)
		}
	})

	gi.NewLabel(tv, "colon").SetType(gi.LabelDisplayLarge).SetText(":")

	minute := gi.NewTextField(tv, "minute").
		SetText(strconv.Itoa(tv.Time.Minute()))
	minute.Style(func(s *styles.Style) {
		s.Font.Size.Dp(57)
		s.SetFixedWidth(units.Dp(96))
	})
	minute.OnChange(func(e events.Event) {
		min, err := strconv.Atoi(minute.Text())
		// TODO(kai/snack)
		if err != nil {
			slog.Error(err.Error())
		}
		// we set our minute and keep everything else
		tt := tv.Time
		tv.Time = time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), min, tt.Second(), tt.Nanosecond(), tt.Location())
		if tv.TmpSave != nil {
			tv.TmpSave.SetValue(tv.Time)
		}
	})

	if !gi.Prefs.Clock24 {
		sw := gi.NewSwitches(tv, "am-pm").SetMutex(true).SetType(gi.SwitchSegmentedButton).SetLayout(gi.LayoutVert).SetItems([]string{"AM", "PM"})
		if tv.Time.Hour() < 12 {
			tv.PM = false
			sw.SelectItemAction(0)
		} else {
			tv.PM = true
			sw.SelectItemAction(1)
		}
		sw.OnChange(func(e events.Event) {
			si := sw.SelectedItem()
			tt := tv.Time
			switch si {
			case "AM":
				tv.PM = false
				tv.Time = time.Date(tt.Year(), tt.Month(), tt.Day(), tv.Hour, tt.Minute(), tt.Second(), tt.Nanosecond(), tt.Location())
			case "PM":
				tv.PM = true
				tv.Time = time.Date(tt.Year(), tt.Month(), tt.Day(), tv.Hour+12, tt.Minute(), tt.Second(), tt.Nanosecond(), tt.Location())
			default:
				// must always have something valid selected
				tv.PM = false
				sw.SelectItem(0)
				tv.Time = time.Date(tt.Year(), tt.Month(), tt.Day(), tv.Hour, tt.Minute(), tt.Second(), tt.Nanosecond(), tt.Location())
			}
			if tv.TmpSave != nil {
				tv.TmpSave.SetValue(tv.Time)
			}
		})
	}

	tv.UpdateEnd(updt)
}

var shortMonths = []string{"Jan", "Feb", "Apr", "Mar", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// DateView is a view for selecting a date
type DateView struct {
	gi.Frame

	// the time that we are viewing
	Time time.Time `set:"-"`

	// value view that needs to have SaveTmp called on it whenever a change is made to one of the underlying values -- pass this down to any sub-views created from a parent
	TmpSave Value `json:"-" xml:"-"`

	// a record of parent View names that have led up to this view -- displayed as extra contextual information in view dialog windows
	ViewPath string
}

// SetTime sets the source time and updates the view
func (tv *DateView) SetTime(tim time.Time) *DateView {
	updt := tv.UpdateStart()
	tv.Time = tim
	tv.UpdateEndRender(updt)
	tv.SendChange()
	return tv
}

func (tv *DateView) ConfigWidget(sc *gi.Scene) {
	if tv.HasChildren() {
		return
	}
	updt := tv.UpdateStart()

	tv.SetLayout(gi.LayoutVert)

	trow := gi.NewLayout(tv).SetLayout(gi.LayoutHoriz)

	sms := make([]any, len(shortMonths))
	for i, sm := range shortMonths {
		sms[i] = sm
	}
	gi.NewChooser(trow, "month").SetItems(sms).SetCurIndex(int(tv.Time.Month() - 1))

	yr := tv.Time.Year()
	yrs := []any{}
	// we go 100 in each direction from the current year
	for i := yr - 100; i <= yr+100; i++ {
		yrs = append(yrs, i)
	}
	gi.NewChooser(trow, "year").SetItems(yrs).SetCurVal(yr)

	grid := gi.NewLayout(tv, "grid").SetLayout(gi.LayoutGrid)
	grid.Style(func(s *styles.Style) {
		s.Columns = 7
	})

	// start of the month
	som := tv.Time.AddDate(0, 0, -tv.Time.Day()+1)
	// end of the month
	eom := tv.Time.AddDate(0, 1, -tv.Time.Day())
	// start of the week containing the start of the month
	somw := som.AddDate(0, 0, -int(som.Weekday()))
	// end of the week containing the end of the month
	eomw := eom.AddDate(0, 0, int(7-eom.Weekday()))

	for yd := somw.YearDay(); yd <= eomw.YearDay(); yd++ {
		yd := yd
		yds := strconv.Itoa(yd)
		// actual time of this date
		dt := somw.AddDate(0, 0, yd-somw.YearDay())
		ds := strconv.Itoa(dt.Day())
		lb := gi.NewLabel(grid, "day-"+yds).SetText(ds)
		lb.Style(func(s *styles.Style) {
			s.Padding.Set(units.Dp(4))
			s.Border.Radius = styles.BorderRadiusFull
			s.AlignV = styles.AlignCenter
			s.Text.Align = styles.AlignCenter
			if dt.Month() != som.Month() {
				s.Color = colors.Scheme.OnSurfaceVariant
			}
			if yd == time.Now().YearDay() {
				s.Border.Width.Set(units.Dp(1))
				s.Border.Color.Set(colors.Scheme.Primary.Base)
			}
			if yd == tv.Time.YearDay() {
				s.BackgroundColor.SetSolid(colors.Scheme.Primary.Base)
				s.Color = colors.Scheme.Primary.On
			}
		})
	}

	tv.UpdateEnd(updt)
}

// TimeValue presents two text fields for editing a date and time,
// both of which can pull up corresponding picker view dialogs.
type TimeValue struct {
	ValueBase
}

func (vv *TimeValue) WidgetType() *gti.Type {
	vv.WidgetTyp = gi.LayoutType
	return vv.WidgetTyp
}

// TimeVal decodes Value into a *time.Time value -- also handles FileTime case
func (vv *TimeValue) TimeVal() *time.Time {
	tmi := laser.PtrValue(vv.Value).Interface()
	switch v := tmi.(type) {
	case *time.Time:
		return v
	case *filecat.FileTime:
		return (*time.Time)(v)
	}
	return nil
}

func (vv *TimeValue) UpdateWidget() {
	if vv.Widget == nil {
		return
	}
	fr := vv.Widget.(*gi.Layout)
	tm := vv.TimeVal()

	fr.ChildByName("date").(*gi.TextField).SetText(tm.Format("1/2/2006"))
	fr.ChildByName("time").(*gi.TextField).SetText(tm.Format(gi.Prefs.TimeFormat()))
}

func (vv *TimeValue) ConfigWidget(w gi.Widget, sc *gi.Scene) {
	if vv.Widget == w {
		vv.UpdateWidget()
		return
	}
	// need TmpSave
	if vv.TmpSave == nil {
		tt := vv.TimeVal()
		vv.TmpSave = NewSoloValue(tt)
	}
	vv.Widget = w
	vv.StdConfigWidget(w)
	ly := vv.Widget.(*gi.Layout)
	ly.SetLayout(gi.LayoutHoriz)

	if len(ly.Kids) > 0 {
		return
	}

	dt := gi.NewTextField(ly, "date").SetTooltip("The date").
		SetLeadingIcon(icons.CalendarToday, func(e events.Event) {
			d := gi.NewDialog(w).Title("Select date")
			NewDateView(d).SetTime(*vv.TimeVal()).SetTmpSave(vv.TmpSave)
			d.OnAccept(func(e events.Event) {
				tt := vv.TmpSave.Val().Interface().(*time.Time)
				vv.SetValue(tt)
				vv.UpdateWidget()
			}).Cancel().Ok().Run()
		})
	dt.SetMinPrefWidth(units.Em(8))
	dt.SetReadOnly(vv.IsReadOnly())
	dt.OnChange(func(e events.Event) {
		d, err := time.Parse("01/02/2006", dt.Text())
		if err != err {
			// TODO(kai/snack)
			slog.Error(err.Error())
			return
		}
		tv := vv.TimeVal()
		// new date and old time
		*tv = time.Date(d.Year(), d.Month(), d.Day(), tv.Hour(), tv.Minute(), tv.Second(), tv.Nanosecond(), tv.Location())
	})
	dt.Config(sc)

	tm := gi.NewTextField(ly, "time").SetTooltip("The time").
		SetLeadingIcon(icons.Schedule, func(e events.Event) {
			d := gi.NewDialog(w).Title("Edit time")
			NewTimeView(d).SetTime(*vv.TimeVal()).SetTmpSave(vv.TmpSave)
			d.OnAccept(func(e events.Event) {
				tt := vv.TmpSave.Val().Interface().(*time.Time)
				vv.SetValue(tt)
				vv.UpdateWidget()
			}).Cancel().Ok().Run()
		})
	tm.SetMinPrefWidth(units.Em(8))
	tm.SetReadOnly(vv.IsReadOnly())
	tm.OnChange(func(e events.Event) {
		t, err := time.Parse(gi.Prefs.TimeFormat(), tm.Text())
		if err != err {
			// TODO(kai/snack)
			slog.Error(err.Error())
			return
		}
		tv := vv.TimeVal()
		// old date and new time
		*tv = time.Date(tv.Year(), tv.Month(), tv.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tv.Location())
	})
	dt.Config(sc)

	vv.UpdateWidget()
}

var durationUnits = []string{
	"nanoseconds",
	"microseconds",
	"milliseconds",
	"seconds",
	"minutes",
	"hours",
	"days",
	"weeks",
	"months",
	"years",
}

var durationUnitsMap = map[string]time.Duration{
	"nanoseconds":  time.Nanosecond,
	"microseconds": time.Microsecond,
	"milliseconds": time.Millisecond,
	"seconds":      time.Second,
	"minutes":      time.Minute,
	"hours":        time.Hour,
	"days":         24 * time.Hour,
	"weeks":        7 * 24 * time.Hour,
	"months":       30 * 24 * time.Hour,
	"years":        365 * 24 * time.Hour,
}

// DurationValue presents a spinner and unit chooser for a [time.Duration]
type DurationValue struct {
	ValueBase
}

func (vv *DurationValue) WidgetType() *gti.Type {
	vv.WidgetTyp = gi.LayoutType
	return vv.WidgetTyp
}

func (vv *DurationValue) UpdateWidget() {
	if vv.Widget == nil {
		return
	}
	npv := laser.NonPtrValue(vv.Value)
	dur := npv.Interface().(time.Duration)
	un := "seconds"
	undur := time.Duration(0)
	for _, u := range durationUnits {
		v := durationUnitsMap[u]
		if v > dur {
			break
		}
		un = u
		undur = v
	}
	adur := dur
	if undur != 0 {
		adur = dur / undur
	}

	ly := vv.Widget.(*gi.Layout)
	ly.ChildByName("value").(*gi.Spinner).SetValue(float32(adur))
	ly.ChildByName("unit").(*gi.Chooser).SetCurVal(un)
}

func (vv *DurationValue) ConfigWidget(w gi.Widget, sc *gi.Scene) {
	if vv.Widget == w {
		vv.UpdateWidget()
		return
	}
	vv.Widget = w
	vv.StdConfigWidget(w)
	fr := vv.Widget.(*gi.Layout)
	fr.SetLayout(gi.LayoutHoriz)

	if len(fr.Kids) > 0 {
		return
	}

	var ch *gi.Chooser

	sp := gi.NewSpinner(fr, "value").SetTooltip("The value of time").SetStep(1).SetPageStep(10)
	sp.OnChange(func(e events.Event) {
		vv.SetValue(sp.Value * float32(durationUnitsMap[ch.CurLabel]))
	})
	sp.Config(sc)

	units := []any{}
	for _, u := range durationUnits {
		units = append(units, u)
	}

	ch = gi.NewChooser(fr, "unit").SetTooltip("The unit of time").SetItems(units)
	ch.OnChange(func(e events.Event) {
		// we update the value to fit the unit
		npv := laser.NonPtrValue(vv.Value)
		dur := npv.Interface().(time.Duration)
		sp.SetValue(float32(dur) / float32(durationUnitsMap[ch.CurLabel]))
	})

	vv.UpdateWidget()
}
