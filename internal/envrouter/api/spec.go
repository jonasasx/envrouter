// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
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

	"H4sIAAAAAAAC/8xYwW7jNhD9FUPt0Yicbg8L37zpIjBQBIHT9rLIgZHGMXclkeGMEhiG/70gaUWkRFnS",
	"xjb2lFgaDue9eTMcahclIpeigIIwmu8iTDaQM/PvQsqMJ4y4KPRPqYQERRzMy4LloP/SVkI0j5AUL56j",
	"/TRSIAVyEmp712XyBk8bIX4E3pn1LyVXkEbzb3aTx2llJZ6+Q0Law43Ic07tqFhJG6GCm+aAyJ7DAeGG",
	"BZ8TzwGJ5bI/VO0iGKmCFAriLMMHSBTQ3xxpSZCPoNQ+GETWwXhQJCt4KQEDLP6AbTAOyRDfhEqDL0sE",
	"1YGgEah279g7fkNhfy1euRJFDgUNpmyojJYFEisSCAjJ134LbWIE+NChG/Bjbr0/UjzrowqAosw1ohRk",
	"JrbGf42rA715O63k4YY29WAeI+hepKfmKFHACNJ/eAcVP82hZKpqZpwgxyOERkwptjWLNgy7UsJStxye",
	"hMiAFceyhcQUQRpedHjZgTtc091Jq+KrEPi81qGEkruy4ftJTd476+8GXvRbXB8Q8eF0iA/99wgH9THQ",
	"j1L78FZU2ukK+wsvUu1otCT7RBVG0wj3eDbWHUG7dDQobzbmYGgc7+AtLKnOUihVNlhj2nYaiKWNRjvg",
	"xVpo1ylgori0hOterURJoCaL+6U+izhlEHj+Cgrtiuur2dVMhyokFEzyaB59uppdfTKnAm0MPzGTPH69",
	"jh2mzfNny5Sm0jxcptE8ugVaSP7f9cI11oBRigIt4X/MZlbqBR2E4LiOv6MVj9W610aOlYQ7LbUajKbM",
	"p+oWaJJxpIlYTzxge2McwhzvdKr2RkAltdn/V6aMYLLwBOmzc1+22bmrDmHFciBQGM2/7SKuPeoUVEfH",
	"vBJKrR1SJUwdnpo6e7TGgPRF2B46mPLBTO+tnD+Q3RFbNZP4UCYJYDNnrRrqV+tNe8klNNs9nI5UcACy",
	"WSEFBoR6Y46oSb37BG2vaclV4FF+Tq+tziE5KLTr8+1bZ+LnZec0jBQyIGgL8C/zvIPjM3cGj8o/2zIJ",
	"w3RO3/66+uoaX6Ki3NvKyBrygPmYnznFCtb9eG85rbTdJaDq8XEAxHAWLaLDTMQB453/3WAf2wkQ4x1u",
	"2H4Y8NrdynNmp1XUN5EhakZjN1zM06CXxneQj1XH6U60anIfmihe3wD75bd0jS+hQfd+OrLcKmBdgIej",
	"vSzUU+FU79eofqQrxzZcQS8lmHvbQfz+pai/dBrL/RvW6OX2Pnm++hraHqtr6siMuSt7hyjHuHN68tN3",
	"jrHJD/mcg1KTnCE9zD1oBmjdMb6MUN4/DQwUCsuyiYdJm2GZ50xtO01ON4y3GDqHnmpOzq0nd6fxeho3",
	"Y7vU/ULTtZYPqNcqBvPZKNoQyXkcZyJh2UYgzT/PPs/iaP+4/z8AAP//asKBZTQaAAA=",
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
