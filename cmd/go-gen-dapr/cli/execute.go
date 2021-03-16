package cli

import (
	"flag"
	"fmt"
	"github.com/purefun/go-gen-dapr/generator"
	"log"
	"os"
)

var (
	pkg      string
	pkgValue string = "."
	pkgUsage string = "the package contains the type"
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
		PackageName: "echo",
		ServicePkg:  pkg,
		ServiceType: typeName,
		GenComment:  true,
	})

	out, err := g.Generate()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}
