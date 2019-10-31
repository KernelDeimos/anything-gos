package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/rosewoodmedia/gofaast/faast"
	"github.com/rosewoodmedia/gofaast/faastfmt"
	"github.com/sirupsen/logrus"

	"github.com/KernelDeimos/anything-gos/genner"
	"github.com/KernelDeimos/anything-gos/interp_a"
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
	logrus.Info("program genfor-interp-a v1.1.0")
	files := os.Args[1:]

	G := genner.NewGennerWithBuiltinsTest()

	if len(files) >= 2 && files[0] == "--hackprefix" {
		logrus.Warn("Interpreting next argument as comment prefx")
		ii := interp_a.InterpreterFactoryA{}.MakeExec()

		G = genner.NewGennerHackyRegexCat(ii, files[1])
		files = files[2:]
	}

	// Note: in the future, perhaps instead of writing the function signature
	//       around every operation, there could be a template. There should
	//       also be a corresponding editor extension to display the function
	//       so the programmer doesn't need to look through different files.

	G.Interp.AddOperation("gen-binding", func(
		args []interface{}) ([]interface{}, error) {

		result := []interface{}{}

		//::gen verify-args gen-binding iname string fname string call string args2 []string returns []string
		if len(args) < 5 {
			return nil, errors.New("gen-binding requires at least 5 arguments")
		}

		var iname string
		var fname string
		var call string
		var args2 []string
		var returns []string
		{
			var ok bool
			iname, ok = args[0].(string)
			if !ok {
				return nil, errors.New("gen-binding: argument 0: iname; must be type string")
			}
			fname, ok = args[1].(string)
			if !ok {
				return nil, errors.New("gen-binding: argument 1: fname; must be type string")
			}
			call, ok = args[2].(string)
			if !ok {
				return nil, errors.New("gen-binding: argument 2: call; must be type string")
			}
			args2, ok = args[3].([]string)
			if !ok {
				return nil, errors.New("gen-binding: argument 3: args2; must be type []string")
			}
			returns, ok = args[4].([]string)
			if !ok {
				return nil, errors.New("gen-binding: argument 4: returns; must be type []string")
			}
		}
		//::end

		toCheck := [][]string{}  // size: [n][2]
		toReturn := [][]string{} // size: [n][2]

		for i := 0; i < len(args2); i += 2 {
			varName := args2[i]
			varType := args2[i+1]
			toCheck = append(toCheck, []string{varName, varType})
		}

		for i := 0; i < len(returns); i += 2 {
			varName := returns[i]
			varType := returns[i+1]
			toReturn = append(toReturn, []string{varName, varType})
		}

		GenerateBinding(
			iname, fname, call, toCheck, toReturn, &result)

		logrus.Debug("it's okay; everything will be fine...")

		return result, nil
	})

	G.Interp.AddOperation("bind-file", func(
		args []interface{}) ([]interface{}, error) {

		result := []interface{}{}

		//::gen verify-args gen-bindings-from-file iname string filename string
		if len(args) < 2 {
			return nil, errors.New("gen-bindings-from-file requires at least 2 arguments")
		}

		var iname string
		var filename string
		{
			var ok bool
			iname, ok = args[0].(string)
			if !ok {
				return nil, errors.New("gen-bindings-from-file: argument 0: iname; must be type string")
			}
			filename, ok = args[1].(string)
			if !ok {
				return nil, errors.New("gen-bindings-from-file: argument 1: filename; must be type string")
			}
		}
		//::end

		// Create the AST by parsing src.
		fset := token.NewFileSet() // positions are relative to fset
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		f, err := parser.ParseFile(fset, filename, string(data), 0)
		if err != nil {
			panic(err)
		}

		faastp := faastfmt.SimplePrinter{}

		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			toCheck := [][]string{}  // size: [n][2]
			toReturn := [][]string{} // size: [n][2]

			call := fn.Name.Name
			// TODO: code generator to create enumerated types for
			//       every available function in strcase, then use
			//       enum type as parameter to define this behaviour.
			fname := strcase.ToKebab(call)

			if fn.Type.Params != nil {
				params, _ := faast.FieldListToFaast(fn.Type.Params.List)
				for _, param := range params {
					toCheck = append(toCheck, []string{
						param.Name.Name,
						faastp.PrintType(param.Type),
					})
				}
			}

			if fn.Type.Results != nil {
				results, _ := faast.FieldListToFaast(fn.Type.Results.List)
				for _, result := range results {
					toReturn = append(toReturn, []string{
						result.Name.Name,
						faastp.PrintType(result.Type),
					})
				}
			}

			GenerateBinding(
				iname, fname, call, toCheck, toReturn, &result)
		}

		return result, nil
	})

	G.Interp.AddOperation("lua-binding", func(
		args []interface{}) ([]interface{}, error) {

		result := []interface{}{}

		//::gen verify-args lua-binding iname string fname string call string args2 []string returns []string
		if len(args) < 5 {
			return nil, errors.New("lua-binding requires at least 5 arguments")
		}

		var iname string
		var fname string
		var call string
		var args2 []string
		var returns []string
		{
			var ok bool
			iname, ok = args[0].(string)
			if !ok {
				return nil, errors.New("lua-binding: argument 0: iname; must be type string")
			}
			fname, ok = args[1].(string)
			if !ok {
				return nil, errors.New("lua-binding: argument 1: fname; must be type string")
			}
			call, ok = args[2].(string)
			if !ok {
				return nil, errors.New("lua-binding: argument 2: call; must be type string")
			}
			args2, ok = args[3].([]string)
			if !ok {
				return nil, errors.New("lua-binding: argument 3: args2; must be type []string")
			}
			returns, ok = args[4].([]string)
			if !ok {
				return nil, errors.New("lua-binding: argument 4: returns; must be type []string")
			}
		}
		//::end

		toCheck := [][]string{}  // size: [n][2]
		toReturn := [][]string{} // size: [n][2]

		for i := 0; i < len(args2); i += 2 {
			varName := args2[i]
			varType := args2[i+1]
			toCheck = append(toCheck, []string{varName, varType})
		}

		for i := 0; i < len(returns); i += 2 {
			varName := returns[i]
			varType := returns[i+1]
			toReturn = append(toReturn, []string{varName, varType})
		}

		GenerateLuaBinding(
			iname, fname, call, toCheck, toReturn, &result)

		logrus.Debug("it's okay; everything will be fine...")

		return result, nil
	})

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
		fname := args[0].(string)

		toCheck := [][]string{} // size: [n][2]

		for i := 1; i < len(args); i += 2 {
			varName := args[i].(string)
			varType := args[i+1].(string)
			toCheck = append(toCheck, []string{varName, varType})
		}

		GenerateVerifyArgs(fname, toCheck, &result)

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
					b := strings.Split(repl, " ")
					lastPartList := b[1:]
					lastPartStr := strings.Title(strings.Join(lastPartList, " "))
					lastPartList = strings.Split(lastPartStr, " ")
					lastPartStr = strings.Join(lastPartList, "")
					replacements["$lcc-"+strconv.Itoa(i+1)] = b[0] + lastPartStr
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

	G.Interp.AddOperation("repeat", func(
		args []interface{}) ([]interface{}, error) {
		if len(args) < 1 {
			return nil, errors.New("repeat requires at least 1 arguments")
		}

		var args0 int
		{
			var ok bool
			args0, ok = args[0].(int)
			if !ok {
				return nil, errors.New(
					"repeat: argument 0: args0; must be type int")
			}
		}

		var args1 string
		{
			var ok bool
			args1, ok = args[1].(string)
			if !ok {
				return nil, errors.New(
					"repeat: argument 1: args1; must be type string")
			}
		}

		output := []interface{}{}

		for i := 0; i < args0; i++ {
			output = append(output, args1)
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
