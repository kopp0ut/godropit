/*
Copyright © 2023 Phil Kopp
*/
package cmd

import (
	"log"

	"github.com/Epictetus24/godropit/internal/gengo"
	"github.com/spf13/cobra"
)

var loop bool

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Execute in the local process",
	Long:  `Executes in the local process, if running as an exe I'd recommend using -loop to hold the process open indefinitely.`,
	Run: func(cmd *cobra.Command, args []string) {
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}
		gengo.NewLocalDropper(input, output, domain, name, time, false, shared, arch, loop)
	},
}

func init() {
	newCmd.AddCommand(localCmd)

	localCmd.MarkFlagRequired("in")
	localCmd.MarkFlagRequired("name")
	localCmd.MarkFlagRequired("out")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	localCmd.Flags().BoolVarP(&loop, "loop", "l", false, "Add a for loop to keep the process alive. Warning: Will potentially keep process alive even if shellcode fails.")
}
