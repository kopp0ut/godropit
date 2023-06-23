package local

import (
	"godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateFiber", "CreateThread", "CreateThreadNative", "EtwpCreateETWThread", "NtQueueAPCThreadExLocal", "goSyscall", "UUIDFromStringA"}

var LeetDroppers = []string{"[HellsGate] BananaPhone", "[Callback] EnumerateChildWindows", "[Callback] EnumerateLoadedModules", "[Callback] CreateThreadPoolWait"}

var Hold = `
	for {

	}
`

func SelectLocal(leet bool) (Dlls, Inject, Import, Extra string) {
	if leet {
		Droppers = append(Droppers, LeetDroppers...)
	}
	_, selected, _ := dropfmt.PromptList(Droppers, "Select the local dropper you would like to use:")

	switch selected {
	case "CreateFiber":
		Dlls = CreateFiberDlls
		Inject = CreateFiber
		Import = CreateFiberImports
	case "CreateThread":
		Dlls = CreateThreadDlls
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
		Dlls = BananaPhoneDlls
		Inject = BananaPhone
		Import = BananaPhoneImports
		Extra = BananaPhoneExtra

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
		Extra = "//notreq"

	}

	return Dlls, Inject, Import, Extra
}
