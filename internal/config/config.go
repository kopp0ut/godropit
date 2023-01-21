package config

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func GoGetEnv(outdir string) {
	cmd := exec.Command("go", "env")

	// open the out file for writing
	outfile, err := os.Create(outdir + "goenv.txt")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	cmd.Wait()
}

func ReadEnv(envar string, outdir string) string {
	envFile := filepath.Join(outdir, "goenv.txt")

	data, err := ioutil.ReadFile(envFile)
	if err != nil {
		color.Red("Could not read goenv.txt, compilation may fail.")
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, envar) {
			envarout := strings.Split(line, "=")
			envarout[1] = strings.Trim(envarout[1], "\"")

			return envarout[1]

		}

	}

	color.Red("Couldn't Find go env %s, compilation may fail.", envar)

	return "failed"

}
