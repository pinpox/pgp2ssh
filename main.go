package main

import (
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/eddsa"
	"github.com/ProtonMail/go-crypto/openpgp/packet"

	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"reflect"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func readEntity(keypath string) (*openpgp.Entity, error) {
	f, err := os.Open(keypath)
	if err != nil {
		log.Println("Error opening file")
		return nil, err
	}
	defer f.Close()
	block, err := armor.Decode(f)
	if err != nil {
		log.Println("decoding")
		return nil, err
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

var (
	UnsupportedKeyType = errors.New("only ed25519 and rsa keys are supported")
)

func getEDDSAKey(castkey *eddsa.PrivateKey) []byte {
	log.Println("public key type:", reflect.TypeOf(castkey.PublicKey))
	var pubkey ed25519.PublicKey = castkey.PublicKey.X

	sshPub, err := ssh.NewPublicKey(pubkey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("public SSH key:\n" + string(ssh.MarshalAuthorizedKey(sshPub)))

	var privkey = ed25519.NewKeyFromSeed(castkey.D)

	privPem, err := ssh.MarshalPrivateKey(&privkey, "")
	if err != nil {
		log.Fatal(err)
	}
	return pem.EncodeToMemory(privPem)
}

func getRSAKey(castkey *rsa.PrivateKey) []byte {

	log.Println("public key type:", reflect.TypeOf(castkey.PublicKey))
	var pubkey rsa.PublicKey = castkey.PublicKey

	sshPub, err := ssh.NewPublicKey(&pubkey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("public SSH key:\n" + string(ssh.MarshalAuthorizedKey(sshPub)))

	// var privkey = ed25519.NewKeyFromSeed(castkey.D)

	privPem, err := ssh.MarshalPrivateKey(castkey, "")
	if err != nil {
		log.Fatal(err)
	}
	return pem.EncodeToMemory(privPem)
}

func main() {
	var keyfile string
	log.Println("Enter path to private PGP key (default: ./priv.asc):")
	_, err := fmt.Scanf("%s", &keyfile)
	if err != nil && err.Error() == "unexpected newline" {
		keyfile = "./priv.asc"
	} else if err != nil {
		log.Fatal(err)
	}

	e, err := readEntity(keyfile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Keys:")
	log.Println("[0]", e.PrimaryKey.KeyIdString()+" (primary)")
	for i := 0; i < len(e.Subkeys); i++ {
		log.Println(fmt.Sprintf("[%d]", i+1), e.Subkeys[i].PrivateKey.KeyIdString()+" (subkey)")
	}

	log.Println("Choose key by index (default: 0):")

	var keyIndex int
	_, err = fmt.Scanf("%d", &keyIndex)
	if err != nil && err.Error() == "unexpected newline" {
		keyIndex = 0
	} else if err != nil {
		log.Fatal(err)
	}

	var targetKey *packet.PrivateKey
	if keyIndex == 0 {
		log.Println(fmt.Sprintf("Continuing with key [%d]", keyIndex), e.PrimaryKey.KeyIdString())
		targetKey = e.PrivateKey
	} else if keyIndex > 0 {
		var subkey = e.Subkeys[keyIndex-1]
		log.Println(fmt.Sprintf("Continuing with key [%d]", keyIndex), subkey.PrivateKey.KeyIdString())
		targetKey = subkey.PrivateKey
	} else {
		log.Fatal("Invalid key index")
	}

	if targetKey.Encrypted {
		log.Println("Please enter passphrase to decrypt PGP key:")
		bytePassphrase, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal(err)
		}
		targetKey.Decrypt(bytePassphrase)
	}
	log.Println("private key type:", reflect.TypeOf(targetKey.PrivateKey))
	castkey_eddsa, ok_eddsa := targetKey.PrivateKey.(*eddsa.PrivateKey)
	if ok_eddsa {
		privateKeyPem := getEDDSAKey(castkey_eddsa)
		log.Println("Private SSH key:\n" + string(privateKeyPem))
		return
	}
	castkey_rsa, ok_rsa := targetKey.PrivateKey.(*rsa.PrivateKey)
	if ok_rsa {
		privateKeyPem := getRSAKey(castkey_rsa)
		log.Println("Private SSH key:\n" + string(privateKeyPem))
		return
	}
}
