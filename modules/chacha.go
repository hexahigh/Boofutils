package modules

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"context"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20"
)

func Chacha_main(password string, decrypt bool, file string, outFile string, mute bool, keyfile string, verbose bool) {
	v := verbose
	// Start the music and console logging
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if !mute {
		go PlayAudioMult(ctx, "valiant.wav,move.wav,legacy.wav,honor.wav,carry_on.wav,back_to_work.wav")
	}
	if file == "" {
		fmt.Println("No file provided")
		os.Exit(1)
	}

	if keyfile != "" {
		verbosePrintln("Using keyfile: "+keyfile, v)
		key, err := os.ReadFile(keyfile)
		if err != nil {
			fmt.Println("Error reading keyfile:", err)
			os.Exit(1)
		}
		password = base64.StdEncoding.EncodeToString(key)
	}

	if password == "" {
		fmt.Println("No password provided, using default.")
		fmt.Println(ColorRedHighIntensity24bit, "THIS SHOULD ONLY BE USED FOR TESTING", ColorReset)
		password = "cb62kZQ6si3fwvTEAvpJUg5KywN6YBurJKr8C7at5y6BtshnoqYSva3wktNfXzkfDDNH4zZGmdJ9w55bVLeYBdWZVParZHXks2otJ4rUdG2VU4rn6CcuCSdwRKhvFRzj"
	}

	if len(password) < 16 {
		fmt.Println(ColorRedBoldHighIntensity24bit, "Password is less than 16 characters, consider using a stronger password", ColorReset)
	}

	if outFile == "" && !decrypt {
		outFile = file + ".chachacha"
	}
	if outFile == "" && decrypt {
		outFile = strings.TrimSuffix(file, ".chachacha")
	}
	verbosePrintln("Output file: "+outFile, v)

	// Call the appropriate function
	if decrypt {
		verbosePrintln("Decrypting file: "+file, v)
		if err := decryptFile(file, password, outFile); err != nil {
			panic(err)
		}
	} else {
		verbosePrintln("Encrypting file: "+file, v)
		if err := encryptFile(file, password, outFile); err != nil {
			panic(err)
		}
	}
	fmt.Println("Done! Press enter to exit.")
	fmt.Scanln()
	cancel()
}

func encryptFile(filePath string, password string, outFile string) error {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	// Derive a key from the password using Argon2
	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Create a new ChaCha20 cipher
	cipher, err := chacha20.NewUnauthenticatedCipher(key, make([]byte, chacha20.NonceSize))
	if err != nil {
		return err
	}

	// Encrypt the data
	cipher.XORKeyStream(data, data)

	// Write the salt and the encrypted data back to the file
	return os.WriteFile(outFile, append(salt, data...), 0644)
}

func decryptFile(filePath string, password string, outFile string) error {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Split the salt and the encrypted data
	salt, data := data[:16], data[16:]

	// Derive a key from the password using Argon2
	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Create a new ChaCha20 cipher
	cipher, err := chacha20.NewUnauthenticatedCipher(key, make([]byte, chacha20.NonceSize))
	if err != nil {
		return err
	}

	// Decrypt the data
	cipher.XORKeyStream(data, data)

	// Write the decrypted data back to the file
	return os.WriteFile(outFile, data, 0644)
}

// TODO: Add Quad encryption
func encryptFileQuad(filePath string, password string, outFile string) error {
	return nil
}

func decryptFileQuad(filePath string, password string, outFile string) error {
	return nil
}
