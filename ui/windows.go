package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var (
	globalMainWin  *MainWin
	globalMainTabs *container.DocTabs
)

type MainWin struct {
	app           fyne.App
	win           fyne.Window
	menu          *ProxyMainMenu
	intercepterUI *intercepterUI
	replayUI      replayUI
}

func NewMainWin() *MainWin {

	mainWin := &MainWin{
		app:           app.New(),
		menu:          NewProxyMainMenuUI(),
		intercepterUI: NewIntercepterUI(),
		replayUI:      replayUI{},
	}

	globalMainWin = mainWin

	//app.SetIcon(icon)

	mainWin.win = mainWin.app.NewWindow("Vbro Proxy 助安社区珍藏版")
	mainWin.win.Resize(fyne.NewSize(800, 600))
	mainWin.win.CenterOnScreen()

	mainWin.menu.setMainMenu()

	globalMainTabs = container.NewDocTabs(container.NewTabItem(mainWin.intercepterUI.name(),
		mainWin.intercepterUI.view()))
	mainWin.win.SetContent(globalMainTabs)

	mainWin.win.ShowAndRun()

	return mainWin
}
