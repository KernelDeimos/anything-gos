package genner

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/KernelDeimos/anything-gos/interp_a"
	"github.com/KernelDeimos/gottagofast/toolparse"
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

type GennerEnvironment interface {
	SetInputBuffer([]string)
}

type TaskGenerator struct {
	// Line types
	IsRunLine *regexp.Regexp // "//::declare"
	IsGenLine *regexp.Regexp // "//::gen"
	IsEndLine *regexp.Regexp // "//::end"

	// Indentation characters
	Indents IndentationSet

	// Main parsing operation
	Exec interp_a.Operation

	// API to environment
	Env GennerEnvironment
}

// Genner is a wrapper for the task generator which also provides access to the
// interpreter.
type Genner struct {
	Do     TaskGenerator
	Interp interp_a.HybridEvaluator
}

func constructGenner(ii interp_a.HybridEvaluator) Genner {
	env := &GennerEnvironmentDefault{}
	ii.AddOperation("DATA", env.OpReadInput)

	tg := TaskGenerator{}
	tg.IsRunLine = regexp.MustCompile(`^\/\/::run`)
	tg.IsGenLine = regexp.MustCompile(`^\/\/::gen`)
	tg.IsEndLine = regexp.MustCompile(`^\/\/::end`)

	tg.Indents = IndentationSet{
		runes: []rune{' ', '\t'},
	}

	tg.Exec = ii.OpEvaluate
	tg.Env = env

	return Genner{
		Do:     tg,
		Interp: ii,
	}
}

func NewGennerHackyRegexCat(ii interp_a.HybridEvaluator,
	prefix string) Genner {
	env := &GennerEnvironmentDefault{}
	ii.AddOperation("DATA", env.OpReadInput)

	tg := TaskGenerator{}
	tg.IsRunLine = regexp.MustCompile(`^` + prefix + `::run`)
	tg.IsGenLine = regexp.MustCompile(`^` + prefix + `::gen`)
	tg.IsEndLine = regexp.MustCompile(`^` + prefix + `::end`)

	tg.Indents = IndentationSet{
		runes: []rune{' ', '\t'},
	}

	tg.Exec = ii.OpEvaluate
	tg.Env = env

	return Genner{
		Do:     tg,
		Interp: ii,
	}
}

func NewDefaultGenner() Genner {
	ii := interp_a.InterpreterFactoryA{}.MakeEmpty()
	return constructGenner(ii)
}

func NewGennerWithBuiltinsTest() Genner {
	ii := interp_a.InterpreterFactoryA{}.MakeExec()
	return constructGenner(ii)
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
	log.Infof("Rewriting code: '%s' with mode %#o", sourcefile, sourcePermissions)
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

	// Code to output
	output := []string{}

	// Buffer for copying lines of code (i.e. templates)
	buffer := []string{}
	var execLine string // command to run after copying code

	const StateLookingForGen int = 1
	const StateLookingForEnd int = 2
	const StateCopyDataBlock int = 3
	state := StateLookingForGen

	for _, line := range lines {
		// Determine what type of iteration
		switch state {
		case StateLookingForGen:
			goto LblLookForGen
		case StateLookingForEnd:
			fallthrough
		case StateCopyDataBlock:
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
			locRun := genner.IsRunLine.FindStringIndex(line)
			locGen := genner.IsGenLine.FindStringIndex(line)
			if locGen != nil {
				log.Info("gen location found")
				// Get one or more tokens after "//::gen", or report error
				genLine := line[locGen[1]:]
				parts, err := toolparse.ParseListSimple(genLine)
				if err != nil {
					return output, err
				}

				// Get code to insert using the code generation operation
				genLines := []string{}
				{
					// Run the code generator, get results as []interface{}
					results, err := genner.Exec(parts)
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
			} else if locRun != nil {
				log.Info("run location found")
				buffer = []string{}
				execLine = line[locRun[1]:]
				state = StateCopyDataBlock
				continue
			}

		}
		continue // end LblLookForGen

	LblLookForEnd:
		{
			line, indentation := genner.Indents.TrimAndKeep(line)
			if genner.IsEndLine.MatchString(line) {
				log.Info("ending")
				// If state was CopyDataBlock, execute the command now
				if state == StateCopyDataBlock {
					// Update buffer in interpreter environment
					genner.Env.SetInputBuffer(buffer)
					log.Info("here")
					parts, err := toolparse.ParseListSimple(execLine)
					if err != nil {
						return output, err
					}

					// Run the code generator
					results, err := genner.Exec(parts)

					// Print results to screen
					// (results do not go in the source file, because mixing
					//  output text with input text is probably bad         )
					log.Warn(results)

					if err != nil {
						return output, err
					}
				}
				// Add end indicator to output to keep it in the file
				output = append(output, indentation+line)
				state = StateLookingForGen
			} else if state == StateCopyDataBlock {
				// This else block means we're in the mode where we copy lines
				// of data (i.e. block under "::run" comment)
				buffer = append(buffer, line)

				// Also we should keep these lines in the file
				output = append(output, indentation+line)
			}
		}
		continue // end lblLookForEnd

	}

	// Report final output with no error :)
	return output, nil
}
