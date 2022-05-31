/*

Gosoap generates Go code from a WSDL file.

Usage: gosoap [clientOption] soapApi.wsdl
  -o string
        File where the generated code will be saved (default "soapApi.go")
  -p string
        Package under which code will be generated (default "soapApi")
  -v    Shows gosoap version

Features

Supports only Document/Literal wrapped services, which are WS-I (http://ws-i.org/) compliant.

Supports WSDL 1.1, XML Schema 1.0, SOAP 1.1.

Resolves external XML Schemas

Supports providing WSDL HTTP URL as well as a local WSDL file.

Not supported

UDDI.

TODO

If WSDL file is local, resolve external XML schemas locally too instead of failing due to not having a URL to download them from.

Resolve XSD element references.

Support for generating namespaces.

*/

package main

import (
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-aegian/gosoap"
)

// Version is initialized in compilation time by go build.
var Version string

// Name is initialized in compilation time by go build.
var Name string

var version = flag.Bool("v", false, "display gosoap version")
var pkg = flag.String("p", "soapProxy", "package name for the soap proxy")
var outFile = flag.String("o", "soap-proxy.go", "output file name for the the soap proxy")
var dir = flag.String("d", "./", "output directory of the soap proxy file")
var insecure = flag.Bool("i", false, "skip TLS verification")
var makePublic = flag.Bool("make-public", true, "generates go types with public/exported")

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.SetPrefix("")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [Option] services.wsdl\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		log.Println(Version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	wsdlPath := os.Args[len(os.Args)-1]

	if *outFile == wsdlPath {
		log.Fatalln("Output file cannot be the same wsdl file")
	}

	builder, err := gosoap.New(wsdlPath, *pkg, *insecure, *makePublic)
	if err != nil {
		log.Fatalln(err)
	}

	soapCode, err := builder.Build()
	if err != nil {
		log.Fatalln(err)
	}

	pkg := filepath.Join(*dir, *pkg)
	err = os.Mkdir(pkg, 0744)

	writeFile(filepath.Join(pkg, *outFile), soapCode["header"])

	writeFile(filepath.Join(pkg, strings.Replace(*outFile, ".", "_types.", 1)), soapCode["types"])

	writeFile(filepath.Join(pkg, strings.Replace(*outFile, ".", "_operations.", 1)), soapCode["operations"])

	log.Println("Done")
}

func writeFile(fileName string, data []byte) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	source, err := format.Source(data)
	if err != nil {
		file.Write(data)
		log.Fatalln(err)
	}

	file.Write(source)

}
