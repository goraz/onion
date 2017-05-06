package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/fzerorubigd/onion.v3"
)

func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	conf := onion.New()

	def := onion.NewDefaultLayer()
	err := def.SetDefault("port", 6998)
	panicOnErr(err)

	err = conf.AddLayer(def)
	panicOnErr(err)

	log.Printf("the port is %d (default layer)", conf.GetInt("port"))

	err = conf.AddLayer(onion.NewFileLayer(filepath.Join(pwd(), "test.json")))
	panicOnErr(err)

	log.Printf("the port is %d (file layer)", conf.GetInt("port"))

	err = conf.AddLayer(onion.NewEnvLayer("PORT"))
	panicOnErr(err)

	log.Printf("the final port is %d, Try to set PORT in env and try again", conf.GetInt("port"))
}
