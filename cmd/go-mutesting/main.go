package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/zimmski/go-tool/importing"
	"github.com/zimmski/osutil"

	"github.com/zimmski/go-mutesting"
	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
	_ "github.com/zimmski/go-mutesting/mutator/branch"
	_ "github.com/zimmski/go-mutesting/mutator/expression"
	_ "github.com/zimmski/go-mutesting/mutator/statement"
)

const (
	returnOk = iota
	returnHelp
	returnBashCompletion
	returnError
)

type options struct {
	General struct {
		Debug                bool `long:"debug" description:"Debug log output"`
		DoNotRemoveTmpFolder bool `long:"do-not-remove-tmp-folder" description:"Do not remove the tmp folder where all mutations are saved to"`
		Help                 bool `long:"help" description:"Show this help message"`
		Verbose              bool `long:"verbose" description:"Verbose log output"`
	} `group:"General options"`

	Files struct {
		Blacklist []string `long:"blacklist" description:"List of MD5 checksums of mutations which should be ignored. Each checksum must end with a new line character."`
		ListFiles bool     `long:"list-files" description:"List found files"`
		PrintAST  bool     `long:"print-ast" description:"Print the ASTs of all given files and exit"`
	} `group:"File options"`

	Mutator struct {
		DisableMutators []string `long:"disable" description:"Disable mutator by their name or using * as a suffix pattern"`
		ListMutators    bool     `long:"list-mutators" description:"List all available mutators"`
	} `group:"Mutator options"`

	Filter struct {
		Match string `long:"match" description:"Only functions are mutated that confirm to the arguments regex"`
	} `group:"Filter options"`

	Exec struct {
		Exec    string `long:"exec" description:"Execute this command for every mutation (by default the built-in exec command is used)"`
		NoExec  bool   `long:"no-exec" description:"Skip the built-in exec command and just generate the mutations"`
		Timeout uint   `long:"exec-timeout" description:"Sets a timeout for the command execution (in seconds)" default:"10"`
	} `group:"Exec options"`

	Test struct {
		Recursive bool `long:"test-recursive" description:"Defines if the executer should test recursively"`
	} `group:"Test options"`

	Remaining struct {
		Targets []string `description:"Packages, directories and files even with patterns (by default the current directory)"`
	} `positional-args:"true" required:"true"`
}

func checkArguments(args []string, opts *options) (bool, int) {
	p := flags.NewNamedParser("go-mutesting", flags.None)

	p.ShortDescription = "Mutation testing for Go source code"

	if _, err := p.AddGroup("go-mutesting", "go-mutesting arguments", opts); err != nil {
		return true, exitError(err.Error())
	}

	completion := len(os.Getenv("GO_FLAGS_COMPLETION")) > 0

	_, err := p.ParseArgs(args)
	if (opts.General.Help || len(args) == 0) && !completion {
		p.WriteHelp(os.Stdout)

		return true, returnHelp
	} else if opts.Mutator.ListMutators {
		for _, name := range mutator.List() {
			fmt.Println(name)
		}

		return true, returnOk
	}

	if err != nil {
		return true, exitError(err.Error())
	}

	if completion {
		return true, returnBashCompletion
	}

	if opts.General.Debug {
		opts.General.Verbose = true
	}

	return false, 0
}

func debug(opts *options, format string, args ...interface{}) {
	if opts.General.Debug {
		fmt.Printf(format+"\n", args...)
	}
}

func verbose(opts *options, format string, args ...interface{}) {
	if opts.General.Verbose || opts.General.Debug {
		fmt.Printf(format+"\n", args...)
	}
}

func exitError(format string, args ...interface{}) int {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)

	return returnError
}

type mutatorItem struct {
	Name    string
	Mutator mutator.Mutator
}

type mutationStats struct {
	passed     int
	failed     int
	duplicated int
	skipped    int
}

func (ms *mutationStats) Score() float64 {
	total := ms.Total()

	if total == 0 {
		return 0.0
	}

	return float64(ms.passed) / float64(total)
}

func (ms *mutationStats) Total() int {
	return ms.passed + ms.failed
}

