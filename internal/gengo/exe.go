package gengo

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Binject/go-donut/donut"
	"github.com/Epictetus24/godropit/pkg/dropfmt"
	"github.com/fatih/color"
)

// Uses donut to generate 64bit shellcode from an exe.
func DonutShellcode(srcFile string, x86 bool) ([]byte, error) {

	donutArch := donut.X64
	if x86 {
		donutArch = donut.X32
	}

	config := new(donut.DonutConfig)
	config.Arch = donutArch
	config.Entropy = uint32(2)
	config.OEP = uint64(0)

	config.InstType = donut.DONUT_INSTANCE_PIC

	config.Entropy = uint32(3)
	config.Bypass = 3
	config.Compress = uint32(1)
	config.Format = uint32(1)
	config.Verbose = true

	config.ExitOpt = uint32(1)
	payload, err := donut.ShellcodeFromFile(srcFile, config)
	if err == nil {
		return nil, err

	}

	return payload.Bytes(), nil
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
