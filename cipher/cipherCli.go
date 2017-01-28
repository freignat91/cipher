package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type cipherCLI struct {
	verbose bool
	debug   bool
}

var (
	RootCmd = &cobra.Command{
		Use:   `cipher [OPTIONS] COMMAND [arg...]`,
		Short: "cipher",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}
)

func cli() {
	RootCmd.PersistentFlags().BoolVarP(&cipherCli.verbose, "verbose", "v", false, `Verbose output`)
	RootCmd.PersistentFlags().BoolVar(&cipherCli.debug, "debug", false, `Silence output`)
	cobra.OnInitialize(func() {
	})

	// versionCmd represents the agrid version
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version number of cipher",
		Long:  `Display the version number of cipher`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("cipher version: %s, build: %s)\n", Version, Build)
		},
	}
	RootCmd.AddCommand(versionCmd)

	//Execute commad
	cmd, _, err := RootCmd.Find(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error during: %s: %v\n", cmd.Name(), err)
		os.Exit(1)
	}

	os.Exit(0)
}
