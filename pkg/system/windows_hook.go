// +build windows

package system

import (
	"syscall"
	"unsafe"
)

// Callback function that is used for events
type HookProc func(int, WPARAM, LPARAM) LRESULT

// idHook values for setting up event listener
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw#parameters
const (
	whCallWndProc     = 4
	whCallWndProcRet  = 12
	whCbt             = 5
	whDebug           = 9
	whForegroundIdle  = 11
	whGetMessage      = 3
	whJournalPlayback = 1
	whJournalRecord   = 0
	whKeyboard        = 2
	whKeyboardLL      = 13
	whMouse           = 7
	whMouseLL         = 14
	whMsgFilter       = -1
	whShell           = 10
	whSysMsgFilter    = 6
)

const (
	wmKeyDown    = 256
	wmSysKeyDown = 260
	wmKeyUp      = 257
	wmSysKeyUp   = 261
	wmKeyFirst   = 256
	wmKeyLast    = 264
)

// Installs an application-defined hook procedure into a hook chain. You would install a hook procedure to monitor the
// system for certain types of events. These events are associated either with a specific thread or with all threads in
// the same desktop as the calling thread.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
func setWindowsHookEx(idHook int, lpfn HookProc, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)

	return HHOOK(ret)
}

// Passes the hook information to the next hook procedure in the current hook chain. A hook procedure can call this
// function either before or after processing the hook information.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-callnexthookex
func callNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)

	return LRESULT(ret)
}

// Removes a hook procedure installed in a hook chain by the SetWindowsHookEx function.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwindowshookex
func unhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

// Contains information about a low-level keyboard input event.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-kbdllhookstruct
type kbDllHookStruct struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// Retrieves a message from the calling thread's message queue. The function dispatches incoming sent messages until a
// posted message is available for retrieval.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessagew
func getMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

// This function setups a listener on the channel that send back the byte type of the key press
func (t *target) ListenKeyboard(channel chan byte) {
	var keyboardHook HHOOK

	keyboardHook = setWindowsHookEx(whKeyboardLL, func(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
		if nCode == 0 && wParam == wmKeyDown {
			structure := (*kbDllHookStruct)(unsafe.Pointer(lParam))
			channel <- byte(structure.VkCode)
		}
		return callNextHookEx(keyboardHook, nCode, wParam, lParam)
	}, 0, 0)

	var msg MSG
	for getMessage(&msg, 0, 0, 0) != 0 {
		// Need to call get message proc to receive keyboard press
	}

	unhookWindowsHookEx(keyboardHook)
	keyboardHook = 0
}
