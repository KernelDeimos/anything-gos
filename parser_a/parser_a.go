package parser_a

import (
	// "bufio"
	// "fmt"
	// "github.com/sirupsen/logrus"
	// "os"

	"github.com/KernelDeimos/anything-gos/interp_a"
	// "github.com/KernelDeimos/gottagofast/toolparse"
)

type TokenList []interface{}

type Context struct {
	tokens TokenList
}

type ParserState struct {
	ContextStack []Context
}

func hasDataFunction(args []interface{}) bool {
	var contains func([]interface{}) bool
	contains = func(args []interface{}) bool {
		if len(args) < 1 {
			return false
		}
		if args[0] == "data" {
			return true
		}
		for _, arg := range args[1:] {
			switch a := arg.(type) {
			case []interface{}:
				if contains(a) {
					return true
				}
			}
		}
		return false
	}
	return contains(args)
}

func parse(lines []string) {
	for _, line := range lines {

		// hasData := hasDataFunction()
		// ii := interp_a.InterpreterFactoryA{}.MakeEmpty()
		// ii.AddOperation("data")
	}
}
