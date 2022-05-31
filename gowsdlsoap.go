package gowsdlsoap

import "github.com/go-aegian/gowsdlsoap/builder"

// New creates the builder.
func New(file, pkg string, ignoreTLS bool, exportAllTypes bool) (*builder.Builder, error) {
	return builder.New(file, pkg, ignoreTLS, exportAllTypes)
}
