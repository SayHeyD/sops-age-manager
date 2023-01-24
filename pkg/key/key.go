package key

import (
	"log"
	"os"
	"strings"
)

type Key struct {
	Name       string
	FileName   string
	PublicKey  string
	PrivateKey string
}

func NewKey(name string, path string, fileContents string) *Key {
	return &Key{
		Name:       name,
		FileName:   path,
		PublicKey:  getPubKeyFromFileContents(fileContents),
		PrivateKey: getPrivateKeyFromFileContents(fileContents),
	}
}

func getPrivateKeyFromFileContents(contents string) string {
	privateKeyPrefix := "AGE-SECRET-KEY"

	_, privateKeyWithoutPrefix, _ := strings.Cut(contents, privateKeyPrefix)

	privateKeyWithPrefix := privateKeyPrefix + privateKeyWithoutPrefix

	return strings.Trim(privateKeyWithPrefix, "\n\t ")
}

func getPubKeyFromFileContents(contents string) string {
	publicKeyPrefix := "public key: "
	privateKeyPrefix := "AGE-SECRET-KEY"

	_, afterPublicKeyString, _ := strings.Cut(contents, publicKeyPrefix)
	publicKey, _, _ := strings.Cut(afterPublicKeyString, privateKeyPrefix)

	publicKey = publicKey[:len(publicKey)-1]

	return publicKey
}

func (k *Key) SetActive() {
	err := os.Setenv("SOPS_AGE_KEY", k.PrivateKey)
	if err != nil {
		log.Fatalf("Cannot set env SOPS_AGE_KEY: %v", err)
	}
}
