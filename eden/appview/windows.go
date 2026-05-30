//go:build windows

package appview

import "syscall"

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procShowWindow       = user32.NewProc("ShowWindow")
	procSetForegroundWin = user32.NewProc("SetForegroundWindow")
)

const swShow = 5

// showWindow re-shows the window hidden under nvim jobstart; see #2.
func (a *AppView) showWindow() {
	hwnd := a.Window()
	if hwnd == nil {
		return
	}
	procShowWindow.Call(uintptr(hwnd), swShow)
	procSetForegroundWin.Call(uintptr(hwnd))
}
