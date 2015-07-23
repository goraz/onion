package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fzerorubigd/onion"
)

func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func main() {

	conf := onion.New()

	err := conf.AddLayer(onion.NewFileLayer(filepath.Join(pwd(), "test.json")))
	if err != nil {
		log.Fatal(err)
	}

	err = conf.AddLayer(onion.NewEnvLayer("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("the port is %d", conf.GetInt("port", 0))
}
