package web

import (
	"html/template"
	"net/http"
)

type MenuItem interface {
	http.Handler
	Url() string
	Items() []MenuItem
	Visible() bool
	Icon() string
}

type Menu struct {
	Tmpl    *template.Template
	Items   []MenuItem
	ItemMap map[string]MenuItem
}

func NewMenu(tmpl *template.Template, items ...MenuItem) *Menu {
	menu := &Menu{
		Tmpl:    tmpl,
		ItemMap: make(map[string]MenuItem),
	}

	menu.AddItems(items...)
	return menu
}

func (menu *Menu) AddItems(items ...MenuItem) {
	for _, item := range items {
		menu.ItemMap[item.Url()] = item
		v := item.Items()
		if len(v) > 0 {
			menu.AddItems(v...)
		}
	}
}
