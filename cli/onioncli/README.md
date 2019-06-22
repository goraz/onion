# Onion CLI 

Onion cli is a simple tool to move configuration from a layer to another layer. 

Currently it can read/write data from file and etcd keys. also it support for encrypt/decrypt data using PGP. 

## Installation 

```
go get -u github.com/fzerorubigd/onion/cli/onioncli
```

if you want to encrypt/decrypt data using PGP, you need to create a private/public key pair using gpg (or any other tool)
for example this is a fast way to create a `TEST` key pair (not protected with password) : 

```bash
export EMAIL="joe@foo.bar"
export NAME="app"
export GNUPGHOME="$(mktemp -d)"
cat >foo <<EOF
     %echo Generating a basic OpenPGP key
     Key-Type: default
     Subkey-Type: default
     Name-Real: ${NAME}
     Name-Comment: app configuration key, no passphrase
     Name-Email: ${EMAIL}
     Expire-Date: 0
     %no-protection
     # Do a commit here, so that we can later print "done" :-)
     %commit
     %echo done
EOF
gpg --batch --generate-key foo
gpg --export --armor "${EMAIL}" > .pubring.gpg
gpg --export-secret-keys --armor "${EMAIL}" > .secring.gpg

```

This should create two file, `.pubring.gpg` and `.secring.gpg` contains your testing (respectively) public and private keys. 

## Usage 

for testing, create a plain `config.yaml` 
```bash
cat > config.yaml <<EOF
---
example: string
number: 100
EOF

```

Read the file and encrypt it using PGP and print the result in stdout :

```bash 
onioncli -s config.yaml -d- --pk=.pubring.gpg
```

Read the file and put data in `/app/data` key in etcd (make sure you have an etcd instance running) :

```bash
onioncli -s config.yaml -d etcd://127.0.0.1:2379/app/data --pk=.pubring.gpg
```

if the `--pk` passed to the cli, then cli encrypt the data before putting it into the destination. 

Read the data from etcd and show it in stdout : 

```bash
onioncli -s etcd://127.0.0.1:2379/app/data -d-
```

If you want to see the actual data (not the base64/PGP encrypted data) you should provide the secret key with `--sk` flag :
```bash
onioncli -s etcd://127.0.0.1:2379/app/data -d- --sk=.secring.gpg
```

in `-s` and `-d` you can use `-` for stdin and stdout.

