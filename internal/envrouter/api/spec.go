// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xYUW/bNhD+KwK3RyOy1z0UenOzojAwBEGy7aXIAyOda7YSyZAnB4bh/z6QsizSIm2p",
	"s4091SWPx/u++3h3ypbkopKCA0dNsi3R+Qoqan/OpSxZTpEJbv4rlZCgkIHd5LQC8y9uJJCMaFSMfyO7",
	"CVEghWYo1OYhZvIOryshfgT27Pm3mikoSPa1ueRl0lqJ1++Qo/Fwr6AAjoyW+hlyBfgn07hAqEYE2iwM",
	"CmFvPCiSJ3irQWM/kB+wCcYhqdbvQhXBzVqDiiA4CtS4d+wdv6GwP/M1U4JXwHEwZUOTs+AaKc+h75j6",
	"iuqhzUVVMXxe0eAu+DH39k9IcnlSAcDryiAqQJZiY/13uCLo7e6klYcb2sSDeYqgR1FcmqNcAUUo/mIR",
	"Kn6aQ7miOsYuLVxlvwpRAuWniNdIFUIRPrTfjEAIP884/218LQKfoi6UUJ6eYPmJ8cLcOzpN54gOU3OE",
	"7jSsZSTotv72g86Pi1UwNKYf4D2cm6g8alUOTpaxnQRi6aMxDhhfCuO6AJ0rJhvCTf1SokZQyfxxYeoz",
	"wxIC62tQujkxu5veTU2oQgKnkpGMfLib3n2wlRJXlp+USpauZ6nDtF3/1jBlqLSLi4Jk5AvgXLJ/ZnPX",
	"2ADWUnDdEP7bdGp5Fxz3QnBcp991I56m55pfDKGyB3+1+iC/pF13TvetOXX78qGKEaoU3TSU+VR9AUxK",
	"pjERy8QDtrPGIczp1qRqZwVUY5/9v2VBEZK5J0ifnce6z85D25gUrQBBaZJ93RJmPJoUtOU0a4XSaQdV",
	"DROHp2OdvTTGoPGTaIrRYMoHM71r5PwfsjviquMkPtd5Dvo4Z703dF6t9/0jt9BsfGAbqeAAZHtCCh0Q",
	"6r2t9Ul3e6KbWtOTq9An+bm8tqKDY1Bos+vd22Xi52XnFIwCSkDoC/APux7h+MqVwaPy975MwjCd7nv+",
	"XX12jW/xotwJfuQb8oD5mFk3mp7HvHCNb4HZHZxHYm6BxQAPR3tbqJfCqQ6z7HmkT45t+FG+1aA23av0",
	"J9P4Y5yEj/tj7ujjhsVxJeAK+XK+FUZmzD15tpM5xtEW5qfvGr3LD/ma3eqYnCF1+/AXKAZDtO4Y30Yo",
	"h++zgUKhZZl4mIyZrquKmm+8iMnlJqIeQ9fQU8fJtfXk3jReT+MGHZe6/9GIY+QDat3GYL/dyQpRZmla",
	"ipyWK6Ex+zj9OE3J7mX3bwAAAP//2tVD7SMWAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

