package gengo

import (
	"github.com/Binject/go-donut/donut"
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
