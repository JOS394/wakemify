//go:build darwin

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx           context.Context
	power         *PowerManager
	onExpire      func()
	windowVisible bool
	lastBlurTime  time.Time
}

func NewApp() *App {
	return &App{
		power: NewPowerManager(),
	}
}

func (a *App) SetOnExpire(cb func()) {
	a.onExpire = cb
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	runtime.WindowHide(ctx)

	a.power.SetOnExpire(func() {
		runtime.EventsEmit(a.ctx, "timer-expired")
		if smInstance != nil {
			smInstance.UpdateUI(false)
		}
		if a.onExpire != nil {
			a.onExpire()
		}
	})

}

func (a *App) HandleWindowBlur() {
	a.windowVisible = false
	a.lastBlurTime = time.Now()
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "window-visible", false)
	}
}

func (a *App) Toggle() (bool, error) {
	if a.power.GetStatus() {
		err := a.power.Deactivate()
		if err != nil {
			return false, err
		}
		if smInstance != nil {
			smInstance.UpdateUI(false)
		}
		return false, nil
	}
	err := a.power.Activate()
	if err != nil {
		return false, err
	}
	if smInstance != nil {
		smInstance.UpdateUI(true)
	}
	return true, nil
}

func (a *App) ActivateForMinutes(minutes int) (bool, error) {
	err := a.power.ActivateWithDuration(minutes)
	if err != nil {
		return false, err
	}
	if smInstance != nil {
		smInstance.UpdateUI(true)
	}
	return true, nil
}

func (a *App) Deactivate() (bool, error) {
	err := a.power.Deactivate()
	if err != nil {
		return false, err
	}
	if smInstance != nil {
		smInstance.UpdateUI(false)
	}
	return false, nil
}

func (a *App) GetStatus() bool {
	return a.power.GetStatus()
}

func (a *App) GetRemainingTime() int {
	d := a.power.GetRemainingTime()
	return int(d.Seconds())
}

func (a *App) QuitApp() {
	a.power.Deactivate()
	runtime.Quit(a.ctx)
}

func (a *App) RestoreMenuPosition() {
	runtime.WindowSetPosition(a.ctx, int(windowX), int(windowY))
}

func (a *App) SetLaunchAtStartup(enabled bool) bool {
	home, _ := os.UserHomeDir()
	label := "com.wakemify.launcher"
	plistPath := filepath.Join(home, "Library", "LaunchAgents", label+".plist")

	if enabled {
		exePath := os.Args[0]
		// Try to resolve symlinks to get the real path
		if resolved, err := filepath.EvalSymlinks(exePath); err == nil {
			exePath = resolved
		}
		plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <false/>
</dict>
</plist>`, label, exePath)
		if err := os.WriteFile(plistPath, []byte(plist), 0644); err != nil {
			return false
		}
		exec.Command("launchctl", "load", plistPath).Run()
		return true
	}

	exec.Command("launchctl", "unload", plistPath).Run()
	os.Remove(plistPath)
	return false
}

func (a *App) IsLaunchAtStartup() bool {
	home, _ := os.UserHomeDir()
	plistPath := filepath.Join(home, "Library", "LaunchAgents", "com.wakemify.launcher.plist")
	_, err := os.Stat(plistPath)
	return err == nil
}

func (a *App) OpenURL(url string) error {
	return exec.Command("open", url).Start()
}
