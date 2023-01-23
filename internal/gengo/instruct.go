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
	getSalmod := exec.Command(goBinPath, "get", "github.com/salukikit/go-util/pkg/box")
	getSalmod.Env = Env
	getSalmod.Stdout = &out
	getSalmod.Stderr = &stderr
	err = getSalmod.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy:\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	//setup to use dll
	if dll {
		/*
			//write dll file
			dllHeader := "dllmain.c"
			dllHeaderFile, err := os.Create(dllHeader)
			if err != nil {
				log.Fatalf("Error creating dropper file: %v ", err)
			}

			dllHeaderFile.WriteString(delivery.DllProcAttach)

			dllHeaderFile.Close()
		*/

		command = []string{"build", "-o", outname + ".dll", "-trimpath", "-buildmode=c-shared", `-ldflags="-w -s -H=windowsgui"`, fname}
	} else {
		command = []string{"build", "-o", outname + ".exe", "-trimpath", `-ldflags="-w -s -H=windowsgui"`, fname}
	}

	buildcmd := fmt.Sprintf("GOOS=windows GOARCH=%s go ", arch) + strings.Join(command, " ")

	color.Green("Prep done! go.mod and go.sum may need updating before compilation.")
	color.Cyan("This version of GoDropIt does not perform compilation to ensure that the end-user verifies the code.\n")
	color.Cyan("To compile navigate to %s and run:\n\n", outdir)
	fmt.Printf("%s\n\n", buildcmd)

	color.Cyan("Or if you're really cool use garble: https://github.com/burrowers/garble")

	return nil

}
