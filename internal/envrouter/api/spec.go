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

	"H4sIAAAAAAAC/8xYwW7jNhD9FYPt0Yicbg8L3bzpIjBQBEHS9rLIgZHGEXclkuGMHBiG/70gZcWURNnS",
	"xjb2lJgckvPevBkOtWGJKrSSIAlZvGGYZFBw9+9c61wknISS9qc2SoMhAW5S8gLsX1prYDFDMkK+sO2U",
	"GdAKBSmzvuszeYPnTKkfgTm3/rUUBlIWf6sOeZrWVur5OyRkd7hRRSGo6xUvKVMmeGgBiPwl7BBmPDhO",
	"ogAkXujjrtotgp4aSEGS4Dk+QmKA/hZIC4JiBKXVwCCydsaDPHmA1xIwwOIPWAf90BzxTZk0OFkimB4E",
	"LUft9p69t2/I7a9yJYySBUgaTNlQGS0kEpcJBITU1H4HbeIE+NijG2j63Jk/kDzLgwoAWRYWUQo6V2u3",
	"/x5XD3o3O63l4bs2bcA8RNC9Sk/NUWKAE6T/iB4qfppDzU1dzARBgQcIZdwYvnaLMo59IeGpnw7PSuXA",
	"5aFoIXFDkIYX7SZ7cIdzuj9otX81giave1dCwX2A5RchU3vu6Ngei06Ymha6w7CWPU7X10vX6aRd4YKu",
	"CbyDt3BsejVVmnxwsKztNOBLF43dQMilslungIkRuiLcFj2jSgIzmd8vbFEXlENgfAUGqxXXV7OrmXVV",
	"aZBcCxazT1ezq0+uvFLm+Im4FtHqOvKYduMvFVOWSje4SFnMboHmWvx3PfeNLWDUSmJF+B+zmeNdSdoJ",
	"wds6+o6VeKqWopGPvzt9sN+iffMR7TqPyG87OplqKWtSdQs0yQXSRC0nDWBbZxzCHG1sqLZOQCV12f9X",
	"p5xgMm8IssnOfdll566+zQwvgMAgi79tmLA72hDUNTiuhbLXDpkSph5PbZ09VcaA9EVVxWgw5YOZ3lZy",
	"/kB0RxzVDuJjmSSA7Zh1cui4Wm+6Sy6h2f4ub6SCA5DdCq0wINQbV+sn+9MnWNWajlwVHuTn9Nrq7TaD",
	"Qrs+37n7SPy87LyCkUIOBF0B/uXGezg+c2VoUPlnVyZhmN7tezyvvvrGl8gov+0fmUMNYE3ML4JaN4H3",
	"y0ZpG1VdK0YbzPj2KC+3gvwrYN7crXqmom2Bh0Qfnd3w4E+Du7QQfVBOp7sCdm/2wWko9m+P4/pc+MaX",
	"0Kf/MhqpzxpYH+DhaC8L9VQ4zfu74zjSB882nEKvJZh1UP1sSO60ljefJKOXWxbPmV+D4uW960ZGzF95",
	"tOvwjHvbjWb4ztFnNF0+Z2fRJmdIDXv/GCpgiNY948sI5f0tPVAoPM8nDUzWDMui4PY93mNyuu61w9A5",
	"9LTn5Nx68k8ar6dxTalP3S/Ujlr5gFnVPrjvLCwj0nEU5SrheaaQ4s+zz7OIbZ+2/wcAAP//F1Eu964Y",
	"AAA=",
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

