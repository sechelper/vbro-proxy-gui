package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/sechelper/vbro/utils"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type intercepterUI struct {
	valveBtn       *widget.Button // 拦截 / 放行
	forwardBtn     *widget.Button // 转发
	dropBtn        *widget.Button // 丢弃
	requestTex     *widget.Entry
	requestTexChan chan []byte
	intercept      bool
	intercepter    *websocket.Conn
}

func (ui *intercepterUI) icon() *fyne.Resource {
	return nil

}

func (ui *intercepterUI) name() string {
	return IntercepterMenuText
}

func NewIntercepterUI() *intercepterUI {
	return &intercepterUI{}
}

func (ui *intercepterUI) view() *fyne.Container {
	ui.requestTexChan = make(chan []byte)
	ui.requestTex = widget.NewMultiLineEntry()
	ui.requestTex.Wrapping = fyne.TextWrapWord

	ui.valveBtn = widget.NewButton(InterceptButtonText, ui.valveTapped)

	ui.forwardBtn = widget.NewButton(ForwardButtonText, ui.forwardTapped)

	ui.dropBtn = widget.NewButton(DropButtonText, ui.dropTapped)

	go func() {

		u := url.URL{
			Scheme: "ws",
			Host:   "localhost:8081",
			Path:   "/intercepter",
		}

		time.Sleep(2 * time.Second)
		dialer := websocket.Dialer{
			NetDial: func(network, addr string) (net.Conn, error) {
				return net.Dial("unix", "vbro-proxy.sock")
			},
		}
		// Authorization: Basic base64(vbro)
		_i, _, err := dialer.Dial(u.String(), http.Header{"Authorization": []string{"dmJybw=="}})
		if err != nil {
			log.Error().Msg(err.Error())
		}
		ui.intercepter = _i

		// 接收ID
		_, b, err := ui.intercepter.ReadMessage()
		if err != nil {
			panic(err)
		}
		fmt.Println("id: ", string(b))
		go func() {
			for {
				_, b, err = ui.intercepter.ReadMessage()
				if err != nil {
					log.Error().Msg("ui.intercepter.ReadMessage error" + err.Error())
					return
				}

				ui.requestTex.SetText(string(b))
			}

		}()
	}()

	return container.NewBorder(nil, nil, nil, container.NewVBox(
		ui.valveBtn,
		ui.forwardBtn,
		ui.dropBtn,
	), ui.requestTex)
}

func (ui *intercepterUI) action() {
	item := container.NewTabItem(ui.name(), ui.view())
	globalMainTabs.Append(item)
	globalMainTabs.Select(item)
}

func (ui *intercepterUI) valveTapped() {
	if ui.valveBtn.Text == PassButtonText {
		ui.valveBtn.SetText(InterceptButtonText)
		ui.intercept = false
	} else if ui.valveBtn.Text == InterceptButtonText {
		ui.valveBtn.SetText(PassButtonText)
		ui.intercept = true
	}
	var valve = "off"
	if ui.intercept {
		valve = "on"
	}
	formData := url.Values{}
	formData.Set("valve", valve)

	res, err := utils.Request("vbro-proxy.sock", "POST",
		"http://unix/intercepter/valve", strings.NewReader(formData.Encode()))
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	defer res.Body.Close()

}

func (ui *intercepterUI) forwardTapped() {
	// TODO 处理请求
	if len(ui.requestTex.Text) == 0 {
		return
	}
	if err := ui.intercepter.WriteMessage(websocket.TextMessage, []byte(ui.requestTex.Text)); err != nil {
		log.Error().Str("Error sending message:", err.Error())
		return
	}
	ui.requestTex.SetText("")
}

func (ui *intercepterUI) dropTapped() {
	if err := ui.intercepter.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
		log.Error().Str("Error sending message:", err.Error())
		return
	}
	ui.requestTex.SetText("")
}
