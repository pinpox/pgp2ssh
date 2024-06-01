# pgp2ssh

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
❯ go build
❯ ./pgp2ssh
```

**Nix/NixOS Users**

A flake is provided for Nix users. Just use `nix run` instead of building and
running manually.

It'll ask you for the path to your private PGP key, followed by choosing the key/subkey and if your PGP key is encrypted it'll ask for the passphrase.

In the output, verify that the public SSH key printed matches the one in `metadata.json`.
If it matches, the last part of the output it will print the matching private SSH key.
You can save the key to a file and use how you want.

### Example: Decrypt age files

If you want to decrypt a file that was encryptd by `age` with your public SSH key, you can just use `age` as normal to decrypt the file using the SSH private key that we've got in the previous step:

```sh
❯ age --decrypt --identity ./ssh-secret-key --output decrypted ./testfile.txt.age
```

## Troubleshooting

If the conversion fails with the error:

```
2024/03/27 22:09:09 openpgp: invalid data: user ID signature with wrong type
```

You might be missing the private key of your subkeys. When running `gpg -K` you
should **NOT** see a `>` infront of the keys like this:

```
ssb>  ed25519/0xB68746238E59B548 2018-07-09 [S] [expires: 2026-01-02]
      Keygrip = C89E5AABCBF7142DBC26E68FB3121DE12DCBF4FF
ssb>  cv25519/0x65CD5E0200C56C17 2018-07-09 [E] [expires: 2026-01-02]
      Keygrip = 867EA9F6ADBEBE18ED98253B884F53CBD53C526B
ssb>  ed25519/0xF36CF32DF9B09855 2018-07-09 [A] [expires: 2026-01-02]
      Keygrip = 553D56865642B05AB3C5B62DC68795691702B960
```
The `>` (corner of a card) indicates, that the private part is on a smart card
or not available. This may also be caused by expired keys. For possible
solutions see https://github.com/pinpox/pgp2ssh/issues/6


### Support & Donations

This project was built with lots of headaches by [pinpox](https://github.com/pinpox/) & [felschr](https://github.com/felschr/). If you need help, feel free to contact us.

And if you want to thank us, you can send us any crypto or token to our Ethereum / Polygon wallets 😊:  
pinpox: `0xde031f16976AFcaC613087B6213Eb521F63d3A49`
felschr: `0xD66753D737603E18018281E298Df86DE402d313E`

<a href="https://www.buymeacoffee.com/pinpox"><img src="https://img.buymeacoffee.com/button-api/?text=Buy me a coffee&emoji=😎&slug=pinpox&button_colour=82aaff&font_colour=000000&font_family=Inter&outline_colour=000000&coffee_colour=FFDD00"></a>
