package internal

import (
	"fmt"
	"net/url"
	"path"
	"test"

	"test/logger"

	"github.com/xeipuuv/gojsonschema"
)

func Reference() error {
	p1 := path.Join(test.RootDir(), "jsonschema/internal/schema.json")
	p2 := path.Join(test.RootDir(), "jsonschema/internal/document.json")
	return reference(p1, p2)
}

func uri(s string) string {
	u, err := url.Parse("file:///")
	if err != nil {
		logger.Panic(err)
	}
	u.Path = path.Join(u.Path, s)
	return u.String()
}

func reference(want, got string) error {

	schemaLoader := gojsonschema.NewReferenceLoader(uri(want))
	documentLoader := gojsonschema.NewReferenceLoader(uri(got))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		logger.Error(err)
		return err
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	return nil
}

func Reference1() error {
	p1 := path.Join(test.RootDir(), "jsonschema/internal/s1.json")
	p2 := path.Join(test.RootDir(), "jsonschema/internal/input1.json")
	return reference(p1, p2)
}
