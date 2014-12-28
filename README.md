# go-mutesting [![GoDoc](https://godoc.org/github.com/zimmski/go-mutesting?status.png)](https://godoc.org/github.com/zimmski/go-mutesting) [![Build Status](https://travis-ci.org/zimmski/go-mutesting.svg?branch=master)](https://travis-ci.org/zimmski/go-mutesting) [![Coverage Status](https://coveralls.io/repos/zimmski/go-mutesting/badge.png?branch=master)](https://coveralls.io/r/zimmski/go-mutesting?branch=master)

go-mutesting is a framework for performing mutation testing on Go source code. The definition of mutation testing is best quoted from Wikipedia:

> Mutation testing (or Mutation analysis or Program mutation) is used to design new software tests and evaluate the quality of existing software tests. Mutation testing involves modifying a program in small ways. Each mutated version is called a mutant and tests detect and reject mutants by causing the behavior of the original version to differ from the mutant. This is called killing the mutant. Test suites are measured by the percentage of mutants that they kill. New tests can be designed to kill additional mutants.
> <br/>-- <cite>[https://en.wikipedia.org/wiki/Mutation_testing](https://en.wikipedia.org/wiki/Mutation_testing)</cite>

Zusätzliches Zitat von Wikipedia:

> Tests can be created to verify the correctness of the implementation of a given software system, but the creation of tests still poses the question whether the tests are correct and sufficiently cover the requirements that have originated the implementation.


TODO

- Text fertig schreiben
	+ ...
	+ Die drei Projekte erwähnen und Danken, und dann schreiben dass die Lösungen einfach nicht passen für mich waren und ich lieber eine allgemeinere Lösung haben wollte die zudem mehr automatisierbar ist, und leichter erweiterbar TODO auch code vorher anscahuen ob etwas brauchbares drinnen ist
		* http://verboselogging.com/2013/03/11/manbearpig-mutation-testing-for-go
		* https://github.com/StefanSchroeder/Golang-Mutation-testing
		* drittes -> siehe git clone ordner
		* https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/R2VCbCPA6ZE auch erwähnen
DONE- Wir wollen folgendes unterstützen
	+ Statements löschen die löschbar sind, das heißt
		* Keine Deklarationen aller Art
		* Keine Kommentare, weil das wäre unnötig
		* Also rein nur Zuweisungen und Aufrüfe
DONE	+ Blöcke löschen, daher
DONE		* if, else if und else branches leeren
	+ Condition Terme löschen vom ersten Condition level wenn mehr als ein Term da ist, für
		* ifs und elseifs
DONE - Jede Art von "Mutator" ist ein "Plugin"
DONE- CMD
DONE	+ Kann wie Golint auf alles mögliche angewendet werden: einzelne Dateien, Ordner, Packages, mit Patterns
DONE	+ Standardmäßig sind alle Mutator aktiviert, man kann sie aber einzeln oder per Pattern deaktivieren
DONE	+ Mutatoren werden deterministisch ausgeführt, und zwar nach name sortiert
DONE	+ Skript kann angegeben werden welches die Mutation überprüfen soll, hat drei Antwortmöglichkeiten: Passed, Failed und wenn etwas mit der Mutation nicht passt "Skipped"
DONE	+ mutation score = number of mutants killed / total number of mutants
- Testscript schreiben für basis sachen -> projekt vom aktuellen ordner testen + datei austauschen bzw austausch rückgängig machen
