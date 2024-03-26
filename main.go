package main

import (
	"crypto/sha512"
	"fmt"
	"log"
	"os"
  "syscall"
	"strings"
  "encoding/pem"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
  "github.com/ProtonMail/go-crypto/openpgp/eddsa"
	"github.com/ProtonMail/go-crypto/openpgp/packet"

	"crypto/ed25519"
	"errors"
	"github.com/Mic92/ssh-to-age/bech32"
	"github.com/davecgh/go-spew/spew"
  "golang.org/x/term"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ssh"
	// "bytes"
	// "golang.org/x/crypto/ssh"
	// "golang.org/x/crypto/curve25519"
	// "https://pkg.go.dev/crypto/ed25519#PrivateKey
	// "crypto/ed25519"ccc1be8d-24dc-41ad-9d66-b657711419d7
	"reflect"
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

func SSHPrivateKeyToAge(bytes []byte) (*string, error) {
	s, err := bech32.Encode("AGE-SECRET-KEY-", bytes)
	if err != nil {
		return nil, err
	}
	s = strings.ToUpper(s)
	return &s, nil
}

func main() {
  // TODO turn these into CLI inputs
	keyfile := "./priv-gpg"

	e, err := readEntity(keyfile)
	if err != nil {
		log.Fatal(err)
	}

	spew.Config.MaxDepth = 2
	spew.Config.Indent = "     "

  log.Println("Keys:")
  log.Println("[0]", e.PrimaryKey.KeyIdString() + " (primary)")
  for i := 0; i < len(e.Subkeys); i++ {
    log.Println(fmt.Sprintf("[%d]", i + 1), e.Subkeys[i].PublicKey.KeyIdString() + " (subkey)")
  }

  log.Println("Please choose a key by index (default: 0):")

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
    var subkey = e.Subkeys[keyIndex - 1]
    log.Println(fmt.Sprintf("Continuing with key [%d]", keyIndex), subkey.PublicKey.KeyIdString())
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
  castkey, ok := targetKey.PrivateKey.(*eddsa.PrivateKey)
  if !ok {
    log.Fatal("failed to cast")
  }
	// spew.Dump(castkey)

  // get public key as OpenSSH key
  log.Println("public key type:", reflect.TypeOf(castkey.PublicKey))
	var pubkey ed25519.PublicKey = castkey.PublicKey.X

	agePub, err := bech32.Encode("age", pubkey)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("public age key:", string(agePub))

  sshPub, err := ssh.NewPublicKey(pubkey)
  if err != nil {
    log.Fatal(err)
  }
  // log.Println("public key SSH key wire format:", sshPub.Marshal())
  // log.Println("public key SHA256:", ssh.FingerprintSHA256(sshPub))
  log.Println("public SSH key:", string(ssh.MarshalAuthorizedKey(sshPub)))

	// TODO: are these the correct bytes?
	var privkey ed25519.PrivateKey = castkey.D
	// var privkey ed25519.PrivateKey = castkey.MarshalByteSecret()
  // var privkey = ed25519.NewKeyFromSeed(castkey.D)

  // TODO is this right?
	bytes, err := ed25519PrivateKeyToCurve25519(privkey)
	if err != nil {
		log.Fatal(err)
	}

  // TODO trying to get private key as age key
	agekey, err := SSHPrivateKeyToAge(bytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*agekey)

  // TODO trying to get private key as OpenSSH key
  privPem, err := ssh.MarshalPrivateKey(privkey, "")
  if err != nil {
    log.Fatal(err)
  }
  if err := pem.Encode(os.Stdout, privPem); err != nil {
		log.Fatal(err)
	}

  // TODO make sure public key is still the same
  var priv ed25519.PrivateKey = bytes
  var pubkey2 = priv.Public()
  // var pubkey2 = privkey.Public()
  sshPub2, err := ssh.NewPublicKey(pubkey2)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("verify public SSH key:", string(ssh.MarshalAuthorizedKey(sshPub2)))
}
