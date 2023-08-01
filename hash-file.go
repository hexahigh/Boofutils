package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Don't touch any of this code. It barely works.

func hf_main() {
	fmt.Println("What is the path of the file?")
	filePath := askInput()
	//md5
	fmt.Println("Calculating MD5...")
	md5hash, err := hash_file_md5(filePath)
	if err != nil {
		fmt.Println("Error calculating MD5:", err)
	}
	//sha1
	fmt.Println("Calculating SHA1...")
	sha1hash, err := hash_file_sha1(filePath)
	if err != nil {
		fmt.Println("Error calculating SHA1:", err)
	}
	//sha256
	fmt.Println("Calculating SHA256...")
	sha256hash, err := hash_file_sha256(filePath)
	if err != nil {
		fmt.Println("Error calculating SHA256:", err)
	}
	fmt.Println("Done! Here are the hashes:")
	fmt.Println("MD5:", md5hash)
	fmt.Println("SHA1:", sha1hash)
	fmt.Println("SHA256:", sha256hash)
}

func hash_file_md5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func hash_file_sha256(filePath string) (string, error) {
	var returnSHA256String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA256String, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA256String, err
	}
	hashInBytes := hash.Sum(nil)[:32] // SHA256 hash is 32 bytes
	returnSHA256String = hex.EncodeToString(hashInBytes)
	return returnSHA256String, nil
}

func hash_file_sha1(filePath string) (string, error) {
	var returnSHA1String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String, err
	}
	defer file.Close()
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String, err
	}
	hashInBytes := hash.Sum(nil)[:20] // SHA1 hash is 20 bytes
	returnSHA1String = hex.EncodeToString(hashInBytes)
	return returnSHA1String, nil
}
