
Example key provided in `./gnupg`

```
❯ gpg --homedir ./gnupg_testkey/ -K
/home/pinpox/code/github.com/pinpox/gpg2age/./gnupg_testkey/pubring.kbx
-----------------------------------------------------------------------
sec   ed25519 2024-03-25 [C]
      9FE4D484B69DB9F5C7AA208E76188CF30717B54E
uid           [ultimate] test (test) <test@test.com>
ssb   ed25519 2024-03-25 [S]
ssb   cv25519 2024-03-25 [E]
ssb   ed25519 2024-03-25 [A]
```

### Get age key and encrypt test file

```sh
❯ gpg --homedir ./gnupg --export-ssh-key 9FE4D484B69DB9F5C7AA208E76188CF30717B54E
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICYvKXGcG4a19tTT0Rycbn+D0r/YlKltLJ9dY2gR/Fjx openpgp:0x47C9F3FF
```

```sh
❯ gpg --homedir ./gnupg --export-ssh-key 9FE4D484B69DB9F5C7AA208E76188CF30717B54E | ssh-to-age                                                                                                                             impure ❄ ssh-to-age
age18s8m9hvlrwvltgys4lafyyqe356ntc7e06t4kd2nccqm5amsaa2s878mju # saved as age-public-key
```

```sh
❯ age --encrypt -R age-public-key testfile.txt > testfile.txt.age
```

### Try to decrypt

```sh
❯ age --decrypt --identity age-secret-key --output decrypted testfile.txt.age                                                                                                                                          impure ❄ ssh-to-age age
age: error: no identity matched any of the recipients
age: report unexpected or unhelpful errors at https://filippo.io/age/report
```

FAIL :(
