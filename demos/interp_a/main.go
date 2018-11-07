package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

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
	ii := interp_a.InterpreterFactoryA{}.MakeExec()

	// Defining a custom function: sayhello
	ii.AddOperation("sayhello", func(args []interface{}) ([]interface{}, error) {
		return []interface{}{"hello"}, nil
	})

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
