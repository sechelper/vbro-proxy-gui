package ui

import "fyne.io/fyne/v2"

type ProxyGUI interface {
	view() *fyne.Container
	action()
	icon() *fyne.Resource
	name() string
}
