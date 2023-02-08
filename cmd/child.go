package cmd

import (
	"log"

	"godropit/internal/gengo"

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
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}
		gengo.NewChildDropper(input, output, domain, name, proc, procArgs, time, false, shared, arch)
	},
}

func init() {
	newCmd.AddCommand(childCmd)

	// Here you will define your flags and configuration settings.

	childCmd.MarkFlagRequired("in")
	childCmd.MarkFlagRequired("name")
	childCmd.MarkFlagRequired("out")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	childCmd.Flags().StringVarP(&proc, "proc", "p", "c:\\\\windows\\\\system32\\\\werfault.exe", "Child Process to execute in.")
	childCmd.Flags().StringVar(&procArgs, "args", "args", "Arguments to pass child proc, embed in quotes or escape spaces pls.")
}
