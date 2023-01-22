/*
Copyright © 2023 PWSK info@pwsk.uk
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var input string
var output string
var time int
var domain string
var shared bool
var arch bool
var name string

//var sgn bool

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new dropper",
	Long:  `Creates a new dropper, options are local, remote or child. Each with their own flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please choose a dropper type like so: 'godropit new <local|remote|child> <flags>'")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	newCmd.PersistentFlags().StringVarP(&input, "in", "i", "", "input file, shellcode.bin or exe.")
	newCmd.MarkFlagRequired("in")
	newCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "droppername, don't include file extension.")
	newCmd.MarkFlagRequired("name")
	newCmd.PersistentFlags().StringVarP(&output, "out", "o", "", "output directory for your generated files.")
	newCmd.MarkFlagRequired("out")
	newCmd.PersistentFlags().IntVarP(&time, "time", "t", 10, "delay in seconds before decryption and execution of shellcode.")
	newCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "")
	newCmd.PersistentFlags().BoolVarP(&shared, "shared", "s", false, "Export dropper as DLL. Default is false")
	newCmd.PersistentFlags().BoolVar(&arch, "x86", false, "Attempts to generate an x86 dropper, completely untested. Use at own risk.")
	//newCmd.PersistentFlags().BoolVar(&sgn, "SGN", false, "Uses nextgen shikata ga nai to encode the shellcode.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
