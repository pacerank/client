// +build windows

package system

import (
	"errors"
	"syscall"
	"unsafe"
)

type target struct{}

// Windows API functions
var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	procModule32First            = modKernel32.NewProc("Module32FirstW")
)

const (
	MaxModuleName32 = 255
	MaxPath         = 260
)

// processEntry is the Windows API structure that contains a process's information.
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/ns-tlhelp32-processentry32w
type processEntry struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MaxPath]uint16
}

// moduleEntry is the Windows API structure that contains data of modules that belongs to a specific process.
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/ns-tlhelp32-moduleentry32w
type moduleEntry struct {
	Size         uint32
	ModuleID     uint32
	ProcessID    uint32
	GlblcntUsage uint32
	ProccntUsage uint32
	ModBaseAddr  uint64
	ModBaseSize  uint32
	Module       syscall.Handle
	ModuleName   [MaxModuleName32 + 1]uint16
	ExePath      [MaxPath]uint16
}

func (t *target) Processes() ([]Process, error) {
	var result = make([]Process, 0)

	// Create a process snap handler, with TH32CS_SNAPPROCESS (0x00000002)
	hProcessSnap, _, _ := procCreateToolhelp32Snapshot.Call(0x00000002, 0)
	if hProcessSnap < 0 {
		return result, syscall.GetLastError()
	}

	defer procCloseHandle.Call(hProcessSnap)

	var process processEntry
	process.Size = uint32(unsafe.Sizeof(process))

	// Get the first process in the list
	if ok, _, _ := procProcess32First.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
		return nil, errors.New("could not retrieve process info")
	}

	for {
		var skip bool
		if process.ProcessID == 0 {
			skip = true
		}

		if process.ExeFile[0] == 0 {
			skip = true
		}

		if skip {
			if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
				break
			}

			continue
		}

		end := 0
		for {
			if process.ExeFile[end] == 0 {
				break
			}

			end++
		}

		fileName := syscall.UTF16ToString(process.ExeFile[:end])

		// Check if this process is a part of main process
		for index, p := range result {
			if p.Pid == int64(process.ParentProcessID) && p.FileName == fileName {
				result[index].ModulePid = append(result[index].ModulePid, int64(process.ProcessID))
				skip = true
				break
			}
		}

		if skip {
			if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
				break
			}

			continue
		}

		path, err := t.getProcessPath(process.ProcessID)
		if err != nil {
			if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
				break
			}

			continue
		}

		cs, err := checksum(path)
		if err != nil {
			if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
				break
			}

			continue
		}

		result = append(result, Process{
			Pid:        int64(process.ProcessID),
			Ppid:       int64(process.ParentProcessID),
			Checksum:   cs,
			Executable: path,
			FileName:   fileName,
		})

		if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
			break
		}
	}

	return result, nil
}

// Get executable path for a ProcessID in string format
func (t *target) getProcessPath(processID uint32) (string, error) {
	// Create a module snap handler with TH32CS_SNAPMODULE (0x00000008) for given process ID
	hModuleSnap, _, _ := procCreateToolhelp32Snapshot.Call(0x00000008, uintptr(processID))
	defer procCloseHandle.Call(hModuleSnap)

	var module moduleEntry
	module.Size = uint32(unsafe.Sizeof(module))

	// Get the first module, as it is the link to the executable. Other modules(DLLs) are irrelevant
	if ok, _, _ := procModule32First.Call(hModuleSnap, uintptr(unsafe.Pointer(&module))); ok == 0 {
		return "", errors.New("could not read module for process")
	}

	end := 0
	for {
		if module.ExePath[end] == 0 {
			break
		}

		end++
	}

	path := syscall.UTF16ToString(module.ExePath[:end])

	if path == "" {
		return "", errors.New("process does not have any executable path")
	}

	return path, nil
}
