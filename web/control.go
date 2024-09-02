package web

import "net/http"

type V4lControl struct {
	url        string
	icon       string
	multiplier int32
	visible    bool
	controls   []*V4lControl
	items      []MenuItem
}

var _ http.Handler = (*V4lControl)(nil)
var _ MenuItem = (*V4lControl)(nil)

func (handler *V4lControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func NewV4lControl(url string, icon string, multiplier int32, controls []*V4lControl) (ctrl *V4lControl) {
	ctrl = &V4lControl{
		url:        url,
		icon:       icon,
		multiplier: multiplier,
		controls:   controls,
		items:      make([]MenuItem, len(controls)),
	}
	for i, control := range controls {
		ctrl.items[i] = MenuItem(control)
	}
	return
}

func (ctrl *V4lControl) Url() string       { return ctrl.url }
func (ctrl *V4lControl) Icon() string      { return ctrl.icon }
func (ctrl *V4lControl) Visible() bool     { return ctrl.visible }
func (ctrl *V4lControl) Items() []MenuItem { return ctrl.items }
