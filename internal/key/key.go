package key

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Key struct {
	Name       string
	PublicKey  string
	PrivateKey string
}

func NewKey(name string, fileContents string) *Key {
	return &Key{
		Name:       name,
		PublicKey:  getPubKeyFromFileContents(fileContents),
		PrivateKey: getPrivKeyFromFileContents(fileContents),
	}
}

func getPrivKeyFromFileContents(contents string) string {
	privateKeyPrefix := "AGE-SECRET-KEY"

	_, privKeyWithoutPrefix, _ := strings.Cut(contents, privateKeyPrefix)

	return privateKeyPrefix + privKeyWithoutPrefix
}

func getPubKeyFromFileContents(contents string) string {
	pubKeyPrefix := "public key: "
	privateKeyPrefix := "AGE-SECRET-KEY"

	_, afterPubKeyString, _ := strings.Cut(contents, pubKeyPrefix)
	pubKey, _, _ := strings.Cut(afterPubKeyString, privateKeyPrefix)

	return pubKey
}

func (k *Key) SetActive() {
	err := os.Setenv("SOPS_AGE_KEY", k.PrivateKey)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot set env SOPS_AGE_KEY: %v", err))
	}
}
