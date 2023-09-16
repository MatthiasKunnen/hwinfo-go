package sharedmemory

import (
	"golang.org/x/sys/windows"
)

func openMutex(name string) (windows.Handle, error) {
	nameUtf16, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	mutexHandle, err := windows.OpenMutex(windows.SYNCHRONIZE, false, nameUtf16)
	if err != nil {
		return 0, err
	}

	return mutexHandle, nil
}
