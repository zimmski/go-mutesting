# go-mutesting [![GoDoc](https://godoc.org/github.com/zimmski/go-mutesting?status.png)](https://godoc.org/github.com/zimmski/go-mutesting) [![Build Status](https://travis-ci.org/zimmski/go-mutesting.svg?branch=master)](https://travis-ci.org/zimmski/go-mutesting) [![Coverage Status](https://coveralls.io/repos/zimmski/go-mutesting/badge.png?branch=master)](https://coveralls.io/r/zimmski/go-mutesting?branch=master)

go-mutesting is a framework for performing mutation testing on Go source code.

## <a name="table-of-content"></a>Table of content

- [What is mutation testing?](#what-is-mutation-testing)
- [How do I use go-mutesting?](#how-do-i-use-go-mutesting)
- [How do I write my own mutation exec commands?](#write-mutation-exec-commands)
- [Which mutators are implemented?](#list-of-mutators)
- [Other mutation testing projects and their flaws](#other-projects)
- [Can I make feature requests and report bugs and problems?](#feature-request)

## <a name="what-is-mutation-testing"></a>What is mutation testing?

The definition of mutation testing is best quoted from Wikipedia:

> Mutation testing (or Mutation analysis or Program mutation) is used to design new software tests and evaluate the quality of existing software tests. Mutation testing involves modifying a program in small ways. Each mutated version is called a mutant and tests detect and reject mutants by causing the behavior of the original version to differ from the mutant. This is called killing the mutant. Test suites are measured by the percentage of mutants that they kill. New tests can be designed to kill additional mutants.
> <br/>-- <cite>[https://en.wikipedia.org/wiki/Mutation_testing](https://en.wikipedia.org/wiki/Mutation_testing)</cite>

> Tests can be created to verify the correctness of the implementation of a given software system, but the creation of tests still poses the question whether the tests are correct and sufficiently cover the requirements that have originated the implementation.
> <br/>-- <cite>[https://en.wikipedia.org/wiki/Mutation_testing](https://en.wikipedia.org/wiki/Mutation_testing)</cite>

Although the definition states that the main purpose of mutation testing is finding implementation cases which are not covered by tests, other implementation flaws can be found too. Mutation testing can for example uncover dead and unneeded code.

Mutation testing is also especially interesting for comparing automatically generated test suites with hand written test suites. This was the original intention of go-mutesting which is used to evaluate the generic fuzzing and delta-debugging framework [Tavor](https://github.com/zimmski/tavor).

## <a name="how-do-i-use-go-mutesting"></a>How do I use go-mutesting?

go-mutesting includes a binary which is go-getable.

```bash
go get -t -v github.com/zimmski/go-mutesting/cmd/go-mutesting
```

The binary's help can be invoked by executing the binary without arguments or with the `--help` option.

```bash
go-mutesting --help
```



TODO
	+ Timeouts sind erfolgreiche Fehler, daher wenn ein Timeout passiert wegen eine Mutation ist es ein PASS für die mutation


## <a name="write-mutation-exec-commands"></a>How do I write my own mutation exec commands?

A mutation exec command is invoked for every mutation which is necessary to test a mutation. Commands should handle at least the following phases.

1. **Setup** the source to include the mutation.
2. **Test** the source by invoking the test suite and possible other test functionality.
3. **Cleanup** revert all changes and remove all temporary assets.
4. **Report** if the mutation was detected.

It is important to note that each invocation should be isolated and therefore stateless. This means that an invocation must not disturb other invocations.

The command is given a set of environment variables which define exactly one mutation.

| Name            | Description                                                    |
| :-------------- | :------------------------------------------------------------- |
| MUTATE_ORIGINAL | Defines the filepath to the original file which was mutated.   |
| MUTATE_CHANGED  | Defines the filepath to the mutation of the original file.     |
| MUTATE_TIMEOUT  | Defines a timeout which should be honored by the exec command. |

A command must exit with an appropriate exit code.

| Exit code | Description                                                                                             |
| :------   | :--------                                                                                               |
| 0         | The mutation was killed. Which means that the test led to a failed test after the mutation was applied. |
| 1         | The mutation is alive. Which means that this could be a flaw.                                           |
| > 1       | The mutation was skipped, since there are other problems e.g. compilation errors.                       |

Examples for exec commands can be found in the [scripts](/scripts) directory.

## <a name="list-of-mutators"></a>Which mutators are implemented?

### Branch mutators

| Name        | Description                                    |
| :---------- | :--------------------------------------------- |
| branch/if   | Empties branches of if and else if statements. |
| branch/else | Empties branches of else statements.           |

### Expression mutators

| Name                | Description                                    |
| :------------------ | :--------------------------------------------- |
| expression/remove   | Searches for `&&` and <code>\|\|</code> expressions and makes each term of the expression unnecessary with using `true` or `false` as replacements. |

## <a name="write-mutators"></a>How do I write my own mutators?

Each mutator must implement the `Mutator` interface of the [github.com/zimmski/go-mutesting/mutator](https://godoc.org/github.com/zimmski/go-mutesting/mutator#Mutator) package. The methods of the interface are described in detail in the source code documentation.

Examples for mutators can be found in the [github.com/zimmski/go-mutesting/mutator](https://godoc.org/github.com/zimmski/go-mutesting/mutator) package and its sub-packages.

## <a name="other-projects"></a>Other mutation testing projects and their flaws

go-mutesting is not the first project to implement mutation testing for Go source code. A quick search search uncovers the following projects.

- https://github.com/darkhelmet/manbearpig
- https://github.com/kisielk/mutator
- https://github.com/StefanSchroeder/Golang-Mutation-testing

All of them have significant flaws in comparison to go-mutesting:

- Only one type (or even one case) of mutation is implemented
- Can only be used for one mutator at a time (manbearpig, Golang-Mutation-testing)
- Mutation is done by content which can lead to lots of invalid mutations (Golang-Mutation-testing)
- New mutators are not easily implemented and integrated
- Can only be used for one package or file at a time
- Other scenarios as `go test` cannot be applied
- Do not properly clean up or handle fatal failures
- No automatic tests to ensure that the algorithms are working at all
- Uses another language (Golang-Mutation-testing)

## <a name="feature-request"></a>Can I make feature requests and report bugs and problems?

Sure, just submit an [issue via the project tracker](https://github.com/zimmski/go-mutesting/issues/new) and I will see what I can do. Please note that I do not guarantee to implement anything soon and bugs and problems are more important to me than new features. If you need something implemented or fixed right away you can contact me via mail <mz@nethead.at> to do contract work for you.
