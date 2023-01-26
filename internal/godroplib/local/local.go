package local

import (
	"github.com/Epictetus24/godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateFiber", "CreateThread", "CreateThreadNative", "EtwpCreateETWThread", "NtQueueAPCThreadExLocal", "goSyscall", "UUIDFromStringA", "[HeavensGate] BananaPhone", "[Callback] EnumerateChildWindows", "[Callback] EnumerateLoadedModules", "[Callback] CreateThreadPoolWait"}

var Hold = `
	for {

	}
`

func SelectLocal() (Dlls, Inject, Import, Extra string) {
	_, selected, _ := dropfmt.PromptList(Droppers, "Select the local dropper you would like to use:")

	switch selected {
	case "CreateFiber":
		Dlls = CreateFiberDlls
		Inject = CreateFiber
		Import = CreateFiberImports
	case "CreateThread":
		Dlls = ""
		Inject = CreateThread
		Import = CreateThreadImports
	case "CreateThreadNative":
		Dlls = CreateThreadNativeDlls
		Inject = CreateThreadNative
		Import = CreateThreadNativeImports
	case "EtwpCreateETWThread":
		Dlls = EtwpCreateETWThreadDlls
		Inject = EtwpCreateETWThread
		Import = EtwpCreateETWThreadImports
	case "NtQueueAPCThreadExLocal":
		Dlls = NtQueueAPCThreadExLocalDlls
		Inject = NtQueueAPCThreadExLocal
		Import = NtQueueAPCThreadExLocalImports
		Extra = NtQueueAPCThreadExLocalExtra
	case "goSyscall":
		Dlls = goSyscallDlls
		Inject = goSyscall
		Import = goSyscallImports
	case "UUIDFromStringA":
		Dlls = UUIDFromStringADlls
		Inject = UUIDFromStringA
		Import = UUIDFromStringAImports
		Extra = UUIDFromStringAExtra
	case "[HellsGate] BananaPhone":
		Dlls = ""
		Inject = bananaPhone
		Import = bananaPhoneImports
		Extra = bananaPhoneExtra

	case "[Callback] EnumerateChildWindows":
		Dlls = EnumChildWindowsDlls
		Inject = EnumChildWindows
		Import = EnumChildWindowsImports
		Extra = ""

	case "[Callback] EnumerateLoadedModules":
		Dlls = EnumerateLoadedModulesDlls
		Inject = EnumerateLoadedModules
		Import = EnumerateLoadedModulesImports
		Extra = ""

	case "[Callback] CreateThreadPoolWait":
		Dlls = CreateThreadPoolWaitDlls
		Inject = CreateThreadPoolWait
		Import = CreateThreadPoolWaitImports
		Extra = ""

	}

	return Dlls, Inject, Import, Extra
}
