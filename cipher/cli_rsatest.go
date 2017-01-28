package main

import (
	"fmt"
	"github.com/freignat91/cipher/rsa"
	"github.com/spf13/cobra"
)

var RSATestCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  `test`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cipherCli.test(cmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(RSATestCmd)
}

func (m *cipherCLI) test(cmd *cobra.Command, args []string) error {
	rsa.TestEncriptDecript()
	return nil
}
