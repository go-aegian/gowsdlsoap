package builder

import (
	"net/url"
	"path/filepath"
)

// A location encapsulate information about the loc of WSDL/XSD.
// It could be either URL or an absolute file path.
type location struct {
	url  *url.URL
	file string
}

// NewLocation parses a raw location into a location structure.
// If path is URL then it should be absolute.
// If path is a file then relative file path will be converted into an absolute path.
func NewLocation(path string) (*location, error) {
	u, _ := url.Parse(path)
	if u.Scheme != "" {
		return &location{url: u}, nil
	}

	absURI, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &location{file: absURI}, nil
}

// Parse parses path in the context of the receiver. The provided path may be relative or absolute.
func (r *location) Parse(ref string) (*location, error) {
	if r.url != nil {
		u, err := r.url.Parse(ref)
		if err != nil {
			return nil, err
		}
		return &location{url: u}, nil
	}

	if filepath.IsAbs(ref) {
		return &location{file: ref}, nil
	}

	if u, err := url.Parse(ref); err == nil {
		if u.Scheme != "" {
			return &location{url: u}, nil
		}
	}

	return &location{file: filepath.Join(filepath.Dir(r.file), ref)}, nil
}

func (r *location) IsFile() bool {
	return r.file != ""
}

func (r *location) IsURL() bool {
	return r.url != nil
}

func (r *location) String() string {
	if r.IsFile() {
		return r.file
	}
	if r.IsURL() {
		return r.url.String()
	}
	return ""
}
