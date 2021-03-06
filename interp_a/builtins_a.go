package interp_a

//go:generate genfor-interp-a $GOFILE

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
)

func BuiltinFormat(args []interface{}) ([]interface{}, error) {
	//::gen verify-args format args0 string
	if len(args) < 1 {
		return nil, errors.New("format requires at least 1 arguments")
	}

	var args0 string
	{
		var ok bool
		args0, ok = args[0].(string)
		if !ok {
			return nil, errors.New("format: argument 0: args0; must be type string")
		}
	}
	//::end
	result := fmt.Sprintf(args0, args[1:]...)

	return []interface{}{result}, nil
}

func BuiltinStrings(args []interface{}) ([]interface{}, error) {
	stringArgs := []string{}
	for _, arg := range args {
		stringArgs = append(stringArgs, fmt.Sprint(arg))
	}

	return []interface{}{stringArgs}, nil
}

func BuiltinPassthrough(args []interface{}) ([]interface{}, error) {
	return args, nil
}

func BuiltinTie(args []interface{}) ([]interface{}, error) {
	return []interface{}{args}, nil
}

func BuiltinCat(args []interface{}) ([]interface{}, error) {
	result := ""
	for _, arg := range args {
		result = result + fmt.Sprint(arg)
	}
	return []interface{}{result}, nil
}

func BuiltinJoinLF(args []interface{}) ([]interface{}, error) {
	strs := []string{}
	for _, arg := range args {
		strs = append(strs, fmt.Sprint(arg))
	}
	return []interface{}{strings.Join(strs, "\n")}, nil
}

func BuiltinCatFromList(args []interface{}) ([]interface{}, error) {
	result := ""
	for _, arg := range args {
		l, ok := arg.([]interface{})
		if !ok {
			return nil, errors.New("unexpected non-list in cat[]")
		}
		for _, part := range l {
			result = result + fmt.Sprint(part)
		}
	}
	return []interface{}{result}, nil
}

func BuiltinCatRepeat(args []interface{}) ([]interface{}, error) {
	//::gen verify-args ditto times int
	if len(args) < 1 {
		return nil, errors.New("ditto requires at least 1 arguments")
	}

	var times int
	{
		var ok bool
		times, ok = args[0].(int)
		if !ok {
			return nil, errors.New("ditto: argument 0: times; must be type int")
		}
	}
	//::end
	result := ""
	for i := 0; i < times; i++ {
		for _, arg := range args[1:] {
			result = result + fmt.Sprint(arg)
		}
	}
	return []interface{}{result}, nil
}

func BuiltinStore(args []interface{}) ([]interface{}, error) {
	op := func(_ []interface{}) ([]interface{}, error) {
		return args, nil
	}
	return []interface{}{Operation(op)}, nil
}

func BuiltinCodeCallsData(args []interface{}) ([]interface{}, error) {
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

	return []interface{}{contains(args)}, nil
}

func BuiltinUnfile(args []interface{}) ([]interface{}, error) {
	//::gen verify-args unfile args0 string
	if len(args) < 1 {
		return nil, errors.New("unfile requires at least 1 arguments")
	}

	var args0 string
	{
		var ok bool
		args0, ok = args[0].(string)
		if !ok {
			return nil, errors.New("unfile: argument 0: args0; must be type string")
		}
	}
	//::end
	data, err := ioutil.ReadFile(args0)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	// TODO: there ought to be a gen function for this
	linesI := []interface{}{}
	for _, line := range lines {
		linesI = append(linesI, line)
	}
	return linesI, nil
}

func BuiltinSlurp(args []interface{}) ([]interface{}, error) {
	//::gen verify-args slurp args0 string
	if len(args) < 1 {
		return nil, errors.New("slurp requires at least 1 arguments")
	}

	var args0 string
	{
		var ok bool
		args0, ok = args[0].(string)
		if !ok {
			return nil, errors.New("slurp: argument 0: args0; must be type string")
		}
	}
	//::end
	data, err := ioutil.ReadFile(args0)
	if err != nil {
		return nil, err
	}
	return []interface{}{string(data)}, nil
}

func BuiltinJsonEncodeOne(args []interface{}) ([]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("json-encode-one requires at least 1 arguments")
	}

	value := args[0]

	bytes, err := json.Marshal(value)
	if err != nil {
		return []interface{}{"error"}, err
	}
	return []interface{}{string(bytes)}, nil
}

func BuiltinJsonDecodeOne(args []interface{}) ([]interface{}, error) {
	//::gen verify-args json-decode-one jsonText string
	if len(args) < 1 {
		return nil, errors.New("json-decode-one requires at least 1 arguments")
	}

	var jsonText string
	{
		var ok bool
		jsonText, ok = args[0].(string)
		if !ok {
			return nil, errors.New("json-decode-one: argument 0: jsonText; must be type string")
		}
	}
	//::end

	var result interface{}

	err := json.Unmarshal([]byte(jsonText), &result)
	if err != nil {
		return []interface{}{"error"}, err
	}
	return []interface{}{result}, nil
}

func BuiltinFnTemplate(args []interface{}) ([]interface{}, error) {
	strArgs := []string{}
	for _, value := range args {
		strArgs = append(strArgs, fmt.Sprint(value))
	}
	templateText := strings.Join(strArgs, "\n")

	fmt.Println(">>>" + templateText + "<<<")

	t, err := template.New("").Parse(templateText)
	if err != nil {
		return nil, err
	}

	fn := func(args []interface{}) ([]interface{}, error) {
		buf := bytes.NewBufferString("")
		err := t.Execute(buf, args)
		interfaceLines := []interface{}{}
		strLines := strings.Split(buf.String(), "\n")
		for _, value := range strLines {
			fmt.Println(">>" + value + "<<")
			interfaceLines = append(interfaceLines, value)
		}
		return interfaceLines, err
	}

	return []interface{}{Operation(fn)}, nil
}
