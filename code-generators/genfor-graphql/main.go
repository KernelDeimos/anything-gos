package main

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/KernelDeimos/anything-gos/genner"
)

func append_str(inI *[]interface{}, inS ...string) {
	for _, str := range inS {
		*inI = append(*inI, str)
	}
}

func printf(inI *[]interface{}, format string, args ...interface{}) {
	*inI = append(*inI, fmt.Sprintf(format, args...))
}

func main() {
	files := os.Args[1:]

	G := genner.NewGennerWithBuiltinsTest()

	G.Interp.AddOperation("graphql-fields", func(
		args []interface{}) ([]interface{}, error) {

		result := []interface{}{}

		for _, arg := range args {
			tuple, ok := arg.([]interface{})
			if !ok {
				return result, errors.New("expected tuple")
			}
			if len(tuple) < 2 {
				return result, errors.New("incorrect length of tuple")
			}

			var identifier, typ string
			{
				var ok bool
				identifier, ok = tuple[0].(string)
				if !ok {
					return result, errors.New("identifier must be string")
				}
				typ, ok = tuple[0].(string)
				if !ok {
					return result, errors.New("identifier must be string")
				}
			}

			fieldType := "String"
			switch typ {
			case "bool":
				fieldType = "Boolean"
			}

			printf(&result, `"%s": &graphql.Field{`, identifier)
			printf(&result, "\t"+`Type: graphql.%s,`, fieldType)
			printf(&result, `},`)
		}

		return result, nil
	})

	if len(files) == 0 {
		logrus.Warn("No files specified; doing nothing")
		return
	}

	for _, file := range files {
		logrus.Infof("=== %s ===", file)
		err := G.Do.GenerateUpdateFile(file)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
