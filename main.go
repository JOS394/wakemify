//go:build darwin

package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	systrayMgr := NewSystrayManager(app)

	app.SetOnExpire(func() {
		systrayMgr.UpdateUI(false)
	})

	go systrayMgr.Run()

	err := wails.Run(&options.App{
		Title:     "Wakemify",
		Width:     240,
		Height:    290,
		MinWidth:  240,
		MinHeight: 290,
		MaxWidth:  240,
		MaxHeight: 290,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 30, G: 30, B: 30, A: 255},
		OnStartup:        app.startup,
		OnDomReady: func(ctx context.Context) {
			removeWindowButtons()
		},
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WindowIsTranslucent:  false,
			WebviewIsTransparent: false,
			About: &mac.AboutInfo{
				Title:   "Wakemify",
				Message: "Keep your Mac awake",
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
