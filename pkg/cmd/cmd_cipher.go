package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type CipherEncryptCmd struct {
	TargetFile string
	SaveTo     string
	Password   string
}

type CipherDecryptCmd struct {
	Password   string
	TargetFile string
	SaveTo     string
}

func (c *CipherEncryptCmd) Run(cmd *cobra.Command, args []string) error {
	target, err := os.ReadFile(c.TargetFile)
	if err != nil {
		return fmt.Errorf("read file %s: %w", c.TargetFile, err)
	}
	encrypted, err := encrypt(target, c.Password)

	if c.SaveTo == "" {
		c.SaveTo = filepath.Base(c.TargetFile) + ".tzzencrypted"
	}

	if err := os.WriteFile(c.SaveTo, []byte(encrypted), os.ModePerm); err != nil {
		return fmt.Errorf("write to %s error: %w", c.SaveTo, err)
	}

	fmt.Printf("Encrypted to file: %s\n", c.SaveTo)

	return nil
}
func (c *CipherDecryptCmd) Run(cmd *cobra.Command, args []string) error {
	target, err := os.ReadFile(c.TargetFile)
	if err != nil {
		return fmt.Errorf("read file %s error: %w", c.TargetFile, err)
	}
	decrypted, err := decrypt(target, c.Password)
	if err != nil {
		return fmt.Errorf("decrypt error: %w", err)
	}

	if c.SaveTo == "" {
		targetFileName := filepath.Base(c.TargetFile)
		if strings.HasSuffix(targetFileName, ".tzzencrypted") {
			c.SaveTo = strings.TrimSuffix(targetFileName, ".tzzencrypted")
		} else {
			c.SaveTo = targetFileName + ".tzzdecrypted"
		}
	}

	if err := os.WriteFile(c.SaveTo, decrypted, os.ModePerm); err != nil {
		return fmt.Errorf("write to file %s error: %w", c.SaveTo, err)
	}

	return nil
}

func encrypt(target []byte, password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is empty")
	}
	// Generate key from password
	key := generateKey(password)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, target, nil)

	// Return base64 encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encrypted []byte, password string) ([]byte, error) {
	// Decode the base64 string
	ciphertext, err := base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		return nil, err
	}

	// Generate key from password
	key := generateKey(password)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Get nonce size
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Split nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func generateKey(password string) []byte {
	// Create a 32-byte key from password using SHA-256
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}
