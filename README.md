# onion

[![Build Status](https://travis-ci.org/goraz/onion.svg)](https://travis-ci.org/goraz/onion)
[![Coverage Status](https://coveralls.io/repos/goraz/onion/badge.svg?branch=develop&service=github)](https://coveralls.io/github/goraz/onion?branch=master)
[![GoDoc](https://godoc.org/github.com/goraz/onion?status.svg)](https://godoc.org/github.com/goraz/onion)
[![Go Report Card](https://goreportcard.com/badge/github.com/goraz/onion)](https://goreportcard.com/report/github.com/goraz/onion)

    import "github.com/goraz/onion"

Package onion is a layer based, pluggable config manager for golang.

The current version in `develop` branch is work in progress (see the [milestone](https://github.com/goraz/onion/milestone/1)), for older versions check the `v2` and `v3` branches and use the `gopkg.in/goraz/onion.v1` and `gopkg.in/goraz/onion.v2`
For the next release we use the go module and tagging using semantic version.

```
Shrek: For your information, there's a lot more to ogres than people think.
Donkey: Example?
Shrek: Example... uh... ogres are like onions! 
[holds up an onion, which Donkey sniffs] 
Donkey: They stink? 
Shrek: Yes... No! 
Donkey: Oh, they make you cry? 
Shrek: No! 
Donkey: Oh, you leave 'em out in the sun, they get all brown, start sproutin' little white hairs...
Shrek: [peels an onion] NO! Layers. Onions have layers. Ogres have layers... You get it? We both have layers.
[walks off]
Donkey: Oh, you both have LAYERS. Oh. You know, not everybody like onions. CAKE! Everybody loves cake! Cakes have layers!
Shrek: I don't care what everyone likes! Ogres are not like cakes.
Donkey: You know what ELSE everybody likes? Parfaits! Have you ever met a person, you say, "Let's get some parfait," they say, "Hell no, I don't like no parfait."? Parfaits are delicious!
Shrek: NO! You dense, irritating, miniature beast of burden! Ogres are like onions! End of story! Bye-bye! See ya later.
Donkey: Parfait's gotta be the most delicious thing on the whole damn planet! 
```
## Goals 

The main goal is to have minimal dependency based on usage. if you need normal config files in the file system, 
there should be no dependency to `etcd` or `consul`, if you have only `yaml` files, including `toml` or any other format 
is just not right.

## Usage 

Choose the layer first. normal file layer and json are built-in but for any other type 
you need to import the package for that layer. 

### Example json file layer 

```go
package main

import (
	"fmt"

	"github.com/goraz/onion"
)

func main() {
	// Create a file layer to load data from json file. onion loads the file based on the extension.
	// so the json file should have `.json` ext.
	l1, err := onion.NewFileLayer("/etc/shared.json", nil)
	if err != nil {
		panic(err)
	}

	// Create a layer based on the environment. it loads every environment with APP_ prefix
	// for example APP_TEST_STRING is available as o.Get("test.string")
	l2 := onion.NewEnvLayerPrefix("_", "APP")

	// Create the onion, the final result is union of l1 and l2 but l2 overwrite l1.
	o := onion.New(l1, l2)
	str := o.GetStringDefault("test.string", "empty")
	fmt.Println(str)
	// Now str is the string in this order
	// 1- if the APP_TEST_STRING is available in the env
	// 2- if the shared.json had key like this { "test" : { "string" : "value" }} then the str is "value"
	// 3- the provided default, "empty"
}
```

### Loading other file format 

Currently `onion` support `json` format out-of-the-box, while you need to blank import the loader package of others formats to use them:
* `toml` (for 0.4.0 version)
* `toml-0.5.0` (for 0.5.0 version)
* `yaml`
* `properties`

For example:
```go 
import (
    _ "github.com/goraz/onion/loaders/toml" // Needed to load TOML format
)
``` 

### Watch file and etcd

Also there is other layers, (like `etcd` and `filewatchlayer`) that watches for change. 

```go
package main

import (
	"fmt"

	"github.com/goraz/onion"
	"github.com/goraz/onion/layers/etcdlayer"
	"github.com/goraz/onion/layers/filewatchlayer"
)

func main() {
	// Create a file layer to load data from json file. also it watches for change in the file
	l1, err := filewatchlayer.NewFileWatchLayer("/etc/shared.json", nil)
	if err != nil {
		panic(err)
	}

	l2, err := etcdlayer.NewEtcdLayer("/app/config", "json", []string{"http://127.0.0.1:2379"}, nil)
	if err != nil {
		panic(err)
	}

	// Create the onion, the final result is union of l1 and l2 but l2 overwrite l1.
	o := onion.New(l1, l2)
	// Get the latest version of the key 
	str := o.GetStringDefault("test.string", "empty")
	fmt.Println(str)
}
```

### Encrypted config 

Also if you want to store data in encrypted content. currently only `secconf` (based on the [crypt](https://github.com/xordataexchange/crypt) project)
also the [onioncli](https://github.com/goraz/onion/tree/develop/cli/onioncli) helps you to manage this keys. 

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/goraz/onion"
	"github.com/goraz/onion/ciphers/secconf"
	"github.com/goraz/onion/layers/etcdlayer"
	"github.com/goraz/onion/layers/filewatchlayer"
)

// Normally this should be in a safe place, not here
const privateKey = `PRIVATE KEY`

func main() {
	// The private key should be in the safe place. this is just a demo, also there is a cli tool
	// to create this `go get -u github.com/goraz/onion/cli/onioncli`
	cipher, err := secconf.NewCipher(bytes.NewReader([]byte(privateKey)))
	if err != nil {
		panic(err)
	}

	// Create a file layer to load data from json file. also it watches for change in the file
	// passing the cipher to this make means the file in base64 and pgp encrypted
	l1, err := filewatchlayer.NewFileWatchLayer("/etc/shared.json", cipher)
	if err != nil {
		panic(err)
	}

	// Create a etcd layer. it watches the /app/config key and it should be json file encoded with
	// base64 and pgp
	l2, err := etcdlayer.NewEtcdLayer("/app/config", "json", []string{"http://127.0.0.1:2379"}, cipher)
	if err != nil {
		panic(err)
	}

	// Create the onion, the final result is union of l1 and l2 but l2 overwrite l1.
	o := onion.New(l1, l2)
	// Get the latest version of the key
	str := o.GetStringDefault("test.string", "empty")
	fmt.Println(str)
}
```
