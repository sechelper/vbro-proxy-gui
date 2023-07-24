package main

import (
	_ "embed"
	"github.com/flopp/go-findfont"
	"github.com/goki/freetype/truetype"
	"os"
	"vbro-gui/proxy"
	"vbro-gui/ui"
)

var (
	//go:embed assets/img/Icon.png
	icon []byte
)

func main() {
	go proxy.Run()
	ui.NewMainWin()
}

func init() {

	fontPath, err := findfont.Find("simhei.ttf")
	if err != nil {
		panic(err)
	}
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}
	if _, err = truetype.Parse(fontData); err != nil {
		panic(err)
	}
	if err = os.Setenv("FYNE_FONT", fontPath); err != nil {
		panic(err)
	}
}
