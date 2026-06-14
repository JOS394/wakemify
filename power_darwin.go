//go:build darwin

package main

/*
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include <stdlib.h>
#include <IOKit/pwr_mgt/IOPMLib.h>
#include <CoreFoundation/CoreFoundation.h>

IOPMAssertionID createPMAssertion(const char* reason) {
    CFStringRef reasonStr = CFStringCreateWithCString(kCFAllocatorDefault, reason, kCFStringEncodingUTF8);
    IOPMAssertionID assertionID;
    IOReturn result = IOPMAssertionCreateWithName(
        kIOPMAssertionTypePreventUserIdleDisplaySleep,
        kIOPMAssertionLevelOn,
        reasonStr,
        &assertionID
    );
    CFRelease(reasonStr);
    if (result == kIOReturnSuccess) return assertionID;
    return 0;
}

int releasePMAssertion(IOPMAssertionID assertionID) {
    if (assertionID == 0) return -1;
    IOReturn result = IOPMAssertionRelease(assertionID);
    if (result == kIOReturnSuccess) return 0;
    return -1;
}
*/
import "C"
import (
	"errors"
	"sync"
	"time"
	"unsafe"
)

var (
	errFailedToCreateAssertion  = errors.New("failed to create power assertion")
	errFailedToReleaseAssertion = errors.New("failed to release power assertion")
)

type PowerManager struct {
	mu          sync.Mutex
	active      bool
	assertionID C.IOPMAssertionID
	timer       *time.Timer
	deadline    time.Time
	onExpire    func()
}

func NewPowerManager() *PowerManager {
	return &PowerManager{}
}

func (pm *PowerManager) SetOnExpire(cb func()) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.onExpire = cb
}

func (pm *PowerManager) Activate() error {
	cReason := C.CString("Wakemify - Preventing sleep")
	defer C.free(unsafe.Pointer(cReason))

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.active {
		return nil
	}

	id := C.createPMAssertion(cReason)
	if id == 0 {
		return errFailedToCreateAssertion
	}

	if pm.timer != nil {
		pm.timer.Stop()
		pm.timer = nil
	}

	pm.assertionID = id
	pm.active = true
	pm.deadline = time.Time{}
	return nil
}

func (pm *PowerManager) Deactivate() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.active {
		return nil
	}

	if pm.timer != nil {
		pm.timer.Stop()
		pm.timer = nil
	}

	pm.deadline = time.Time{}
	result := C.releasePMAssertion(pm.assertionID)
	pm.assertionID = 0
	pm.active = false

	if result != 0 {
		return errFailedToReleaseAssertion
	}
	return nil
}

func (pm *PowerManager) ActivateWithDuration(minutes int) error {
	cReason := C.CString("Wakemify - Preventing sleep")
	defer C.free(unsafe.Pointer(cReason))

	pm.mu.Lock()

	if pm.active {
		if pm.timer != nil {
			pm.timer.Stop()
			pm.timer = nil
		}
		result := C.releasePMAssertion(pm.assertionID)
		pm.assertionID = 0
		pm.active = false
		if result != 0 {
			pm.mu.Unlock()
			return errFailedToReleaseAssertion
		}
	}

	id := C.createPMAssertion(cReason)
	if id == 0 {
		pm.mu.Unlock()
		return errFailedToCreateAssertion
	}

	pm.assertionID = id
	pm.active = true
	pm.deadline = time.Now().Add(time.Duration(minutes) * time.Minute)

	onExpire := pm.onExpire
	pm.timer = time.AfterFunc(time.Duration(minutes)*time.Minute, func() {
		pm.mu.Lock()
		if !pm.active {
			pm.mu.Unlock()
			return
		}
		if pm.timer != nil {
			pm.timer.Stop()
			pm.timer = nil
		}
		C.releasePMAssertion(pm.assertionID)
		pm.assertionID = 0
		pm.active = false
		pm.deadline = time.Time{}
		cb := onExpire
		pm.mu.Unlock()
		if cb != nil {
			cb()
		}
	})
	pm.mu.Unlock()

	return nil
}

func (pm *PowerManager) GetStatus() bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.active
}

func (pm *PowerManager) GetRemainingTime() time.Duration {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.active || pm.deadline.IsZero() {
		return 0
	}

	d := time.Until(pm.deadline)
	if d < 0 {
		return 0
	}
	return d
}
