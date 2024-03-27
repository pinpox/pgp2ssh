# gpg2age

Derive private ed25519 SSH key from private PGP key.

GPG itself only supports exporting _public_ SSH keys and other tools don't work for ed25519 keys.

##### Notes:

- A tool exists to do this for RSA keys: [openpgp2ssh](https://manpages.ubuntu.com/manpages/xenial/man1/openpgp2ssh.1.html) but it does not seem to support `ed25519` keys
- Work on `gnupg` was started for this feature, but never finished see this
  issue and commit: https://dev.gnupg.org/T6647

## Instructions

First you need to export your PGP key from GPG:

```sh
❯ gpg2 --export-secret-keys --armor test@test.test >priv-gpg
```

Then identify the public SSH key that was used to encrypt your secret.
You can search for your GitHub username in: https://fluence-dao.s3.eu-west-1.amazonaws.com/metadata.json

If you have multiple subkeys, usually it is the authenticate key highlighted with `[A]` in the output of:

```sh
❯ gpg --list-secret-keys --with-keygrip
```

### Derive private SSH key

```sh
❯ go run main.go
```

It'll ask you for the path to your private PGP key, followed by choosing the key/subkey and if your PGP key is encrypted it'll ask for the passphrase.

In the output, verify that the public SSH key printed matches the one in `metadata.json`.
If it matches, the last part of the output it will print the matching private SSH key.
You can save the key to a file and use how you want.

### Example: Decrypt age files

If you want to decrypt a file that was encryptd by `age` with your public SSH key, you can just use `age` as normal to decrypt the file using the SSH private key that we've got in the previous step:

```sh
❯ age --decrypt --identity ./ssh-secret-key --output decrypted ./testfile.txt.age
```