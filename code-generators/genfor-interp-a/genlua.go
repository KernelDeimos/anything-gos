package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func GenerateVerifyLuaArgs(
	fname string, toCheck [][]string,
	result *[]interface{},
) {

	printf(result, "if L.GetTop() < %d {", len(toCheck))
	printf(result, `L.RaiseError("%s requires at least %d arguments")`,
		fname, len(toCheck),
	)
	printf(result, `return 0`)
	printf(result, "}\n")

	for _, item := range toCheck {
		typName := item[1]
		*result = append(*result,
			fmt.Sprintf("var %s %s", item[0], typName))
	}

	*result = append(*result, "{")
	for i, item := range toCheck {
		luaIndex := i + 1

		vName := item[0]
		vType := item[1]

		switch vType {
		case "int":
			printf(result, "\t%s = L.ToInt(%d)", vName, luaIndex)
		case "string":
			printf(result, "\t%s = L.ToString(%d)", vName, luaIndex)
		default:
			intermediateName := "tmpFor" + strings.Title(vName)
			printf(result, "\tvar ok bool")
			printf(result, "\t%s := L.ToUserData(%d)",
				intermediateName, luaIndex)

			printf(result, "\t%s, ok = %s.Value.(%s)",
				vName,
				intermediateName,
				vType,
			)
			printf(result, "\tif !ok {")
			printf(result,
				`L.RaiseError("%s: argument %d: %s; must be type %s")`,
				fname, i, vName, vType,
			)
			printf(result, `return 0`)
			printf(result, "\t}")
		}
	}
	*result = append(*result, "}")
}

func GenerateLuaBinding(
	iname string, // name of interpreter object
	fname string, // name of function in interpreter
	call string, // Golang function to call
	toCheck [][]string,
	toReturn [][]string,
	result *[]interface{},
) {
	printf(result, iname+`.SetGlobal("`+
		fname+ // TODO: escape quotes
		`", `+iname+`.NewFunction(func(`+"\n\t"+
		`L *lua.LState) int {`+"\n",
	)
	defer printf(result, "}))")

	// Gets populated with the variadic argument's name if it exists
	var variadicArg string

	if len(toCheck) > 0 {
		// Check for variadic argument
		lastItem := toCheck[len(toCheck)-1]
		if strings.HasPrefix(lastItem[1], "...") {
			// Temporary check: warn that only ...interface{} works right now
			if lastItem[1] != "...interface{}" {
				logrus.Warn(
					"Variadic doesn't work in LUA yet.",
				)
			}
			toCheck = toCheck[:len(toCheck)-1]
			variadicArg = lastItem[0]
		}

	}

	if len(toCheck) > 0 {

		GenerateVerifyLuaArgs(fname, toCheck, result)

		if variadicArg != "" {
			printf(result, "%s := []interface{}{}", variadicArg)
			printf(result, "for i := %d; i <= L.GetTop(); i++ {",
				len(toCheck)+1)
			printf(result, "\tlv := L.Get(i)")
			printf(result, "\tswitch lv.Type() {")
			printf(result, "\tcase lua.LTString:")
			printf(result, "\t\t%s = append(%s, L.ToString(i))",
				variadicArg, variadicArg)
			printf(result, "\tcase lua.LTNumber:")
			printf(result, "\t\t%s = append(%s, L.ToNumber(i))",
				variadicArg, variadicArg)
			printf(result, "\t}")
			printf(result, "}")
		}

		// Check for method receiver
		// TODO: this is dupe code (genbindings.go)
		{
			callPath := strings.Split(call, ".")
			if len(callPath) > 1 {
				if callPath[0] == toCheck[0][0] {
					toCheck = toCheck[1:]
				}
			}
		}
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

		if variadicArg != "" {
			argumentNames = append(argumentNames, variadicArg+"...")
		}

		// Generate function call and assignment
		if len(returnNames) > 0 {
			printf(result, "%s := %s(%s)",
				strings.Join(returnNames, ","),
				call,
				strings.Join(argumentNames, ","),
			)
		} else {
			printf(result, "%s(%s)",
				call,
				strings.Join(argumentNames, ","),
			)
		}
	}

	// See if last return value is of type error
	if lastReturnIsTypeError {
		lastReturnValueName := returnNames[len(returnNames)-1]
		printf(result, "if %s != nil {", lastReturnValueName)
		printf(result, "\tL.")
		printf(result, "\treturn 0", lastReturnValueName)
		printf(result, "}")
	}

	if lastReturnIsTypeError {
		returnNames = returnNames[:len(returnNames)-1]
		toReturn = toReturn[:len(toReturn)-1]
	}

	for i := len(toReturn) - 1; i >= 0; i-- {
		returnName := toReturn[i][0]
		returnType := toReturn[i][1]

		// GOTCHA: Do not use `continue` without closing this
		printf(result, "{")

		switch returnType {
		case "int":
			fallthrough
		case "float64":
			printf(result, "L.Push(lua.LNumber(%s))", returnName)
		case "string":
			printf(result, "L.Push(lua.LString(%s))", returnName)
		default:
			printf(result, "tmp := L.NewUserData()")
			printf(result, "tmp.Value = %s", returnName)
			printf(result, "L.Push(tmp)")
		}

		printf(result, "}")
	}

	// Generate return statement
	printf(result, "return %d", len(toReturn))

}