func mainCmd(args []string) int {
	var opts = &options{}
	var mutationBlackList = map[string]struct{}{}

	if exit, exitCode := checkArguments(args, opts); exit {
		return exitCode
	}

	files := importing.FilesOfArgs(opts.Remaining.Targets)
	if len(files) == 0 {
		return exitError("Could not find any suitable Go source files")
	}

	if opts.Files.ListFiles {
		for _, file := range files {
			fmt.Println(file)
		}

		return returnOk
	} else if opts.Files.PrintAST {
		for _, file := range files {
			fmt.Println(file)

			src, _, err := mutesting.ParseFile(file)
			if err != nil {
				return exitError("Could not open file %q: %v", file, err)
			}

			mutesting.PrintWalk(src)

			fmt.Println()
		}

		return returnOk
	}

	if len(opts.Files.Blacklist) > 0 {
		for _, f := range opts.Files.Blacklist {
			c, err := ioutil.ReadFile(f)
			if err != nil {
				return exitError("Cannot read blacklist file %q: %v", f, err)
			}

			for _, line := range strings.Split(string(c), "\n") {
				if line == "" {
					continue
				}

				if len(line) != 32 {
					return exitError("%q is not a MD5 checksum", line)
				}

				mutationBlackList[line] = struct{}{}
			}
		}
	}

	var mutators []mutatorItem

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

		verbose(opts, "Enable mutator %q", name)

		m, _ := mutator.New(name)
		mutators = append(mutators, mutatorItem{
			Name:    name,
			Mutator: m,
		})
	}

	tmpDir, err := ioutil.TempDir("", "go-mutesting-")
	if err != nil {
		panic(err)
	}
	verbose(opts, "Save mutations into %q", tmpDir)

	var execs []string
	if opts.Exec.Exec != "" {
		execs = strings.Split(opts.Exec.Exec, " ")
	}

	stats := &mutationStats{}

	for _, file := range files {
		verbose(opts, "Mutate %q", file)

		src, fset, pkg, info, err := mutesting.ParseAndTypeCheckFile(file)
		if err != nil {
			return exitError(err.Error())
		}

		err = os.MkdirAll(tmpDir+"/"+filepath.Dir(file), 0755)
		if err != nil {
			panic(err)
		}

		tmpFile := tmpDir + "/" + file

		originalFile := fmt.Sprintf("%s.original", tmpFile)
		err = osutil.CopyFile(file, originalFile)
		if err != nil {
			panic(err)
		}
		debug(opts, "Save original into %q", originalFile)

		mutationID := 0

		if opts.Filter.Match != "" {
			m, err := regexp.Compile(opts.Filter.Match)
			if err != nil {
				return exitError("Match regex is not valid: %v", err)
			}

			for _, f := range astutil.Functions(src) {
				if m.MatchString(f.Name.Name) {
					mutationID = mutate(opts, mutators, mutationBlackList, mutationID, pkg, info, file, fset, src, f, tmpFile, execs, stats)
				}
			}
		} else {
			_ = mutate(opts, mutators, mutationBlackList, mutationID, pkg, info, file, fset, src, src, tmpFile, execs, stats)
		}
	}

	if !opts.General.DoNotRemoveTmpFolder {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			panic(err)
		}
		debug(opts, "Remove %q", tmpDir)
	}

	if !opts.Exec.NoExec {
		fmt.Printf("The mutation score is %f (%d passed, %d failed, %d duplicated, %d skipped, total is %d)\n", stats.Score(), stats.passed, stats.failed, stats.duplicated, stats.skipped, stats.Total())
	} else {
		fmt.Println("Cannot do a mutation testing summary since no exec command was executed.")
	}

	return returnOk
}

func mutate(opts *options, mutators []mutatorItem, mutationBlackList map[string]struct{}, mutationID int, pkg *types.Package, info *types.Info, file string, fset *token.FileSet, src ast.Node, node ast.Node, tmpFile string, execs []string, stats *mutationStats) int {
	for _, m := range mutators {
		debug(opts, "Mutator %s", m.Name)

		changed := mutesting.MutateWalk(pkg, info, node, m.Mutator)

		for {
			_, ok := <-changed

			if !ok {
				break
			}

			mutationFile := fmt.Sprintf("%s.%d", tmpFile, mutationID)
			checksum, duplicate, err := saveAST(mutationBlackList, mutationFile, fset, src)
			if err != nil {
				fmt.Printf("INTERNAL ERROR %s\n", err.Error())
			} else if duplicate {
				debug(opts, "%q is a duplicate, we ignore it", mutationFile)

				stats.duplicated++
			} else {
				debug(opts, "Save mutation into %q with checksum %s", mutationFile, checksum)

				if !opts.Exec.NoExec {
					execExitCode := mutateExec(opts, pkg, file, src, mutationFile, execs)

					debug(opts, "Exited with %d", execExitCode)

					msg := fmt.Sprintf("%q with checksum %s", mutationFile, checksum)

					switch execExitCode {
					case 0:
						fmt.Printf("PASS %s\n", msg)

						stats.passed++
					case 1:
						fmt.Printf("FAIL %s\n", msg)

						stats.failed++
					case 2:
						fmt.Printf("SKIP %s\n", msg)

						stats.skipped++
					default:
						fmt.Printf("UNKOWN exit code for %s\n", msg)
					}
				}
			}

			changed <- true

			// Ignore original state
			<-changed
			changed <- true

			mutationID++
		}
	}

	return mutationID
}

