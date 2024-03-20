package main

import (
	// "bytes"
	"crypto/sha512"
	"fmt"
	"log"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/eddsa"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	 "github.com/davecgh/go-spew/spew"
	"golang.org/x/crypto/ssh"

	"strings"

	"crypto/ed25519"
	"errors"
	"reflect"

	// "filippo.io/edwards25519"
	"github.com/Mic92/ssh-to-age/bech32"
	"golang.org/x/crypto/curve25519"
	// "golang.org/x/crypto/curve25519"
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
	UnsupportedKeyType = errors.New("only ed25519 keys are supported")
)

func ed25519PrivateKeyToCurve25519(pk ed25519.PrivateKey) ([]byte, error) {
	h := sha512.New()
	_, err := h.Write(pk.Seed())
	if err != nil {
		return []byte{}, err
	}
	out := h.Sum(nil)
	return out[:curve25519.ScalarSize], nil
}

func SSHPrivateKeyToAge(sshKey, passphrase []byte) (*string, error) {
	var (
		privateKey interface{}
		err        error
	)
	if len(passphrase) > 0 {
		privateKey, err = ssh.ParseRawPrivateKeyWithPassphrase(sshKey, passphrase)
	} else {
		privateKey, err = ssh.ParseRawPrivateKey(sshKey)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse ssh private key: %w", err)
	}

	ed25519Key, ok := privateKey.(*ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("got %s key type but: %w", reflect.TypeOf(privateKey), UnsupportedKeyType)
	}

	bytes, err := ed25519PrivateKeyToCurve25519(*ed25519Key) //THIS
	if err != nil {
		return nil, err
	}

	s, err := bech32.Encode("AGE-SECRET-KEY-", bytes)
	if err != nil {
		return nil, err
	}
	s = strings.ToUpper(s)
	return &s, nil
}

func main() {

	e, err := readEntity("test-key.asc")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(reflect.TypeOf(e.PrivateKey.PrivateKey))

	castkey, ok := e.PrivateKey.PrivateKey.(*eddsa.PrivateKey)
	if !ok {
		log.Fatal("failed to cast")
	}
	spew.Dump(castkey)

	agekey, err := SSHPrivateKeyToAge(castkey.D, []byte{})

	if err != nil {
		log.Fatal(err)
	}
	log.Println(agekey)

}
