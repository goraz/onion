package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/goraz/onion/ciphers/secconf"
	"github.com/ogier/pflag"
	"go.etcd.io/etcd/client"
)

var (
	src    = pflag.StringP("source", "s", "", "Source address to read from")
	dst    = pflag.StringP("destination", "d", "", "Destination address to write into")
	srcKey = pflag.String("sk", "", "Source private key to use for reading data from source, if the source is plain leave it empty")
	dstKey = pflag.String("pk", "", "Destination public key to use for writing data to destination, if the destination is plain leave it empty")
)

func open(path string) (io.ReadCloser, error) {
	if path == "-" {
		return os.Stdin, nil
	}

	return os.Open(path)
}

func create(path string) (io.WriteCloser, error) {
	if path == "-" {
		return os.Stdout, nil
	}

	return os.Create(path)
}

func readAllFile(path string) ([]byte, error) {
	f, err := open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return ioutil.ReadAll(f)
}

func writeAllFile(path *url.URL, data []byte) error {
	f, err := create(path.Path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = fmt.Fprint(f, string(data))
	return err
}

func connectEtcd(u *url.URL) (client.Client, client.KeysAPI, error) {
	cli, err := client.New(client.Config{
		Endpoints: []string{"http://" + u.Host},
	})
	if err != nil {
		return nil, nil, err
	}

	kv := client.NewKeysAPI(cli)
	return cli, kv, nil
}

func readAllEtcd(u *url.URL) ([]byte, error) {
	_, kv, err := connectEtcd(u)
	if err != nil {
		return nil, err
	}
	resp, err := kv.Get(context.TODO(), u.Path, nil)
	if err != nil {
		return nil, err
	}

	return []byte(resp.Node.Value), nil
}

func writeAllEtcd(u *url.URL, data []byte) error {
	_, kv, err := connectEtcd(u)
	if err != nil {
		return err
	}
	if _, err = kv.Set(context.TODO(), u.Path, string(data), nil); err != nil {
		return err
	}

	return nil
}

func read(path *url.URL) ([]byte, error) {
	switch path.Scheme {
	case "":
		return readAllFile(path.Path)
	case "etcd":
		return readAllEtcd(path)
	default:
		return nil, fmt.Errorf("scheme %q is not valid", path.Scheme)
	}
}

func write(path *url.URL, data []byte) error {
	switch path.Scheme {
	case "":
		return writeAllFile(path, data)
	case "etcd":
		return writeAllEtcd(path, data)
	default:
		return fmt.Errorf("scheme %q is not valid", path.Scheme)
	}
}

type transformer func([]byte) ([]byte, error)

func encrypt(fl string) (transformer, error) {
	if fl == "" {
		return func(in []byte) ([]byte, error) {
			return in, nil
		}, nil
	}
	data, err := readAllFile(fl)
	if err != nil {
		return nil, err
	}
	return func(in []byte) ([]byte, error) {
		return secconf.Encode(in, bytes.NewReader(data))
	}, nil
}

func decrypt(fl string) (transformer, error) {
	if fl == "" {
		return func(in []byte) ([]byte, error) {
			return in, nil
		}, nil
	}
	data, err := readAllFile(fl)
	if err != nil {
		return nil, err
	}
	return func(in []byte) ([]byte, error) {
		return secconf.Decode(in, bytes.NewReader(data))
	}, nil
}

func fatalIfErr(message string, err error) {
	if err == nil {
		return
	}

	log.Fatalf(message, err)

}

func main() {
	pflag.Parse()
	if *src == "" || *dst == "" {
		pflag.Usage()
		return
	}

	srcUrl, err := url.Parse(*src)
	fatalIfErr("Parsing source url failed: %q", err)

	dstUrl, err := url.Parse(*dst)
	fatalIfErr("Parsing destination url failed: %q", err)

	enc, err := encrypt(*dstKey)
	fatalIfErr("Fail to create the encrypt function: %q", err)

	dec, err := decrypt(*srcKey)
	fatalIfErr("Fail to create the decrypt function: %q", err)

	in, err := read(srcUrl)
	fatalIfErr("Failed to read from source: %q", err)

	in, err = dec(in)
	fatalIfErr("Failed to decrypt the source: %q", err)

	out, err := enc(in)
	fatalIfErr("Failed to encrypt the data: %q", err)

	err = write(dstUrl, out)
	fatalIfErr("Failed to write data into destination: %q", err)
}
