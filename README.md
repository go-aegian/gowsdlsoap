# gowsdlsoap

[![GoDoc](https://godoc.org/github.com/go-aegian/gowsdlsoap?status.svg)](https://godoc.org/github.com/go-aegian/gowsdlsoap)

Generates GO types based structs for a given service wsdl file, it provides a proxy http client to make request to the given service.

Supports file attachments.

Supports NTLM and Basic Auth authentication methods.

### Install

* [Download binary release](https://github.com/go-aegian/gowsdlsoap/releases)
* Download and build locally: `go get github.com/go-aegian/gowsdlsoap/...`
* Install from go: `go install github.com/go-aegian/gowsdlsoap/...`
* Uninstall: `go clean -i -n github.com/go-aegian/gowsdlsoap`

### Goals
* Generate go code for the wsdl definition
* Support only Document/Literal wrapped services, which are [WS-I](http://ws-i.org/) compliant
* Support:
	* WSDL 1.1
	* XML Schema 1.0
	* SOAP 1.1
* Resolve external XML Schemas
* Support external and local WSDL

### Caveats
* Please keep in mind that the generated code is just a reflection of what the WSDL is like. If your WSDL has duplicated type definitions, your Go code is going to have the same and may not compile.

### Usage
```
Usage: gowsdlsoap [options] services.wsdl
  -o string
        File where the generated code will be saved (default "services-proxy.go")
  -p string
        Package under which code will be generated (default "servicesProxy")
  -i    Skips TLS Verification
  -v    Shows gowsdlsoap version
  ```
