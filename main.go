package main

import (
	"crypto/sha512"
	"fmt"
	"log"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/eddsa"
	"github.com/ProtonMail/go-crypto/openpgp/packet"

	"strings"

	"crypto/ed25519"
	"errors"
	"github.com/Mic92/ssh-to-age/bech32"
	"golang.org/x/crypto/curve25519"
	// "github.com/davecgh/go-spew/spew"
	// "bytes"
	// "golang.org/x/crypto/ssh"
	// "golang.org/x/crypto/curve25519"
	// "reflect"
	// "filippo.io/edwards25519"
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

func SSHPrivateKeyToAge(privatekey ed25519.PrivateKey, passphrase []byte) (*string, error) {

	bytes, err := ed25519PrivateKeyToCurve25519(privatekey)
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

	// log.Println(reflect.TypeOf(e.PrivateKey.PrivateKey))

	castkey, ok := e.PrivateKey.PrivateKey.(*eddsa.PrivateKey)
	if !ok {
		log.Fatal("failed to cast")
	}
	// spew.Dump(castkey)

	// TODO: Not sure if these are the correct bytes ??????
	agekey, err := SSHPrivateKeyToAge(castkey.D, []byte{})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*agekey)

}
