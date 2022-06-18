package main

import (
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"
)

const MODULE_BASE_ADDRESS uint32 = 0x400000 // Where Client.exe starts
const MEMDATABASE_OFFSET uint32 = 0xd4
const LOWEST_BRANCH uint32 = 0x6c
const HIGHEST_BRANCH uint32 = 0x348
const BRANCH_SIZE uint32 = 999 // Each "branch" contains exactly 1,000 items (0-999)
const BRANCH_ITEM_NAME_OFFSET = 0xc

/* In order to verify correctness, we check for the normal, default
 * "Attack" skill, which has a known ID and will be one of a
 * handful of strings (client language support).
 */
const SEARCH_ITEM_ID uint32 = 540000

var SEARCH_ITEM_ALLOWED_NAMES = []string{"Attack", "Angreifen", "Ataque", "Attaque", "Atak"}

/* BranchItemInfo - a  struct from Runes of Magic's client
 * The memory DB in RoM contains an array of these structs.
 * Each (valid) entry would contain a branch ID (array index)
 * and the address in the game client to that branch.
 */
type BranchItemInfo struct {
	_        uint32
	branchId uint32
	_        [14]byte
	address  uint32
	_        uint32
	_        uint32
}

func Scan(handle windows.Handle, procId uint32, chunkStart uint32, chunkEnd uint32) (uint32, error) {
	count := (chunkEnd-chunkStart)/4 + 4

	branchListPtrs, err := getBranchListPtrs(handle, MODULE_BASE_ADDRESS+chunkStart, count)
	if err != nil {
		return 0, err
	}

	for index, branchListPtr := range branchListPtrs {
		if branchListPtr == 0 {
			continue
		}

		var branchListAddress uint32
		err := ReadProcessSimpleValue(handle, branchListPtr+MEMDATABASE_OFFSET, &branchListAddress)

		if err != nil || branchListAddress == 0 || branchListAddress == 1 {
			continue
		}

		for branch := LOWEST_BRANCH; branch < HIGHEST_BRANCH; branch += 4 {
			found := scanBranch(handle, branchListAddress, branch)

			if found { // Found a match!
				return chunkStart + uint32(index*4), nil
			}
		}
	}

	return 0, nil
}

func scanBranch(handle windows.Handle, branchListAddress uint32, branch uint32) bool {
	var branchAddress uint32
	err := ReadProcessSimpleValue(handle, branchListAddress+branch, &branchAddress)
	if err != nil || branchAddress == 0 || branchAddress == 0xFFFFFFFF {
		return false
	}

	var branchItemInfos [BRANCH_SIZE + 1]BranchItemInfo
	var size uint32 = uint32(unsafe.Sizeof(BranchItemInfo{})) * BRANCH_SIZE
	err = ReadProcessMemory(handle, branchAddress, uintptr(unsafe.Pointer(&branchItemInfos)), size)

	if err != nil {
		return false
	}

	for index, branchItemInfo := range branchItemInfos {
		// Read branch ID (index) doesn't match expectations or address is invalid, so skip attempting this one
		if branchItemInfo.branchId != uint32(index) || branchItemInfo.address == 0 {
			continue
		}

		// Verify that the item ID we read from memory is the same as our expectation
		var itemId uint32
		err = ReadProcessSimpleValue(handle, branchItemInfo.address, &itemId)
		if err != nil || itemId != SEARCH_ITEM_ID {
			continue
		}

		// Verify that the name we read from memory matches our expectation
		var nameAddress uint32
		err = ReadProcessSimpleValue(handle, branchItemInfo.address+BRANCH_ITEM_NAME_OFFSET, &nameAddress)
		if err != nil || nameAddress == 0 {
			continue
		}

		var name [32]byte
		ReadProcessMemory(handle, nameAddress, uintptr(unsafe.Pointer(&name)), 32)
		sName := windows.ByteSliceToString(name[:])

		if isValidItemName(sName) {
			return true
		}
	}

	return false
}

func isValidItemName(name string) bool {
	for _, v := range SEARCH_ITEM_ALLOWED_NAMES {
		if v == name {
			return true
		}
	}

	return false
}

func getBranchListPtrs(handle windows.Handle, address uint32, count uint32) ([]uint32, error) {
	buffer := make([]uint32, count)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buffer))

	ReadProcessMemory(
		handle,
		address,
		uintptr(unsafe.Pointer(hdr.Data)),
		4*count,
	)

	return buffer, nil
}
