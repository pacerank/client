// +build windows

package system

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"syscall"
	"unsafe"
)

const (
	MaxModuleName32 = 255
	MaxPath         = 260
)

// processEntry is the Windows API structure that contains a process's information.
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/ns-tlhelp32-processentry32w
type processEntry struct {
	Size              DWORD
	CntUsage          DWORD
	ProcessID         DWORD
	DefaultHeapID     uintptr
	ModuleID          DWORD
	CntThreads        DWORD
	ParentProcessID   DWORD
	PriorityClassBase LONG
	Flags             DWORD
	ExeFile           [MaxPath]uint16
}

// moduleEntry is the Windows API structure that contains data of modules that belongs to a specific process.
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/ns-tlhelp32-moduleentry32w
type moduleEntry struct {
	Size         DWORD
	ModuleID     DWORD
	ProcessID    DWORD
	GlblcntUsage DWORD
	ProccntUsage DWORD
	ModBaseAddr  BYTE
	ModBaseSize  DWORD
	Module       HANDLE
	ModuleName   [MaxModuleName32 + 1]uint16
	ExePath      [MaxPath]uint16
}

func (t *target) Processes() ([]*Process, error) {
	var result = make([]*Process, 0)

	// Create a process snap handler, with TH32CS_SNAPPROCESS (0x00000002)
	hProcessSnap, _, _ := procCreateToolhelp32Snapshot.Call(0x00000002, 0)
	if hProcessSnap < 0 {
		return result, syscall.GetLastError()
	}

	// Close handler after method is ready
	defer func() {
		_, _, _ = procCloseHandle.Call(hProcessSnap)
	}()

	var process processEntry
	process.Size = DWORD(unsafe.Sizeof(process))

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
			if p.Parent == int64(process.ParentProcessID) && p.FileName == fileName {
				result[index].Children = append(result[index].Children, int64(process.ProcessID))
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

		path, err := getProcessPath(int64(process.ProcessID))
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

		result = append(result, &Process{
			ProcessID:  int64(process.ProcessID),
			Parent:     int64(process.ParentProcessID),
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

func (t *target) ActiveProcess() (*Process, error) {
	handle, _, _ := procForegroundWindow.Call()
	if handle == 0 {
		return nil, errors.New("no window is currently active")
	}

	var processID DWORD

	_, _, _ = procWindowThreadProcessId.Call(handle, uintptr(unsafe.Pointer(&processID)))
	if processID == 0 {
		return nil, errors.New("process id could not be found for window handle")
	}

	process, err := getProcess(processID)
	if err != nil {
		return nil, err
	}

	return process, nil
}

// Get process information
func getProcess(processID DWORD) (result *Process, err error) {
	// Create a process snap handler, with TH32CS_SNAPTHREAD (0x00000004)
	hProcessSnap, _, _ := procCreateToolhelp32Snapshot.Call(0x00000002, uintptr(processID))
	if hProcessSnap < 0 {
		return nil, syscall.GetLastError()
	}

	// Close handler after method is ready
	defer func() {
		_, _, _ = procCloseHandle.Call(hProcessSnap)
	}()

	var process processEntry
	process.Size = DWORD(unsafe.Sizeof(process))

	// Get the first process in the list
	if ok, _, _ := procProcess32First.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
		return nil, errors.New("could not retrieve process info")
	}

	var children []int64

	for {
		if process.ParentProcessID == processID {
			children = append(children, int64(process.ProcessID))
			if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
				break
			}

			continue
		}

		if process.ProcessID == processID {
			end := 0
			for {
				if process.ExeFile[end] == 0 {
					break
				}

				end++
			}

			fileName := syscall.UTF16ToString(process.ExeFile[:end])

			path, err := getProcessPath(int64(process.ProcessID))
			if err != nil {
				log.Error().Err(err).Msg("could not get process path")
				if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
					break
				}

				continue
			}

			cs, err := checksum(path)
			if err != nil {
				log.Error().Err(err).Msg("could not get checksum of executable")
				if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
					break
				}

				continue
			}

			result = &Process{
				ProcessID:  int64(process.ProcessID),
				Parent:     int64(process.ParentProcessID),
				Children:   nil,
				FileName:   fileName,
				Checksum:   cs,
				Executable: path,
			}
		}

		if ok, _, _ := procProcess32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&process))); ok == 0 {
			break
		}
	}

	if result == nil {
		return nil, errors.New(fmt.Sprintf("could not find process with process id %d", processID))
	}

	result.Children = children

	return result, nil
}

// Get executable path for a ProcessID in string format
func getProcessPath(processID int64) (string, error) {
	// Create a module snap handler with TH32CS_SNAPMODULE (0x00000008) for given process ID
	hModuleSnap, _, _ := procCreateToolhelp32Snapshot.Call(0x00000008, uintptr(processID))

	// Close handler after method is ready
	defer func() {
		_, _, _ = procCloseHandle.Call(hModuleSnap)
	}()

	var module moduleEntry
	module.Size = DWORD(unsafe.Sizeof(module))

	// Get the first module, as it is the link to the executable. Other modules(DLLs) are irrelevant
	if ok, _, err := procModule32First.Call(hModuleSnap, uintptr(unsafe.Pointer(&module))); ok == 0 {
		return "", err
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
