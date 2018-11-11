package genner

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/KernelDeimos/anything-gos/interp_a"
)

/*
	Note: the reason this file is so verbose is because I'm toying with ideas
	for the code generator as I write this. Eventually, I'd like most of the
	source code in this file to be generated with the very //::gen statements
	which it implements the functionality for.

	Also, the scopey nesty thing will hopefully make everything easier to read
*/

const (
	ErrAutoScopes = "automatic scope limiters have not been implemented yet"
	ErrEmptyGen   = "Nothing specified after ::gen"

	// Thinking of calling this thing "Anycode Preprocessor",
	// but for now it's "genner".
	TmpFilePrefix = "genner_"
)

type GenMode string

const (
	GenModeDefault  GenMode = "default"
	GenModeCollapse GenMode = "collapse"
)

type IndentationSet struct {
	runes []rune
}

func NewIndentationSet(runes []rune) IndentationSet {
	return IndentationSet{
		runes: runes,
	}
}

// TrimAndKeep trims left space as defined by the indentation set, and reports
// both the trimmed string and a string holding the indentation characters.
func (indents IndentationSet) TrimAndKeep(subject string) (
	modified, indent string,
) {
	cutset := string(indents.runes)
	modified = strings.TrimLeft(subject, cutset)
	indent = subject[0 : len(subject)-len(modified)]
	return
}

type TaskGenerator struct {
	// Line types
	IsGenLine *regexp.Regexp // "//::gen"
	IsEndLine *regexp.Regexp // "//::end"

	// Indentation characters
	Indents IndentationSet

	// Main parsing operation
	Exec interp_a.Operation
}

type Genner struct {
	Do     TaskGenerator
	Interp interp_a.HybridEvaluator
}

func NewDefaultGenner() Genner {
	tg := TaskGenerator{}
	tg.IsGenLine = regexp.MustCompile(`^\/\/::gen`)
	tg.IsEndLine = regexp.MustCompile(`^\/\/::end`)

	tg.Indents = IndentationSet{
		runes: []rune{' ', '\t'},
	}

	ii := interp_a.InterpreterFactoryA{}.MakeEmpty()
	tg.Exec = ii.OpEvaluate

	return Genner{
		Do:     tg,
		Interp: ii,
	}
}

func (genner TaskGenerator) GenerateUpdateFile(
	sourcefile string,
) error {
	var sourcePermissions os.FileMode

	// Get lines from input file
	var lines []string
	{
		var input string
		{
			// Open source file and save permisions
			sourceFileOS, err := os.Open(sourcefile)
			if err != nil {
				return err
			}
			stat, err := sourceFileOS.Stat()
			if err != nil {
				return err
			}
			sourcePermissions = stat.Mode()

			// Read entire input file
			inputBytes, err := ioutil.ReadAll(sourceFileOS)
			if err != nil {
				return err
			}

			// Write a temporary backup, since this file will be modified
			tmpf, err := ioutil.TempFile("", TmpFilePrefix)
			if err != nil {
				return errors.New("could not backup source file: " +
					err.Error())
			}

			// Log temporary file name in case someone needs to find it
			// NB: in case anyone tries to parse these log messages, beware
			//     of filenames containing single quotes in environments which
			//     allow this. Verify that the parsed name begins with genner_.
			log.Infof("Writing backup: '%s': '%s'", sourcefile, tmpf.Name())

			// Convert input file bytes to string
			input = string(inputBytes)
		}
		lines = strings.Split(input, "\n")
	}

	// Run the code generation function
	output, err := genner.ProcessLinesAndInsertGeneratedCode(lines)
	if err != nil {
		return err
	}

	// Convert output lines to byte slice
	var outputBytes []byte
	{
		outputBytes = []byte(strings.Join(output, "\n"))
	}

	// Write the modified file contents
	err = ioutil.WriteFile(sourcefile, outputBytes, sourcePermissions)
	if err != nil {
		return err
	}

	// Report no error; the file was modified successfully
	return nil
}

// ProcessLinesAndInsertGeneratedCode searches for comments which indicate
// some code generation operation to perform (which begin with //::gen under
// the default configuration), runs the operation Exec to produce lines of
// generated code, and inserts the generated code into a new list. Using the
// output []string in place of the input []string will effectively replace
// all ::gen-::end sections with newly generated code.
//
// If an error is reported, lines copied/generated before the error occurred
// will be returned. This is done to provide helpful debugging information.
func (genner TaskGenerator) ProcessLinesAndInsertGeneratedCode(
	lines []string) ([]string, error) {

	rIsGenLine := regexp.MustCompile(`^\/\/::gen`)

	output := []string{}

	const StateLookingForGen int = 1
	const StateLookingForEnd int = 2
	state := StateLookingForGen

	for _, line := range lines {
		// Determine what type of iteration
		switch state {
		case StateLookingForGen:
			goto LblLookForGen
		case StateLookingForEnd:
			goto LblLookForEnd
		}

	LblLookForGen:
		{
			// Always copy the line to output
			/*::idea lislang-go
			-	append output line
			*/
			output = append(output, line)

			// Get line without indentation, but keep indentation for later
			line, indentation := genner.Indents.TrimAndKeep(line)

			// Continue to next iteration unless line starts with "//::gen"
			isGenerate := rIsGenLine.MatchString(line)
			if !isGenerate {
				continue
			}

			// Get one or more tokens after "//::gen", or report error
			var parts []string
			{
				parts = strings.Split(line, " ")[1:]
				if len(parts) == 0 {
					log.Warn(ErrEmptyGen)
					return output, errors.New(ErrEmptyGen)
				}
			}

			// Get code to insert using the code generation operation
			genLines := []string{}
			{
				// Need to convert tokens from []string to []interface{}
				partsAsVoid := []interface{}{}
				for _, part := range parts {
					partsAsVoid = append(partsAsVoid, part)
				}

				// Run the code generator, get results as []interface{}
				results, err := genner.Exec(partsAsVoid)
				if err != nil {
					return output, err
				}

				// Interpret results as lines of code to insert
				for _, datum := range results {
					switch line := datum.(type) {
					case string:
						// A single line of generated code
						genLines = append(genLines, line)
					case []string:
						// TODO: implement automatic scope limiters here
						return output, errors.New(ErrAutoScopes)
					}
				}
			}

			for _, genLine := range genLines {
				output = append(output, indentation+genLine)
			}

			state = StateLookingForEnd

		}
		continue // end LblLookForGen

	LblLookForEnd:
		{
			line, indentation := genner.Indents.TrimAndKeep(line)
			if genner.IsEndLine.MatchString(line) {
				// Add end indicator to output to keep it in the file
				output = append(output, indentation+line)
				state = StateLookingForGen
			}
		}
		continue // end lblLookForEnd

	}

	// Report final output with no error :)
	return output, nil
}
