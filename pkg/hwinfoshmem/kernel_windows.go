package hwinfoshmem

import (
	"golang.org/x/sys/windows"
	"os"
	"unsafe"
)

var (
	modkernel32         = windows.NewLazyDLL("kernel32.dll")
	procOpenFileMapping = modkernel32.NewProc("OpenFileMappingW")
)

// openFileMapping implements the OpenFileMappingW Windows function and allows for opening a named file mapping object.
// This is the basis for working with shared memory.
// See https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-openfilemappingw.
func openFileMapping(access uint32, inheritHandle uint32, name string) (handle windows.Handle, err error) {
	namePointer, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	r1, _, err := procOpenFileMapping.Call(uintptr(access), uintptr(inheritHandle), uintptr(unsafe.Pointer(namePointer)))
	handle = windows.Handle(r1)

	if handle == 0 {
		if err == windows.ERROR_FILE_NOT_FOUND {
			err = &os.PathError{Path: name, Op: "OpenFileMapping", Err: err}
		} else {
			err = os.NewSyscallError("OpenFileMapping", err)
		}
	}

	if err == windows.Errno(0) {
		err = nil
	}

	return
}
