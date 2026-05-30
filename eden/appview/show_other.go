//go:build !windows

package appview

// showWindow is a no-op on non-Windows platforms; the SW_HIDE startup quirk
// that hides the webview under Neovim's jobstart() is Windows-specific.
func (a *AppView) showWindow() {}
