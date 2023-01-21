package gengo

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/Epictetus24/godropit/internal/config"
	"github.com/Epictetus24/godropit/pkg/dropfmt"
	"github.com/fatih/color"
)

type Generator struct {
	GenName   string
	Dll       bool
	Shellcode dropfmt.DropFmt
	FileExt   string
}

func initConfig(outdir string) {

	genDir := filepath.Join(outdir, ".godropit")

	os.MkdirAll(genDir, 0666)

	if _, err := os.Stat(genDir + "/goenv.txt"); errors.Is(err, os.ErrNotExist) {
		config.GoGetEnv(genDir)
	}
}

// Compiles the resulting go file into the appropriate format.
func buildFileGo(outdir, fname string, dll bool, x86 bool) (bool, error) {
	var command []string
	var err error
	goBinPath := "/usr/local/go/bin/go"

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}
	err = os.Chdir(outdir)
	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	//go mod init the directory
	initCmd := exec.Command(goBinPath, "mod", "init", "betaDrop")
	initCmd.Stdout = &out
	initCmd.Stderr = &stderr
	err = initCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod init, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return false, err
	}
	//go mod tidy to get depencies
	tidyCmd := exec.Command(goBinPath, "mod", "tidy")
	tidyCmd.Stdout = &out
	tidyCmd.Stderr = &stderr
	err = tidyCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return false, err
	}

	if dll {
		command = []string{"build", "-o", outdir + "dropper.dll", "-trimpath", "-buildmode=c-shared", `-ldflags='-w -s -H=windowsgui'`, fname}
	} else {
		command = []string{"build", "-o", outdir + "dropper.exe", "-trimpath", `-ldflags='-w -s -H=windowsgui'`, fname}
	}

	cmd := exec.Command(goBinPath, command...)
	arch := "amd64"
	if x86 {
		arch = "386"
	}

	cmd.Env = []string{
		fmt.Sprintf("CC=%s", config.ReadEnv("CC", outdir)),
		fmt.Sprintf("CGO_ENABLED=%s", config.ReadEnv("CGO_ENABLED", outdir)),
		fmt.Sprintf("GOOS=%s", "windows"),
		fmt.Sprintf("GOARCH=%s", arch),
		fmt.Sprintf("GOCACHE=%s", config.ReadEnv("GOCACHE", outdir)),
		fmt.Sprintf("GOMODCACHE=%s", config.ReadEnv("GOMODCACHE", outdir)),
		fmt.Sprintf("GOPRIVATE=%s", ""),
		fmt.Sprintf("PATH=%s:%s", path.Join(config.ReadEnv("GOVERSION", outdir), "bin"), os.Getenv("PATH")),
		fmt.Sprintf("GOPATH=%s", config.ReadEnv("GOPATH", outdir)),
	}

	if dll {
		cmd.Env[0] = fmt.Sprintf("CC=%s", "x86_64-w64-mingw32-gcc")
		cmd.Env[1] = fmt.Sprintf("CGO_ENABLED=%s", "1")

	}

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with compiling, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return false, err
	}

	return true, nil

}
