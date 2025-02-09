package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/TicketsBot/export/example/utils"
	"github.com/TicketsBot/export/pkg/validator"
	"os"
)

var (
	keyPath = flag.String("key", "", "Path to the public key file")
	zipPath = flag.String("zip", "", "Path to the zip file")
)

func main() {
	flag.Parse()

	if *keyPath == "" || *zipPath == "" {
		flag.PrintDefaults()
		return
	}

	key, err := utils.LoadPublicKeyFromDisk(*keyPath)
	if err != nil {
		panic(err)
	}

	b, err := os.ReadFile(*zipPath)
	if err != nil {
		panic(err)
	}

	v := validator.NewValidator(key,
		validator.WithMaxUncompressedSize(100*1024*1024),
		validator.WithMaxIndividualFileSize(100*1024*1024))

	reader := bytes.NewReader(b)
	output, err := v.ValidateGuildData(reader, reader.Size())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", output)
}