func mutateExec(opts *options, pkg *types.Package, file string, src ast.Node, mutationFile string, execs []string) (execExitCode int) {
	if len(execs) == 0 {
		debug(opts, "Execute built-in exec command for mutation")

		diff, err := exec.Command("diff", "-u", file, mutationFile).CombinedOutput()
		if err == nil {
			execExitCode = 0
		} else if e, ok := err.(*exec.ExitError); ok {
			execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			panic(err)
		}
		if execExitCode != 0 && execExitCode != 1 {
			fmt.Printf("%s\n", diff)

			panic("Could not execute diff on mutation file")
		}

		defer func() {
			_ = os.Rename(file+".tmp", file)
		}()

		err = os.Rename(file, file+".tmp")
		if err != nil {
			panic(err)
		}
		err = osutil.CopyFile(mutationFile, file)
		if err != nil {
			panic(err)
		}

		pkgName := pkg.Path()
		if opts.Test.Recursive {
			pkgName += "/..."
		}

		test, err := exec.Command("go", "test", "-timeout", fmt.Sprintf("%ds", opts.Exec.Timeout), pkgName).CombinedOutput()
		if err == nil {
			execExitCode = 0
		} else if e, ok := err.(*exec.ExitError); ok {
			execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			panic(err)
		}

		if opts.General.Debug {
			fmt.Printf("%s\n", test)
		}

		switch execExitCode {
		case 0: // Tests passed -> FAIL
			fmt.Printf("%s\n", diff)

			execExitCode = 1
		case 1: // Tests failed -> PASS
			if opts.General.Debug {
				fmt.Printf("%s\n", diff)
			}

			execExitCode = 0
		case 2: // Did not compile -> SKIP
			if opts.General.Verbose {
				fmt.Println("Mutation did not compile")
			}

			if opts.General.Debug {
				fmt.Printf("%s\n", diff)
			}
		default: // Unknown exit code -> SKIP
			fmt.Println("Unknown exit code")
			fmt.Printf("%s\n", diff)
		}

		return execExitCode
	}

	debug(opts, "Execute %q for mutation", opts.Exec.Exec)

	execCommand := exec.Command(execs[0], execs[1:]...)

	execCommand.Stderr = os.Stderr
	execCommand.Stdout = os.Stdout

	execCommand.Env = append(os.Environ(), []string{
		"MUTATE_CHANGED=" + mutationFile,
		fmt.Sprintf("MUTATE_DEBUG=%t", opts.General.Debug),
		"MUTATE_ORIGINAL=" + file,
		"MUTATE_PACKAGE=" + pkg.Path(),
		fmt.Sprintf("MUTATE_TIMEOUT=%d", opts.Exec.Timeout),
		fmt.Sprintf("MUTATE_VERBOSE=%t", opts.General.Verbose),
	}...)
	if opts.Test.Recursive {
		execCommand.Env = append(execCommand.Env, "TEST_RECURSIVE=true")
	}

	err := execCommand.Start()
	if err != nil {
		panic(err)
	}

	// TODO timeout here

	err = execCommand.Wait()

	if err == nil {
		execExitCode = 0
	} else if e, ok := err.(*exec.ExitError); ok {
		execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
	} else {
		panic(err)
	}

	return execExitCode
}

func main() {
	os.Exit(mainCmd(os.Args[1:]))
}

func saveAST(mutationBlackList map[string]struct{}, file string, fset *token.FileSet, node ast.Node) (string, bool, error) {
	var buf bytes.Buffer

	h := md5.New()

	err := printer.Fprint(io.MultiWriter(h, &buf), fset, node)
	if err != nil {
		return "", false, err
	}

	checksum := fmt.Sprintf("%x", h.Sum(nil))

	if _, ok := mutationBlackList[checksum]; ok {
		return checksum, true, nil
	}

	mutationBlackList[checksum] = struct{}{}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return "", false, err
	}

	err = ioutil.WriteFile(file, src, 0666)
	if err != nil {
		return "", false, err
	}

	return checksum, false, nil
}
