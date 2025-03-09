package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func cipherCmd() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use: "cipher",
	}

	encrypt := cmdpkg.CipherEncryptCmd{}
	encryptCmd := &cobra.Command{
		Use:  "encrypt",
		RunE: encrypt.Run,
	}
	encryptCmd.Flags().StringVar(
		&encrypt.TargetFile,
		"target",
		"",
		"target file to be encrypted",
	)
	encryptCmd.Flags().StringVar(
		&encrypt.SaveTo,
		"saveto",
		"",
		"save the encrypted content to this file",
	)
	encryptCmd.Flags().StringVar(
		&encrypt.Password,
		"password",
		"",
		"password",
	)

	decrypt := cmdpkg.CipherDecryptCmd{}
	decryptCmd := &cobra.Command{
		Use:  "decrypt",
		RunE: decrypt.Run,
	}
	decryptCmd.Flags().StringVar(
		&decrypt.TargetFile,
		"target",
		"",
		"target file to be decrypted",
	)
	decryptCmd.Flags().StringVar(
		&decrypt.SaveTo,
		"saveto",
		"",
		"save the decrypted content to this file",
	)
	decryptCmd.Flags().StringVar(
		&decrypt.Password,
		"password",
		"",
		"password",
	)

	cobraCmd.AddCommand(
		encryptCmd,
		decryptCmd,
	)

	return cobraCmd
}
