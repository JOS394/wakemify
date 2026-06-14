//go:build darwin

package main

/*
*/
import "C"
import (
	"os"
	"time"

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
	go func() {
		if smInstance == nil || smInstance.app.ctx == nil {
			return
		}
		ctx := smInstance.app.ctx

		if smInstance.app.windowVisible {
			runtime.WindowHide(ctx)
			smInstance.app.windowVisible = false
			runtime.EventsEmit(ctx, "window-visible", false)
			return
		}

		if time.Since(smInstance.app.lastBlurTime) < 200*time.Millisecond {
			return
		}

		runtime.WindowSetPosition(ctx, int(windowX), int(windowY))
		runtime.WindowShow(ctx)
		smInstance.app.windowVisible = true
		runtime.EventsEmit(ctx, "window-visible", true)
	}()
}

//export goOnQuit
func goOnQuit() {
	if smInstance != nil {
		smInstance.app.power.Deactivate()
	}
	os.Exit(0)
}
