package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jessevdk/go-flags"

	"github.com/zimmski/go-mutesting"
	"github.com/zimmski/go-mutesting/importing"
	"github.com/zimmski/go-mutesting/mutator"
	_ "github.com/zimmski/go-mutesting/mutator/branch"
	_ "github.com/zimmski/go-mutesting/mutator/expression"
)

const (
	returnOk = iota
	returnHelp
	returnBashCompletion
	returnError
)

const (
	execPassed  = 0
	execFailed  = 1
	execSkipped = 2
)

var opts struct {
	General struct {
		Debug                bool `long:"debug" description:"Debug log output"`
		DoNotRemoveTmpFolder bool `long:"do-not-remove-tmp-folder" description:"Do not remove the tmp folder where all mutations are saved to"`
		Help                 bool `long:"help" description:"Show this help message"`
		Verbose              bool `long:"verbose" description:"Verbose log output"`
	} `group:"General options"`

	Files struct {
		ListFiles bool `long:"list-files" description:"List found files"`
		PrintAST  bool `long:"print-ast" description:"Print the ASTs of all given files and exit"`
	} `group:"File options"`

	Mutator struct {
		DisableMutators []string `long:"disable" description:"Disable mutator or mutators using * as a suffix pattern"`
		ListMutators    bool     `long:"list-mutators" description:"List all available mutators"`
	} `group:"Mutator options"`

	Exec struct {
		Exec    string `long:"exec" description:"Execute this command for every mutation"`
		Timeout uint   `long:"exec-timeout" description:"Sets a timeout for the command execution (in seconds)" default:"10"`
	} `group:"Exec options"`

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

func debug(format string, args ...interface{}) {
	if opts.General.Debug {
		fmt.Printf(format+"\n", args...)
	}
}

func verbose(format string, args ...interface{}) {
	if opts.General.Verbose || opts.General.Debug {
		fmt.Printf(format+"\n", args...)
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

	if opts.Files.ListFiles {
		for _, file := range files {
			fmt.Println(file)
		}

		os.Exit(returnOk)
	} else if opts.Files.PrintAST {
		for _, file := range files {
			fmt.Println(file)

			src, _, err := mutesting.ParseFile(file)
			if err != nil {
				exitError("Could not open file %q: %v", file, err)
			}

			mutesting.PrintWalk(src)

			fmt.Println()
		}

		os.Exit(returnOk)
	}

	var mutators []mutator.Mutator

MUTATOR:
	for _, name := range mutator.List() {
		if len(opts.Mutator.DisableMutators) > 0 {
			for _, d := range opts.Mutator.DisableMutators {
				pattern := strings.HasSuffix(d, "*")

				if (pattern && strings.HasPrefix(name, d[:len(d)-2])) || (!pattern && name == d) {
					continue MUTATOR
				}
			}
		}

		debug("Enable mutator %q", name)

		m, _ := mutator.New(name)
		mutators = append(mutators, m)
	}

	tmpDir, err := ioutil.TempDir("", "go-mutesting-")
	if err != nil {
		panic(err)
	}
	debug("Save mutations into %q", tmpDir)

	var execs []string
	if opts.Exec.Exec != "" {
		execs = strings.Split(opts.Exec.Exec, " ")
	}

	passed := 0
	failed := 0
	skipped := 0

	for _, file := range files {
		debug("Mutate %q", file)

		src, fset, err := mutesting.ParseFile(file)
		if err != nil {
			exitError("Could not open file %q: %v", file, err)
		}

		err = os.MkdirAll(tmpDir+"/"+filepath.Dir(file), 0755)
		if err != nil {
			panic(err)
		}

		tmpFile := tmpDir + "/" + file

		originalFile := fmt.Sprintf("%s.original", tmpFile)
		err = copyFile(file, originalFile)
		if err != nil {
			panic(err)
		}
		debug("Save original into %q", originalFile)

		i := 0

		for _, m := range mutators {
			debug("Mutator %s", m)

			changed := mutesting.MutateWalk(src, m)

			for {
				_, ok := <-changed

				if !ok {
					break
				}

				mutationFile := fmt.Sprintf("%s.%d", tmpFile, i)
				err = saveAST(mutationFile, fset, src)
				if err != nil {
					panic(err)
				}
				debug("Save mutation into %q", mutationFile)

				if len(execs) > 0 {
					debug("Execute %q for mutation", opts.Exec.Exec)

					execCommand := exec.Command(execs[0], execs[1:]...)

					execCommand.Stderr = os.Stderr
					execCommand.Stdout = os.Stdout

					execCommand.Env = append(os.Environ(), []string{
						"MUTATE_CHANGED=" + mutationFile,
						"MUTATE_ORIGINAL=" + file,
						fmt.Sprintf("MUTATE_TIMEOUT=%d", opts.Exec.Timeout),
					}...)

					err = execCommand.Start()
					if err != nil {
						panic(err)
					}

					// TODO timeout here

					err = execCommand.Wait()

					var execExitCode int
					if err == nil {
						execExitCode = 0
					} else if e, ok := err.(*exec.ExitError); ok {
						execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
					} else {
						panic(err)
					}

					debug("Exited with %d", execExitCode)

					switch execExitCode {
					case 0:
						fmt.Printf("PASS %q\n", mutationFile)

						passed++
					case 1:
						fmt.Printf("FAIL %q\n", mutationFile)

						failed++
					case 2:
						fmt.Printf("SKIP %q\n", mutationFile)

						skipped++
					default:
						fmt.Printf("UNKOWN exit code for %q\n", mutationFile)
					}
				}

				changed <- true

				// ignore original state
				<-changed
				changed <- true

				i++
			}
		}
	}

	if !opts.General.DoNotRemoveTmpFolder {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			panic(err)
		}
		debug("Remove %q", tmpDir)
	}

	if len(execs) > 0 {
		fmt.Printf("The mutation score is %f (%d passed, %d failed, %d skipped, total is %d)\n", float64(passed)/float64(passed+failed), passed, failed, skipped, passed+failed+skipped)
	} else {
		fmt.Println("Cannot do a mutation testing summary since no exec command was given.")
	}

	os.Exit(returnOk)
}

func copyFile(src string, dst string) (err error) {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		e := s.Close()
		if err == nil {
			err = e
		}
	}()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		e := d.Close()
		if err == nil {
			err = e
		}
	}()

	_, err = io.Copy(d, s)
	if err == nil {
		i, err := os.Stat(src)
		if err != nil {
			return err
		}

		return os.Chmod(dst, i.Mode())
	}

	return nil
}

func saveAST(file string, fset *token.FileSet, node ast.Node) error {
	var buf bytes.Buffer

	err := printer.Fprint(&buf, fset, node)
	if err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, src, 0666)
}
