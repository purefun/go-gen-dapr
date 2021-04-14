package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/purefun/go-gen-dapr/generator"
	"golang.org/x/tools/go/packages"
)

var (
	pkg      string
	pkgValue = "." // current package
	pkgUsage = "the package contains the type"

	target      string
	targetValue = "service"
	targetUsage = "service, pubsub, subscriber"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "\tgo-gen-dapr [flags] [interface]")
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
}

func Execute() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.StringVar(&pkg, "pkg", pkgValue, pkgUsage)
	flag.StringVar(&target, "target", targetValue, targetUsage)

	flag.Parse()

	var out string
	var err error

	switch target {
	case "pubsub":
		out, err = generator.GeneratePublishEvents(generator.PublishEventsOptions{
			Pkg:        pkg,
			GenComment: true,
		})
	case "service", "subscriber":
		args := flag.Args()
		if len(args) == 0 {
			fmt.Println("go-gen-dapr SomeInterfaceName")
			flag.Usage()
			os.Exit(1)
		}
		generateType := generator.GenerateTypeService
		if target == "subscriber" {
			generateType = generator.GenerateTypeSubscriber
		}
		g := generator.NewGenerator(generator.Options{
			ServicePkg:   pkg,
			ServiceType:  args[0],
			GenComment:   true,
			GenerateType: generateType,
		})
		out, err = g.Generate()
	default:
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}

	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedFiles}, pkg)
	if err != nil {
		log.Fatal("load package failed: ", err)
	}
	file := pkgs[0].GoFiles[0]
	outFile := filepath.Join(filepath.Dir(file), target+".dapr.go")

	err = ioutil.WriteFile(outFile, []byte(out), 0644)
	if err != nil {
		log.Fatal("write file err: ", err)
	}
}
