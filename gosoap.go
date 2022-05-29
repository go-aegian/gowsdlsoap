package gosoap

import "github.com/go-aegian/gosoap/builder"

// New creates the builder.
func New(file, pkg string, ignoreTLS bool, exportAllTypes bool) (*builder.Builder, error) {
	return builder.New(file, pkg, ignoreTLS, exportAllTypes)
}
