package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"

	"github.com/KernelDeimos/anything-gos/interp_a"
	"github.com/KernelDeimos/gottagofast/toolparse"
)

const DemoMsg = `Hello! Welcome to the interpreter_a demo!

Try some of the following:
	cat a b c
	format 'Hello %s!' World
	exec cat a b c
	cat (format 'Hello %s' there) !
	sayhello

Use EOF (Ctrl+D) to exit. Use rlwrap for best experience.
`

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	ii := interp_a.InterpreterFactoryA{}.MakeExec()

	// Defining a custom function: sayhello
	ii.AddOperation("sayhello", func(args []interface{}) ([]interface{}, error) {
		return []interface{}{"hello"}, nil
	})

	// Defining default behaviour (no function found)
	/*
		ii.SetDefaultBehaviour(func(args []interface{}) ([]interface{}, error) {
			return []interface{}{"not found"}, nil
		}, interp_a.EntryIsEvaluator)
	*/

	// Alias the main executing function; allows you to prepend "exec" to any
	// line of code with no effect
	_, err := ii.OpAddOperation([]interface{}{
		"exec", interp_a.Operation(ii.OpEvaluate)})
	// ^ could also use AddOperation instead of OpAddOperation like above, and
	//   then not need to wrap in []interface{}
	// ^ OR, could use ii.OpEvaluate([]interface{}{
	//                           ":","exec",interp_a.Operation(ii.OpEvaluate)})
	if err != nil {
		logrus.Fatal(err)
	}

	if len(os.Args) > 1 {
		f := os.Args[1]
		fBytes, err := ioutil.ReadFile(f)
		if err != nil {
			logrus.Fatal(err)
		}
		fStr := string(fBytes)
		fStr = strings.Replace(fStr, "\n", "", -1)
		fStr = fStr
		list, err := toolparse.ParseListSimple(fStr)
		logrus.Debug(list)
		if err != nil {
			logrus.Error(err)
		}
		output, err := ii.OpEvaluate(list)
		if err != nil {
			logrus.Error(err)
		}
		fmt.Println(output)
		return
	}

	fmt.Println(DemoMsg)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			logrus.Fatal(err)
		}
		list, err := toolparse.ParseListSimple(text)
		if err != nil {
			logrus.Error(err)
		}
		output, err := ii.OpEvaluate(list)
		if err != nil {
			logrus.Error(err)
		}
		fmt.Println(output)
	}
}
