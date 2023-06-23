package gengo

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fatih/color"
)

const (
	goBinPath = "go"
)

var Env []string

var Arch string

func init() {
	GoGetEnv()
	Env = []string{
		fmt.Sprintf("CC=%s", ReadEnv("CC")),
		fmt.Sprintf("CGO_ENABLED=%s", "0"),
		fmt.Sprintf("GOCACHE=%s", ReadEnv("GOCACHE")),
		fmt.Sprintf("GOMODCACHE=%s", ReadEnv("GOMODCACHE")),
		fmt.Sprintf("GOPRIVATE=%s", ReadEnv("GOPRIVATE")),
		fmt.Sprintf("PATH=%s:%s", path.Join(ReadEnv("GOVERSION"), "bin"), os.Getenv("PATH")),
		fmt.Sprintf("GOPATH=%s", ReadEnv("GOPATH")),
	}
}

func buildInstruct(outdir, fname string, dll bool, x86 bool) error {
	var command []string
	var err error

	err = os.Chdir(outdir)
	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	fmt.Printf("Shellcode files are now in %s.\n", outdir)
	fmt.Println("Using go for executable. Ensure go is in your PATH.")
	goBinPath := "go"

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	fmt.Printf("Getting go env and saving to %s\n", outdir+"/goenv.txt")
	GoGetEnv()
	os.Remove("go.mod")
	os.Remove("go.sum")
	if err != nil {
		log.Printf("Tried to clean any previous go.mod files but failed:\n%v.\n", err)
	}

	//go mod init the directory

	Arch = "amd64"
	if x86 {
		Arch = "386"
	}

	Env = append(Env, fmt.Sprintf("GOARCH=%s", Arch))
	Env = append(Env, fmt.Sprintf("GOOS=%s", "windows"))

	var out bytes.Buffer
	var stderr bytes.Buffer

	outname := fmt.Sprintf("%s_%s", strings.ReplaceAll(fname, ".go", ""), Arch)

	initCmd := exec.Command(goBinPath, "mod", "init", strings.ReplaceAll(fname, ".go", ""))
	initCmd.Env = Env

	initCmd.Stdout = &out
	initCmd.Stderr = &stderr
	err = initCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod init, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	initCmd.Wait()
	//go mod tidy to get depencies
	tidyCmd := exec.Command(goBinPath, "mod", "tidy")
	tidyCmd.Env = Env
	tidyCmd.Stdout = &out
	tidyCmd.Stderr = &stderr
	err = tidyCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy:\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	tidyCmd.Wait()

	getWinmod := exec.Command(goBinPath, "get", "golang.org/x/sys/windows")
	getWinmod.Env = Env
	getWinmod.Stdout = &out
	getWinmod.Stderr = &stderr
	err = getWinmod.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy:\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	getSalmod := exec.Command(goBinPath, "get", "github.com/epictetus24/go-util/pkg/box")
	getSalmod.Env = Env
	getSalmod.Stdout = &out
	getSalmod.Stderr = &stderr
	err = getSalmod.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy:\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	if dll {
		Env = append(Env, fmt.Sprintf("CGO_ENABLED=%s", "1"))
		command = []string{"build", "-o", outname + ".dll", "-trimpath", "-buildmode=c-shared", `-ldflags="-w -s -H=windowsgui"`, fname}
	} else {
		Env = append(Env, fmt.Sprintf("CGO_ENABLED=%s", "0"))
		command = []string{"build", "-o", outname + ".exe", "-trimpath", `-ldflags="-w -s -H=windowsgui"`, fname}
	}

	buildcmd := fmt.Sprintf("GOOS=windows GOARCH=%s go ", Arch) + strings.Join(command, " ")

	color.Green("Prep done! go.mod and go.sum may need updating before compilation.")
	color.Cyan("This version of GoDropIt does not perform compilation to ensure that the end-user verifies the code.\n")
	color.Cyan("To compile navigate to %s and run:\n\n", outdir)
	fmt.Printf("%s\n\n", buildcmd)

	color.Cyan("Or if you're really cool use garble: https://github.com/burrowers/garble")

	return nil

}

// Compiles the resulting go file into the appropriate format.
func buildFileGo(outdir, fname string, dll bool, x86 bool) (bool, error) {
	var command []string
	var err error

	err = os.Chdir(outdir)
	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	outname := fmt.Sprintf("%s_%s", strings.ReplaceAll(fname, ".go", ""), Arch)

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

	color.Green("Dropper compiled with regular normal go compiler, find it at %s/%s.\n", outdir, outname)

	return true, nil

}

func buildWasm(outdir, fname string) error {
	Arch = "wasm"

	Env = append(Env, fmt.Sprintf("GOARCH=%s", Arch))
	Env = append(Env, fmt.Sprintf("GOOS=%s", "js"))

	var command []string
	var err error

	err = os.Chdir(outdir)
	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	fmt.Printf("Shellcode files are now in %s.\n", outdir)
	fmt.Println("Using go for executable. Ensure go is in your PATH.")
	goBinPath := "go"

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	fmt.Printf("Getting go env and saving to %s\n", outdir+"/goenv.txt")
	GoGetEnv()
	os.Remove("go.mod")
	os.Remove("go.sum")
	if err != nil {
		log.Printf("Tried to clean any previous go.mod files but failed:\n%v.\n", err)
	}

	//go mod init the directory

	var out bytes.Buffer
	var stderr bytes.Buffer

	outname := fmt.Sprintf("%s_%s", strings.ReplaceAll(fname, ".go", ""), Arch)

	initCmd := exec.Command(goBinPath, "mod", "init", strings.ReplaceAll(fname, ".go", ""))
	initCmd.Env = Env

	initCmd.Stdout = &out
	initCmd.Stderr = &stderr
	err = initCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod init, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	initCmd.Wait()
	//go mod tidy to get depencies
	tidyCmd := exec.Command(goBinPath, "mod", "tidy")
	tidyCmd.Env = Env
	tidyCmd.Stdout = &out
	tidyCmd.Stderr = &stderr
	err = tidyCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy:\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	tidyCmd.Wait()

	err = os.Chdir(outdir)
	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	outname = fmt.Sprintf("%s_%s", strings.ReplaceAll(fname, ".go", ""), Arch)

	outname = outname + ".wasm"
	command = []string{"build", "-o", outname, "-trimpath", `-ldflags=-w -s -H=windowsgui`, fname}

	cmd := exec.Command(goBinPath, command...)
	cmd.Env = Env
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with compiling, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	cmd.Wait()

	color.Green("Dropper compiled with regular normal go compiler, find it at %s/%s.\n", outdir, outname)

	return nil

}
