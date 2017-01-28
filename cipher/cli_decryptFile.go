package main

import (
	"fmt"
	"github.com/freignat91/cipher/rsa"
	"github.com/spf13/cobra"
	"time"
)

var DecryptFileCmd = &cobra.Command{
	Use:   "decryptFile [sourcefilePath] [targetFilePath] [privateeKeyFilePath]",
	Short: "decrypt file",
	Long:  `decrypt file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cipherCli.decryptFile(cmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(DecryptFileCmd)
}

func (m *cipherCLI) decryptFile(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage cipher decryptFile [sourcefilePath] [targetFilePath] [privateKeyFilePath]")
	}
	t0 := time.Now()
	if err := rsa.DecryptFile(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Printf("done time=%ds\n", time.Now().Sub(t0).Nanoseconds()/1000000000)
	return nil
}
