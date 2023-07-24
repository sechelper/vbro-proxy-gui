package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type about struct {
	aboutDialog dialog.Dialog
}

func newAbout() *about {
	var a about
	//
	//logo := canvas.NewImageFromResource(theme.ResourceAppIcon)
	////logo.FillMode = canvas.ImageFillContain
	//logo.SetMinSize(fyne.NewSize(200, 200))

	content := widget.NewCard("", "",
		container.NewBorder(nil, nil, nil, nil, widget.NewRichTextFromMarkdown(`
## 代理抓包工具

[vbro-proxy](https://github.com/sechelper/vbro) 由[助安社区](http://secself.com)开源，打造中国最好用的攻防武器库。

【[提交BUG](https://github.com/sechelper/vbro/issues)】

`)))
	a.aboutDialog = dialog.NewCustom("助安社区", "关闭", content, globalMainWin.win)
	return &a
}

func showAbout() {
	newAbout().aboutDialog.Show()
}
