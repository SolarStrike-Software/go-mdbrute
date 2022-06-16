package main

import (
	"fmt"
	"log"

	"github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"
)

const CHUNK_SIZE uint32 = 0x1000
const START_ADDRESS uint32 = 0x00620000
const END_ADDRESS uint32 = 0x00650000

func main() {
	err := winio.EnableProcessPrivileges([]string{"SeDebugPrivilege"})
	if err != nil {
		log.Fatalf("could not elevate process. Reason: %s", err)
	}

	procId, err := GetProcId("Client.exe")

	fmt.Printf("Found proc ID: 0x%X\n", procId)

	if err != nil {
		log.Fatal("Encountered error while trying to find process: %w", err)
	}

	if procId == 0 {
		log.Fatal("Could not find target process")
	}

	handle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, procId)
	defer windows.CloseHandle(handle)

	if err != nil {
		log.Fatalf("could not open process. Reason: %s", err)
	}

	halfAddrRange := (END_ADDRESS - START_ADDRESS) / 2
	for i := START_ADDRESS + halfAddrRange; i <= END_ADDRESS; i += CHUNK_SIZE {
		addr, err := Work(handle, procId, i, i+CHUNK_SIZE)

		if err != nil {
			log.Fatal(err)
		} else if addr == 0 {
			addr, err = Work(handle, procId, i-CHUNK_SIZE, i)
		}

		if addr != 0 && err == nil {
			fmt.Printf("Found address 0x%X\n", addr)
			return
		}
	}
}

// func startWorker(wg *sync.WaitGroup, handle windows.Handle, procId uint32, from uint32, to uint32) {
// 	defer wg.Done()
// 	fmt.Printf("Starting worker\n")
// 	addr, err := Work(handle, procId, from, to)

// 	if err == nil {
// 		fmt.Printf("Found address: 0x%X\n", addr)
// 	} else {
// 		fmt.Println(err)
// 	}
// }
