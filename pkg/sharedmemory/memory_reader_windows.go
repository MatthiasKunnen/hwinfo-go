package sharedmemory

import (
	"errors"
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	hwinfoSensorsMapFilename   = "Global\\HWiNFO_SENS_SM2"
	hwinfoSensorStringLength   = 128
	hwinfoUnitStringLength     = 16
	hwinfoSharedMemoryMaxBytes = 20_000_000
)

// MemoryReader allows for reading the shared memory provided by HWiNFO.
// Create an instance using NewMemoryReader.
// Use the Open function to make MemoryReader ready to start reading.
//
// # Locking
//
// Locking is required to prevent HWiNFO from changing the shared memory while it is being read.
// Locking is performed by using the Lock function.
// While the lock is held, HWiNFO will pause any updates to the shared memory and its own UI.
// Therefore, locks should be released as soon as possible.
//
// By default, the read functions will enforce locking.
// If necessary, it is possible to disable this enforcement by setting DisableLockEnforcement to
// `false`.
// Do note that this causes the risk of receiving garbage data when HWiNFO changes the shared memory
// layout.
//
// # Quick processing
//
// If the data can be processed quickly, use the following:
//
//  1. Open
//  2. Lock
//  3. [Reader.GetHeader]
//  4. [Reader.GetSensors] / [Reader.GetReadings]
//  5. Process the data quickly, this will block HWiNFO
//  6. ReleaseLock
//
// # Slow processing
//
// If the data cannot be processed quickly, use the following:
//
//  1. Open
//  2. Lock
//  3. [Reader.GetHeader]
//  4. Copy
//  5. ReleaseLock
//  6. Process the data, this will not block HWiNFO
type MemoryReader struct {
	DisableLockEnforcement bool
	mutex                  windows.Handle
	mmfHandle              windows.Handle
	mmfPtr                 uintptr
	Reader                 Reader
}

func NewMemoryReader() *MemoryReader {
	memoryReader := &MemoryReader{}
	memoryReader.Reader.GetPointer = func() (uintptr, error) {
		if memoryReader.mmfPtr == 0 {
			return 0, errors.New("shared memory not open, use Open first")
		}

		if err := memoryReader.enforceLock(); err != nil {
			return 0, err
		}

		return memoryReader.mmfPtr, nil
	}
	return memoryReader
}

func (reader *MemoryReader) Close() error {
	if err := reader.closeMapView(); err != nil {
		return err
	}
	if err := reader.closeHandle(); err != nil {
		return err
	}

	return nil
}

func (reader *MemoryReader) Open() error {
	mmf, err := OpenFileMapping(windows.FILE_MAP_READ, 0, hwinfoSensorsMapFilename)

	if err != nil {
		return err
	}

	reader.mmfHandle = mmf

	mmfPtr, err := windows.MapViewOfFile(
		mmf,
		windows.FILE_MAP_READ,
		0,
		0,
		0,
	)
	if err != nil {
		defer reader.closeHandle()
		return err
	}
	reader.mmfPtr = mmfPtr

	return nil
}

// Copy copies the shared memory in order to perform processing after the lock has been released.
// After
func (reader *MemoryReader) Copy(info *HwinfoHeader) *BytesReader {
	offset := reader.mmfPtr
	size := info.ReadingSectionOffset + info.ReadingAmount*info.ReadingSize
	byteSlice := make([]byte, size)

	// See https://stackoverflow.com/a/39215454/2512498
	copy(byteSlice, (*(*[hwinfoSharedMemoryMaxBytes]byte)(unsafe.Pointer(offset)))[:size])

	// Per byte copy, seems 50% slower
	//for i := uint32(0); i < size; i++ {
	//	byteSlice[i] = *(*byte)(unsafe.Pointer(offset + uintptr(i)))
	//}

	return NewBytesReader(byteSlice)
}

// Lock acquires the HWiNFO mutex.
// While holding the mutex, HWiNFO will pause so keep this lock as short as possible.
// Locks should be held until data is processed or copied.
// Release it using ReleaseLock.
func (reader *MemoryReader) Lock() error {
	mutex, err := openMutex("Global\\HWiNFO_SM2_MUTEX")
	if err != nil {
		return err
	}
	reader.mutex = mutex

	_, err = windows.WaitForSingleObject(mutex, 200)
	if err != nil {
		return err
	}

	return nil
}

// ReleaseLock releases the HWiNFO mutex.
// Returns nil when: the mutex is successfully released or the mutex was not held.
//
// After releasing the lock, any results from GetSensors and GetReadings should no longer be used
// as HWiNFO could have changed the memory layout which can turn the results into garbage.
// E.g. HWiNFO could have since detected a new sensor which might cause part of the data to be
// offset.
//
// Errors returned are documented at
// https://learn.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-releasemutex.
func (reader *MemoryReader) ReleaseLock() error {
	if reader.mutex == 0 {
		return nil
	}

	err := windows.ReleaseMutex(reader.mutex)
	reader.mutex = 0
	if err != nil {
		return err
	}

	return nil
}

func (reader *MemoryReader) closeHandle() error {
	if reader.mmfHandle != 0 {
		if err := windows.CloseHandle(reader.mmfHandle); err != nil {
			return err
		}
		reader.mmfHandle = 0
	}

	return nil
}

func (reader *MemoryReader) closeMapView() error {
	if reader.mmfPtr != 0 {
		if err := windows.UnmapViewOfFile(reader.mmfPtr); err != nil {
			return err
		}
		reader.mmfPtr = 0
	}

	return nil
}

func (reader *MemoryReader) enforceLock() error {
	if reader.DisableLockEnforcement {
		return nil
	}

	if reader.mutex != 0 {
		return nil
	}

	return errors.New("lock not acquired. Acquire it using Lock()")
}
