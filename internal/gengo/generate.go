package gengo

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Epictetus24/godropit/pkg/dropfmt"

	"github.com/fatih/color"
)

type Generator struct {
	GenName   string
	Dll       bool
	Shellcode dropfmt.DropFmt
	FileExt   string
}

// Compiles the resulting go file into the appropriate format.
func buildFileGo(outdir, fname string, dll bool, x86 bool) (bool, error) {
	var command []string
	var err error

	err = os.Chdir(outdir)
	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	fmt.Printf("Shellcode files are now in %s.\n", outdir)
	fmt.Println("Using /usr/local/go/bin/go for go executable.")
	goBinPath := "/usr/local/go/bin/go"

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	arch := "amd64"
	if x86 {
		arch = "386"
	}

	Env := []string{
		fmt.Sprintf("CC=%s", "gcc"),
		fmt.Sprintf("CGO_ENABLED=%s", "0"),
		fmt.Sprintf("GOOS=%s", "windows"),
		fmt.Sprintf("GOARCH=%s", arch),
		fmt.Sprintf("GOCACHE=%s", ReadEnv("GOCACHE")),
		fmt.Sprintf("GOMODCACHE=%s", ReadEnv("GOMODCACHE")),
		fmt.Sprintf("GOPRIVATE=%s", ""),
		fmt.Sprintf("PATH=%s:%s", path.Join(ReadEnv("GOVERSION"), "bin"), os.Getenv("PATH")),
		fmt.Sprintf("GOPATH=%s", ReadEnv("GOPATH")),
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	outname := fmt.Sprintf("%s_%s", strings.ReplaceAll(fname, ".go", ""), arch)

	//setup to use dll
	if dll {

		outname = outname + ".dll"
		command = []string{"build", "-x", "-o", outname, "-trimpath", "-buildmode=c-shared", `-ldflags=-w -s -H=windowsgui`, fname}
	} else {
		outname = outname + ".exe"
		command = []string{"build", "-o", outname, "-trimpath", `-ldflags=-w -s -H=windowsgui`, fname}
	}

	cmd := exec.Command(goBinPath, command...)
	cmd.Env = Env
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if dll {
		cmd.Env[0] = fmt.Sprintf("CC=%s", "x86_64-w64-mingw32-gcc")
		cmd.Env[1] = fmt.Sprintf("CGO_ENABLED=%s", "1")
	}

	err = cmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with compiling, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return false, err
	}

	cmd.Wait()

	color.Green("Dropper compiled, find it at %s/%s.\n", outdir, outname)

	return true, nil

}
