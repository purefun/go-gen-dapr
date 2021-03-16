package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/purefun/go-gen-dapr/generator"
)

var (
	pkg      string
	pkgValue = "."
	pkgUsage = "the package contains the type"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "\tgo-gen-dapr [flags] interface")
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
}

func Execute() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.StringVar(&pkg, "pkg", pkgValue, pkgUsage)

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	typeName := args[0]

	g := generator.NewGenerator(generator.Options{
		ServicePkg:  pkg,
		ServiceType: typeName,
		GenComment:  true,
	})

	out, err := g.Generate()

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(pkg+"/service.dapr.go", []byte(out), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
