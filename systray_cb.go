//go:build darwin

package main

/*
*/
import "C"
import (
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var windowX, windowY float64

//export setWindowPosition
func setWindowPosition(x, y float64) {
	windowX = x
	windowY = y
}

//export goOnToggleWindow
func goOnToggleWindow() {
	if smInstance == nil || smInstance.app.ctx == nil {
		return
	}
	ctx := smInstance.app.ctx

	if smInstance.app.windowVisible {
		runtime.WindowHide(ctx)
		smInstance.app.windowVisible = false
		return
	}

	runtime.WindowSetPosition(ctx, int(windowX), int(windowY))
	runtime.WindowShow(ctx)
	smInstance.app.windowVisible = true
}

//export goOnQuit
func goOnQuit() {
	if smInstance != nil {
		smInstance.app.power.Deactivate()
	}
	os.Exit(0)
}
