package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Microsoft/go-winio"
	"github.com/jessevdk/go-flags"
	"golang.org/x/sys/windows"
)

const CHUNK_SIZE uint32 = 0x1000
const START_ADDRESS uint32 = 0x00620000
const END_ADDRESS uint32 = 0x00650000

const GAME_EXE_NAME = "Client.exe"
const GAME_PARENT_NAME = "RoMLauncher.exe"

var opts struct {
	ProcId       uint32 `long:"proc" description:"The Process ID to target"`
	StartAddress uint32 `long:"start" description:"The memory address to start the scan from"`
	EndAddress   uint32 `long:"end" description:"The memory address to end the scan on"`
}

func main() {
	err := winio.EnableProcessPrivileges([]string{"SeDebugPrivilege"})
	if err != nil {
		log.Fatalf("could not elevate process. Reason: %s", err)
	}

	_, err = flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	var procId uint32
	if opts.ProcId != 0 {
		procId = opts.ProcId
	} else {
		procId, err = GetProcId(GAME_EXE_NAME, GAME_PARENT_NAME)
		if err != nil {
			log.Fatal("Encountered error while trying to find process: %w", err)
		}

		if procId == 0 {
			log.Fatal("Could not find target process")
		}
	}

	fmt.Printf("Using process ID: %d\n", procId)
	handle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, procId)
	defer windows.CloseHandle(handle)

	if err != nil {
		log.Fatalf("could not open process. Reason: %s", err)
	}

	addr, err := outwardScan(handle, procId, START_ADDRESS, END_ADDRESS, CHUNK_SIZE)
	if err != nil {
		log.Fatal(err)
	}

	if addr != 0 {
		fmt.Printf("Found address 0x%X\n", addr)
		return
	}

	fmt.Printf("Could not find any valid address.")
}

/* Start the scan in the middle of the range and work outwards.
 * The address is most likely to be near the middle so this
 * can speed things up significantly.
 */
func outwardScan(handle windows.Handle, procId uint32, startAddress uint32, endAddress uint32, chunkSize uint32) (uint32, error) {
	halfAddrRange := (endAddress - startAddress) / 2
	for i := startAddress + halfAddrRange; i <= endAddress; i += chunkSize {
		fmt.Printf("Scanning range 0x%X -> 0x%X\n", i, i+chunkSize)
		addr, err := Scan(handle, procId, i, i+chunkSize)
		if err != nil {
			return 0, err
		}

		if addr != 0 {
			return addr, nil
		}

		fmt.Printf("Scanning range 0x%X -> 0x%X\n", i-chunkSize, i)
		addr, err = Scan(handle, procId, i-chunkSize, i)
		if addr != 0 && err == nil {

			return addr, nil
		}
	}

	return 0, nil
}
