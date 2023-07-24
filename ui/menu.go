package ui

import (
	"fyne.io/fyne/v2"
)

type ProxyMainMenu struct {
}

func NewProxyMainMenuUI() *ProxyMainMenu {
	var proxyMainMenu = &ProxyMainMenu{}

	return proxyMainMenu
}

func (proxyMenu *ProxyMainMenu) setMainMenu() {

	var menuItems []*fyne.MenuItem

	// 拦截器
	menuItems = append(menuItems, fyne.NewMenuItem(IntercepterMenuText, globalMainWin.intercepterUI.action))

	// 重放
	menuItems = append(menuItems, fyne.NewMenuItem(ReplayMenuText, proxyMenu.todo))

	sy := fyne.NewMenu(ToolsMenuText, menuItems...)
	globalMainWin.win.SetMainMenu(fyne.NewMainMenu(sy))
}

func (proxyMenu *ProxyMainMenu) openTab(c *fyne.Container) {

}

func (proxyMenu *ProxyMainMenu) todo() {
	showAbout()
	//utils.ShowErrorDialog("功能暂未实现，敬请关注！", globalMainWin.win)
}
