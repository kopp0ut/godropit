/*
Copyright © 2023 Phil Kopp
*/
package cmd

import (
	"log"

	"godropit/internal/godroplib/local"
	"godropit/pkg/gengo"

	"github.com/spf13/cobra"
)

var loop bool

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Execute in the local process",
	Long:  `Executes in the local process, if running as an exe I'd recommend using -loop to hold the process open indefinitely.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkImg()
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}

		name = check(name)

		var localDrop gengo.Dropper
		localDrop.Dlls, localDrop.Inject, localDrop.Import, localDrop.Extra = local.SelectLocal(gengo.Leet)
		localDrop.Delay = time
		localDrop.Arch = arch
		localDrop.Shared = shared
		localDrop.Debug = debug

		if loop {
			localDrop.Hold = local.Hold
		} else {
			localDrop.Hold = "//notreq"
		}

		localDrop.Dtype = "//notreq"

		gengo.NewDropper(localDrop, name, domain, input, output, stagerurl, imgpath, hostname, ua, sgn)
	},
}

func init() {
	newCmd.AddCommand(localCmd)

	localCmd.MarkFlagRequired("in")
	localCmd.MarkFlagRequired("name")
	localCmd.MarkFlagRequired("out")

	// local specific flags
	localCmd.Flags().BoolVarP(&loop, "loop", "l", false, "Add a for loop to keep the process alive. Warning: Will potentially keep process alive even if shellcode fails.")
}
