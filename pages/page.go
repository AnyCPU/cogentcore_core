// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pages provides an easy way to make content-focused
// sites consisting of Markdown, HTML, and Cogent Core pages.
package pages

//go:generate core generate

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"net/url"
	"path"
	"strings"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/base/iox/tomlx"
	"cogentcore.org/core/base/strcase"
	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/htmlview"
	"cogentcore.org/core/pages/wpath"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/system"
	"cogentcore.org/core/tree"
	"cogentcore.org/core/views"
)

// Page represents a content page with support for navigating
// to other pages within the same source content.
type Page struct {
	core.Frame

	// Source is the filesystem in which the content is located.
	Source fs.FS

	// Context is the page's [htmlview.Context].
	Context *htmlview.Context `set:"-"`

	// The history of URLs that have been visited. The oldest page is first.
	History []string `set:"-"`

	// HistoryIndex is the current place we are at in the History
	HistoryIndex int `set:"-"`

	// PagePath is the fs path of the current page in [Page.Source]
	PagePath string `set:"-"`

	// URLToPagePath is a map between user-facing page URLs and underlying
	// FS page paths.
	URLToPagePath map[string]string `set:"-"`
}

var _ tree.Node = (*Page)(nil)

func (pg *Page) OnInit() {
	pg.Frame.OnInit()
	pg.Context = htmlview.NewContext()
	pg.Context.OpenURL = func(url string) {
		pg.OpenURL(url, true)
	}
	pg.Style(func(s *styles.Style) {
		s.Direction = styles.Column
	})
}

// OpenURL sets the content of the page from the given url. If the given URL
// has no scheme (eg: "/about"), then it sets the content of the page to the
// file specified by the URL. This is either the "index.md" file in the
// corresponding directory (eg: "/about/index.md") or the corresponding
// md file (eg: "/about.md"). If it has a scheme, (eg: "https://example.com"),
// then it opens it in the user's default browser.
func (pg *Page) OpenURL(rawURL string, addToHistory bool) {
	u, err := url.Parse(rawURL)
	if err != nil {
		core.ErrorSnackbar(pg, err, "Invalid URL")
		return
	}
	if u.Scheme != "" {
		system.TheApp.OpenURL(u.String())
		return
	}

	if pg.Source == nil {
		core.MessageSnackbar(pg, "Programmer error: page source must not be nil")
		return
	}

	// if we are not rooted, we go relative to our current fs path
	if !strings.HasPrefix(rawURL, "/") {
		rawURL = path.Join(path.Dir(pg.PagePath), rawURL)
	}

	// the paths in the fs are never rooted, so we trim a rooted one
	rawURL = strings.TrimPrefix(rawURL, "/")

	pg.PagePath = pg.URLToPagePath[rawURL]

	b, err := fs.ReadFile(pg.Source, pg.PagePath)
	if err != nil {
		// we go to the first page in the directory if there is no index page
		if errors.Is(err, fs.ErrNotExist) && (strings.HasSuffix(pg.PagePath, "index.md") || strings.HasSuffix(pg.PagePath, "index.html")) {
			fs.WalkDir(pg.Source, path.Dir(pg.PagePath), func(path string, d fs.DirEntry, err error) error {
				if path == pg.PagePath || d.IsDir() {
					return nil
				}
				pg.PagePath = path
				return fs.SkipAll
			})
			// need to update rawURL with new page path
			for u, p := range pg.URLToPagePath {
				if p == pg.PagePath {
					rawURL = u
					break
				}
			}
			b, err = fs.ReadFile(pg.Source, pg.PagePath)
		}
		if err != nil {
			core.ErrorSnackbar(pg, err, "Error opening page")
			return
		}
	}

	pg.Context.PageURL = rawURL
	if addToHistory {
		pg.HistoryIndex = len(pg.History)
		pg.History = append(pg.History, pg.Context.PageURL)
	}

	btp := []byte("+++")
	if bytes.HasPrefix(b, btp) {
		b = bytes.TrimPrefix(b, btp)
		fmb, content, ok := bytes.Cut(b, btp)
		if !ok {
			slog.Error("got unclosed front matter")
			b = fmb
			fmb = nil
		} else {
			b = content
		}
		if len(fmb) > 0 {
			var fm map[string]string
			errors.Log(tomlx.ReadBytes(&fm, fmb))
			fmt.Println("front matter", fm)
		}
	}

	// need to reset
	NumExamples[pg.Context.PageURL] = 0

	nav := pg.FindPath("splits/nav-frame/nav").(*views.TreeView)
	nav.UnselectAll()
	nav.FindPath(rawURL).(*views.TreeView).Select()

	fr := pg.FindPath("splits/body").(*core.Frame)
	fr.DeleteChildren()
	err = htmlview.ReadMD(pg.Context, fr, b)
	if err != nil {
		core.ErrorSnackbar(pg, err, "Error loading page")
		return
	}
	fr.Update()
}

