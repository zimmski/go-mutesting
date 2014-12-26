package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	//"github.com/zimmski/go-mutesting"
	"github.com/zimmski/go-mutesting/importing"
	"github.com/zimmski/go-mutesting/mutator"
	_ "github.com/zimmski/go-mutesting/mutator/branch"
)

const (
	returnOk = iota
	returnHelp
	returnBashCompletion
	returnError
)

var opts struct {
	General struct {
		Help    bool `long:"help" description:"Show this help message"`
		Verbose bool `long:"verbose" description:"Verbose log output"`
	} `group:"General options"`

	Mutator struct {
		DisableMutators []string `long:"disable" description:"Disable mutator or mutators using * as a suffix pattern"`
		ListMutators    bool     `long:"list-mutators" description:"List all available mutators"`
	} `group:"Mutator options"`

	Remaining struct {
		Targets []string `description:"Packages, directories and files even with patterns"`
	} `positional-args:"true" required:"true"`
}

func checkArguments() {
	p := flags.NewNamedParser("go-mutesting", flags.None)

	p.ShortDescription = "Mutation testing for Go source code"

	if _, err := p.AddGroup("go-mutesting", "go-mutesting arguments", &opts); err != nil {
		exitError(err.Error())
	}

	completion := len(os.Getenv("GO_FLAGS_COMPLETION")) > 0

	_, err := p.Parse()
	if (opts.General.Help || len(os.Args) == 1) && !completion {
		p.WriteHelp(os.Stdout)

		os.Exit(returnHelp)
	} else if opts.Mutator.ListMutators {
		for _, name := range mutator.List() {
			fmt.Println(name)
		}

		os.Exit(returnOk)
	}

	if err != nil {
		exitError(err.Error())
	}

	if completion {
		os.Exit(returnBashCompletion)
	}
}

func exitError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)

	os.Exit(returnError)
}

func main() {
	checkArguments()

	files := importing.FilesOfArgs(opts.Remaining.Targets)
	if len(files) == 0 {
		exitError("Could not find any suitable Go source files")
	}

	var mutators []mutator.Mutator

MUTATOR:
	for _, name := range mutator.List() {
		if len(opts.Mutator.DisableMutators) != 0 {
			for _, d := range opts.Mutator.DisableMutators {
				pattern := strings.HasSuffix(d, "*")

				if (pattern && strings.HasPrefix(name, d[:len(d)-2])) || (!pattern && name == d) {
					continue MUTATOR
				}
			}
		}

		m, _ := mutator.New(name)
		mutators = append(mutators, m)
	}

	os.Exit(returnOk)
}
