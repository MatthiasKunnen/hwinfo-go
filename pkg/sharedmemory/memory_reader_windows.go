package sharedmemory

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	hwinfoMutexName            = "Global\\HWiNFO_SM2_MUTEX"
	hwinfoSensorsMapFilename   = "Global\\HWiNFO_SENS_SM2"
	hwinfoSensorStringLength   = 128
	hwinfoUnitStringLength     = 16
	hwinfoSharedMemoryMaxBytes = 20_000_000
)

// HwinfoSensorStringAscii is a fixed length byte array of 8-bit ASCII encoded characters.
// The specific extended ASCII codepage used depends on the system's locale.
//
// Get the codepage used by your system using this powershell command:
//
//	[System.Text.Encoding]::Default
//
// The string it contains is padded by nul bytes.
type HwinfoSensorStringAscii = [hwinfoSensorStringLength]byte

// HwinfoSensorStringUtf8 is a fixed length byte array of UTF-8 encoded characters.
// The string it contains is padded by nul bytes.
// To convert it to a string, use
// [github.com/MatthiasKunnen/hwinfo-go/pkg/util.NulTerminatedUtf8ByteArrayToString].
// It is used in labels for sensor and reading.
type HwinfoSensorStringUtf8 = [hwinfoSensorStringLength]byte

// HwinfoUnitStringAscii is the same as [HwinfoSensorStringAscii] but used for unit strings such as
// °C and MHz.
type HwinfoUnitStringAscii = [hwinfoUnitStringLength]byte

// HwinfoUnitStringUtf8 is the same as  [HwinfoSensorStringUtf8] but used for unit strings such as
// °C and MHz.
type HwinfoUnitStringUtf8 = [hwinfoUnitStringLength]byte

// MemoryReader allows for reading the shared memory provided by HWiNFO.
// Create an instance using NewMemoryReader.
// Use the [MemoryReader.Open] function to make MemoryReader ready to start reading.
//
// # Locking
//
// Locking is required to prevent HWiNFO from changing the shared memory while it is being read.
// Locking is performed by using the Lock function.
// While the lock is held, HWiNFO will pause any updates to the shared memory and its own UI.
// Therefore, locks should be released as soon as possible.
//
// By default, the read functions will enforce locking.
// If necessary, it is possible to disable this enforcement by setting
// [MemoryReader.DisableLockEnforcement] to `false`.
// Do note that this causes the risk of receiving garbage data when HWiNFO changes the shared memory
// layout.
//
// # Quick processing
//
// If the data can be processed quickly, use the following:
//
//  1. [MemoryReader.Open]
//  2. [MemoryReader.Lock]
//  3. [Reader.GetHeader]
//  4. [Reader.GetSensors] / [Reader.GetReadings]
//  5. Process the data quickly, this will block HWiNFO
//  6. [MemoryReader.ReleaseLock]
//
// # Slow processing
//
// If the data cannot be processed quickly, use the following:
//
//  1. [MemoryReader.Open]
//  2. [MemoryReader.Lock]
//  3. [Reader.GetHeader]
//  4. [MemoryReader.Copy]
//  5. [MemoryReader.ReleaseLock]
//  6. Process the data, this will not block HWiNFO
type MemoryReader struct {
	DisableLockEnforcement bool
	mutex                  windows.Handle
	mmfHandle              windows.Handle
	mmfPtr                 uintptr
	Data                   Reader
}

func NewMemoryReader() *MemoryReader {
	memoryReader := &MemoryReader{}
	memoryReader.Data.GetPointer = func() (uintptr, error) {
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

// Close deallocates the resources used by the MemoryReader.
// Call [MemoryReader.Open] before performing any further operations.
func (reader *MemoryReader) Close() error {
	return errors.Join(
		reader.ReleaseLock(),
		reader.closeMapView(),
		reader.closeMappedFileHandle(),
	)
}

// Open readies the MemoryReader for reading the shared memory.
// Use [MemoryReader.Close] when there will be no more reads.
func (reader *MemoryReader) Open() error {
	mmf, err := openFileMapping(windows.FILE_MAP_READ, 0, hwinfoSensorsMapFilename)

	if err != nil {
		return fmt.Errorf("error opening file mapping: %w", err)
	}

	reader.mmfHandle = mmf

	// https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-mapviewoffile
	mmfPtr, err := windows.MapViewOfFile(
		mmf,
		windows.FILE_MAP_READ,
		0,
		0,
		0,
	)
	if err != nil {
		defer reader.closeMappedFileHandle()
		return fmt.Errorf("error mapping view of file: %w", err)
	}
	reader.mmfPtr = mmfPtr

	return nil
}

// Copy copies the shared memory in order to perform processing after the lock has been released.
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
// Release it using [MemoryReader.ReleaseLock].
func (reader *MemoryReader) Lock() error {
	mutex, err := openMutex(hwinfoMutexName)
	if err != nil {
		return fmt.Errorf("error opening HWiNFO mutex (%s): %w", hwinfoMutexName, err)
	}
	reader.mutex = mutex

	_, err = windows.WaitForSingleObject(mutex, 200)
	if err != nil {
		return fmt.Errorf("error waiting for HWiNFO mutex (%s): %w", hwinfoMutexName, err)
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
		return fmt.Errorf("error releasing HWiNFO lock: %w", err)
	}

	return nil
}

func (reader *MemoryReader) closeMappedFileHandle() error {
	if reader.mmfHandle != 0 {
		if err := windows.CloseHandle(reader.mmfHandle); err != nil {
			return fmt.Errorf("error closing handle to memory mapped file: %w", err)
		}
		reader.mmfHandle = 0
	}

	return nil
}

func (reader *MemoryReader) closeMapView() error {
	if reader.mmfPtr != 0 {
		if err := windows.UnmapViewOfFile(reader.mmfPtr); err != nil {
			return fmt.Errorf("error unmapping MemoryReader view of file: %w", err)
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