func (pg *Page) Config() {
	if pg.HasChildren() {
		return
	}
	sp := core.NewSplits(pg, "splits").SetSplits(0.2, 0.8)

	nav := views.NewTreeViewFrame(sp, "nav").SetText(core.TheApp.Name())
	nav.SetReadOnly(true)
	nav.ParentWidget().Style(func(s *styles.Style) {
		s.Background = colors.C(colors.Scheme.SurfaceContainerLow)
	})
	nav.OnSelect(func(e events.Event) {
		if len(nav.SelectedNodes) == 0 {
			return
		}
		sn := nav.SelectedNodes[0]
		url := "/"
		if sn != nav {
			// we need a slash so that it doesn't think it's a relative URL
			url = "/" + sn.PathFrom(nav)
		}
		pg.OpenURL(url, true)
	})

	pg.URLToPagePath = map[string]string{"": "index.md"}

	errors.Log(fs.WalkDir(pg.Source, ".", func(fpath string, d fs.DirEntry, err error) error {
		// already handled
		if fpath == "" || fpath == "." {
			return nil
		}

		p := wpath.Format(fpath)

		pdir := path.Dir(p)
		base := path.Base(p)

		// already handled
		if base == "index.md" {
			return nil
		}

		ext := path.Ext(base)
		if ext != "" && ext != ".md" {
			return nil
		}

		parent := nav
		if pdir != "" && pdir != "." {
			parent = nav.FindPath(pdir).(*views.TreeView)
		}

		nm := strings.TrimSuffix(base, ext)
		txt := strcase.ToSentence(nm)
		tv := views.NewTreeView(parent, nm).SetText(txt)

		// need index.md for page path
		if d.IsDir() {
			fpath += "/index.md"
		}
		pg.URLToPagePath[tv.PathFrom(nav)] = fpath
		return nil
	}))

	core.NewFrame(sp, "body").Style(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	if pg.PagePath == "" {
		pg.OpenURL("/", true)
	}
}

// AppBar is the default app bar for a [Page]
func (pg *Page) AppBar(tb *core.Toolbar) {
	// ch := tb.AppChooser()

	back := tb.ChildByName("back").(*core.Button)
	back.OnClick(func(e events.Event) {
		if pg.HistoryIndex > 0 {
			pg.HistoryIndex--
			// we reverse the order
			// ch.SelectItem(len(pg.History) - pg.HistoryIndex - 1)
			// we need a slash so that it doesn't think it's a relative URL
			pg.OpenURL("/"+pg.History[pg.HistoryIndex], false)
		}
	})

	// TODO(kai/abc)
	// ch.AddItemsFunc(func() {
	// 	ch.Items = make([]any, len(pg.History))
	// 	for i, u := range pg.History {
	// 		// we reverse the order
	// 		ch.Items[len(pg.History)-i-1] = u
	// 	}
	// })
	// ch.OnChange(func(e events.Event) {
	// 	// we need a slash so that it doesn't think it's a relative URL
	// 	pg.OpenURL("/"+ch.CurrentLabel, true)
	// 	e.SetHandled()
	// })
}
