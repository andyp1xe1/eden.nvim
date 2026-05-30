//go:build windows

package appview

import (
	"syscall"
)

// When eden is launched by Neovim's jobstart(), the child process is created
// with libuv's UV_PROCESS_WINDOWS_HIDE flag, i.e. STARTUPINFO.wShowWindow =
// SW_HIDE. Windows applies that value to the process's *first* ShowWindow call,
// overriding webview's own ShowWindow(SW_SHOW). The window is therefore created
// but never displayed. Calling ShowWindow again here (the override only affects
// the first call) makes it visible.
var (
	user32                = syscall.NewLazyDLL("user32.dll")
	procShowWindow        = user32.NewProc("ShowWindow")
	procSetForegroundWin  = user32.NewProc("SetForegroundWindow")
)

const swShow = 5 // SW_SHOW

func (a *AppView) showWindow() {
	hwnd := a.Window()
	if hwnd == nil {
		return
	}
	procShowWindow.Call(uintptr(hwnd), swShow)
	procSetForegroundWin.Call(uintptr(hwnd))
}
