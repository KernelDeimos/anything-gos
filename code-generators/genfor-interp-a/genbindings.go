package main

import (
	"fmt"
	"strings"
)

func GenerateVerifyArgs(
	fname string, toCheck [][]string,
	result *[]interface{},
) {

	printf(result, "if len(args) < %d {", len(toCheck))
	printf(result, "\treturn nil, "+
		`errors.New("%s requires at least %d arguments")`,
		fname, len(toCheck))
	printf(result, "}\n")

	for _, item := range toCheck {
		var typName string
		switch item[1] {
		case "integer":
			typName = "int"
		default:
			typName = item[1]
		}
		*result = append(*result,
			fmt.Sprintf("var %s %s", item[0], typName))
	}

	*result = append(*result, "{")
	*result = append(*result, "\tvar ok bool")
	for i, item := range toCheck {
		vName := item[0]
		vType := item[1]

		switch vType {
		case "integer":
			*result = append(*result, "\tvar err error")
			printf(result, "\tvar %sStr string", vName)
			printf(result, "\t%sStr, ok = args[%d].(string)", vName, i)
			printf(result, "\tif !ok {")
			printf(result, "\t\treturn nil, "+
				`errors.New("%s: argument %d: %s; must be type %s")`,
				fname, i, vName, "int(string)")
			printf(result, "\t}")
			printf(result, "\t%s, err = strconv.Atoi(%sStr)", vName, vName)
			printf(result, "\tif err != nil {")
			printf(result, "\t\treturn nil, err")
			printf(result, "\t}")
		default:
			printf(result, "\t%s, ok = args[%d].(%s)", vName, i, vType)
			printf(result, "\tif !ok {")
			printf(result, "\t\treturn nil, "+
				`errors.New("%s: argument %d: %s; must be type %s")`,
				fname, i, vName, vType)
			printf(result, "\t}")
		}
	}
	*result = append(*result, "}")
}

func GenerateBinding(
	iname string, // name of interpreter object
	fname string, // name of function in interpreter
	call string, // Golang function to call
	toCheck [][]string,
	toReturn [][]string,
	result *[]interface{},
) {
	printf(result, iname+`.AddOperation("`+
		fname+ // TODO: escape quotes
		`", func(`+"\n\t"+
		`args []interface{}) ([]interface{}, error) {`+"\n",
	)
	defer printf(result, "})")

	// TODO: generate binding thingy
	if len(toCheck) > 0 {
		GenerateVerifyArgs(fname, toCheck, result)
	}

	// List of names of return values, used for:
	// - left-hand side of assignment where function is called
	// - values inside []interface{} that is returned
	returnNames := []string{}
	for _, item := range toReturn {
		name := item[0]
		returnNames = append(returnNames, name)
	}

	var lastReturnIsTypeError = len(toReturn) > 0 &&
		toReturn[len(toReturn)-1][1] == "error"

	{
		// List of names of arguments, used for function call
		argumentNames := []string{}
		for _, item := range toCheck {
			name := item[0]
			argumentNames = append(argumentNames, name)
		}

		// Generate function call and assignment
		printf(result, "%s := %s(%s)",
			strings.Join(returnNames, ","),
			call,
			strings.Join(argumentNames, ","),
		)
	}

	// See if last return value is of type error
	if lastReturnIsTypeError {
		lastReturnValueName := returnNames[len(returnNames)-1]
		printf(result, "if %s != nil {", lastReturnValueName)
		printf(result, "\treturn nil, %s", lastReturnValueName)
		printf(result, "}")
	}

	if lastReturnIsTypeError {
		returnNames = returnNames[:len(returnNames)-1]
	}

	// Generate return statement
	printf(result, "return []interface{}{%s}, nil",
		strings.Join(returnNames, ","),
	)

}
