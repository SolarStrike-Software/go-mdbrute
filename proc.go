package main

import (
	"fmt"
	"unsafe"

	"github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows"
)

var (
	kernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procReadProcessMemory = kernel32.NewProc("ReadProcessMemory")
)

const PROCESS_VM_READ = 0x0010
const PROCESS_ALL_ACCESS = (0x000F0000 | 0x00100000 | 0xFFF)

type SimpleValueType interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int | ~uint | ~int32 | ~uint32 | ~int64 | ~uint64 | ~float32 | ~float64
}

func ReadProcessMemory(handle windows.Handle, address uint32, buffer uintptr, nSize uint32) error {
	var (
		length uint32
	)

	result, _, err := procReadProcessMemory.Call(
		uintptr(handle),
		uintptr(address),
		buffer,
		uintptr(nSize),
		uintptr(unsafe.Pointer(&length)),
	)

	if result == 0 {
		if err != nil {
			return err
		}

		return fmt.Errorf("error reading process memory. Err code %w", windows.GetLastError())
	}

	return nil
}

func ReadProcessSimpleValue[V SimpleValueType](handle windows.Handle, address uint32, output *V) error {
	return ReadProcessMemory(handle, address, uintptr(unsafe.Pointer(output)), uint32(unsafe.Sizeof(*output)))
}

func GetProcId(name string, parentName string) (uint32, error) {

	procs, err := ps.Processes()

	if err != nil {
		return 0, err
	}

	for _, proc := range procs {
		if proc.Executable() == name {
			parent, err := ps.FindProcess(proc.PPid())

			if err == nil {
				if parentName == "" || parent.Executable() == parentName {
					return uint32(proc.Pid()), nil
				}
			}
		}
	}

	return 0, nil
}
