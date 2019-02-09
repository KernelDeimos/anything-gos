package main

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"

	"github.com/KernelDeimos/anything-gos/genner"
	"github.com/KernelDeimos/gottagofast/utilstr"
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

	// Note: in the future, perhaps instead of writing the function signature
	//       around every operation, there could be a template. There should
	//       also be a corresponding editor extension to display the function
	//       so the programmer doesn't need to look through different files.

	G.Interp.AddOperation("verify-args", func(
		args []interface{}) ([]interface{}, error) {

		result := []interface{}{}

		if len(args) < 1 {
			return nil, errors.New(
				"f_verify_args requires at least 1 argument")
		}

		if len(args)%2 == 0 {
			return nil, errors.New(
				"f_verify_args expects an odd number of values")
		}

		for _, arg := range args {
			if _, ok := arg.(string); !ok {
				return nil, errors.New("all arguments must be strings")
			}
		}

		// function name, for descriptive error messages
		fname := args[0]

		toCheck := [][]string{} // size: [n][2]

		for i := 1; i < len(args); i += 2 {
			varName := args[i].(string)
			varType := args[i+1].(string)
			toCheck = append(toCheck, []string{varName, varType})
		}

		printf(&result, "if len(args) < %d {", len(toCheck))
		printf(&result, "\treturn nil, "+
			`errors.New("%s requires at least %d arguments")`,
			fname, len(toCheck))
		printf(&result, "}\n")

		for _, item := range toCheck {
			var typName string
			switch item[1] {
			case "integer":
				typName = "int"
			default:
				typName = item[1]
			}
			result = append(result, fmt.Sprintf("var %s %s", item[0], typName))
		}

		result = append(result, "{")
		result = append(result, "\tvar ok bool")
		for i, item := range toCheck {
			vName := item[0]
			vType := item[1]

			switch vType {
			case "integer":
				result = append(result, "\tvar err error")
				printf(&result, "\tvar %sStr string", vName)
				printf(&result, "\t%sStr, ok = args[%d].(string)", vName, i)
				printf(&result, "\tif !ok {")
				printf(&result, "\t\treturn nil, "+
					`errors.New("%s: argument %d: %s; must be type %s")`,
					fname, i, vName, "int(string)")
				printf(&result, "\t}")
				printf(&result, "\t%s, err = strconv.Atoi(%sStr)", vName, vName)
				printf(&result, "\tif err != nil {")
				printf(&result, "\t\treturn nil, err")
				printf(&result, "\t}")
			default:
				printf(&result, "\t%s, ok = args[%d].(%s)", vName, i, vType)
				printf(&result, "\tif !ok {")
				printf(&result, "\t\treturn nil, "+
					`errors.New("%s: argument %d: %s; must be type %s")`,
					fname, i, vName, vType)
				printf(&result, "\t}")
			}
		}
		result = append(result, "}")

		return result, nil

	})

	G.Interp.AddOperation("ucwords", func(
		args []interface{}) ([]interface{}, error) {

		result := []interface{}{}

		for _, arg := range args {
			if str, ok := arg.(string); ok {
				result = append(result, strings.Title(str))
			} else {
				result = append(result, arg)
			}
		}

		return result, nil
	})

	G.Interp.AddOperation("repcsv", func(
		args []interface{}) ([]interface{}, error) {
		if len(args) < 1 {
			return nil, errors.New("repcsv requires at least 1 arguments")
		}

		var args0 string
		{
			var ok bool
			args0, ok = args[0].(string)
			if !ok {
				return nil, errors.New(
					"repcsv: argument 0: args0; must be type string")
			}
		}

		output := []interface{}{}

		for _, arg := range args[1:] {
			token, ok := arg.(string)
			if !ok {
				return nil, errors.New("setters: all tokens must be string")
			}
			parts := strings.Split(token, ",")

			replacements := map[string]string{}
			replacements["$$"] = "$"

			for i, repl := range parts {
				replacements["$"+strconv.Itoa(i+1)] = repl
				replacements["$uw-"+strconv.Itoa(i+1)] = strings.Title(repl)
				replacements["$u-"+strconv.Itoa(i+1)] = strings.ToUpper(repl)
				replacements["$l-"+strconv.Itoa(i+1)] = strings.ToLower(repl)

				if repl != "" {
					a := strings.Title(repl)
					b := strings.Split(a, " ")
					c := strings.Join(b, "")
					replacements["$ucc-"+strconv.Itoa(i+1)] = c
				}
				if repl != "" {
					a := repl[0:1] + strings.Title(repl[1:])
					b := strings.Split(a, " ")
					c := strings.Join(b, "")
					replacements["$lcc-"+strconv.Itoa(i+1)] = c
				}
				if repl != "" {
					a := repl[0:1] + strings.ToLower(repl[1:])
					b := strings.Split(a, " ")
					c := strings.Join(b, "_")
					replacements["$lcu-"+strconv.Itoa(i+1)] = c
				}
				if repl != "" {
					a := repl[0:1] + strings.ToLower(repl[1:])
					b := strings.Split(a, " ")
					c := strings.Join(b, "-")
					replacements["$lcd-"+strconv.Itoa(i+1)] = c
				}
			}

			result, _ := utilstr.AtomicReplace(args0, replacements)
			output = append(output, result)
		}

		return output, nil
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
