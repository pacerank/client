// +build windows

package system

import (
	"syscall"
)

type target struct {
	processes []*Process
}

// kernel32.dll API calls
var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	procModule32First            = modKernel32.NewProc("Module32FirstW")
)

// User32.dll API calls
var (
	modUser32                 = syscall.NewLazyDLL("User32.dll")
	procForegroundWindow      = modUser32.NewProc("GetForegroundWindow")
	procWindowThreadProcessId = modUser32.NewProc("GetWindowThreadProcessId")
	procSetWindowsHookEx      = modUser32.NewProc("SetWindowsHookExW")
	procCallNextHookEx        = modUser32.NewProc("CallNextHookEx")
	procGetMessage            = modUser32.NewProc("GetMessageW")
	procUnhookWindowsHookEx   = modUser32.NewProc("UnhookWindowsHookEx")
)

// Type structure for windows API
type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	LONG      int32
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
	BYTE      int64
)
