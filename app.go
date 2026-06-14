//go:build darwin

package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx           context.Context
	power         *PowerManager
	onExpire      func()
	windowVisible bool
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
		if a.onExpire != nil {
			a.onExpire()
		}
	})

}

func (a *App) Toggle() (bool, error) {
	if a.power.GetStatus() {
		err := a.power.Deactivate()
		if err != nil {
			return false, err
		}
		return false, nil
	}
	err := a.power.Activate()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *App) ActivateForMinutes(minutes int) (bool, error) {
	err := a.power.ActivateWithDuration(minutes)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *App) Deactivate() (bool, error) {
	err := a.power.Deactivate()
	if err != nil {
		return false, err
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
