package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func ShowErrorDialog(msg string, parent fyne.Window) {
	dialog.ShowInformation("错误信息", msg, parent)
}
