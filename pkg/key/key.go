package key

import (
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"log"
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

func (k *Key) SetActiveEncryption() {
	appConfig, err := config.NewConfigFromFile("")
	if err != nil {
		log.Fatalf("Could not get application config: %v", err)
	}

	appConfig.EncryptionKeyName = k.Name
	err = appConfig.Write("")
	if err != nil {
		log.Fatalf("Could not write application config: %v", err)
	}
}

func (k *Key) SetActiveDecryption() {
	appConfig, err := config.NewConfigFromFile("")
	if err != nil {
		log.Fatalf("Could not get application config: %v", err)
	}

	appConfig.DecryptionKeyName = k.Name
	err = appConfig.Write("")
	if err != nil {
		log.Fatalf("Could not write application config: %v", err)
	}
}
