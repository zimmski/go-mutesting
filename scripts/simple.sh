#!/bin/bash

export GOMUTESTING_DIFF=$(diff -u $MUTATE_ORIGINAL $MUTATE_CHANGED)

mv $MUTATE_ORIGINAL $MUTATE_ORIGINAL.tmp
cp $MUTATE_CHANGED $MUTATE_ORIGINAL

go test ./... > /dev/null

export GOMUTESTING_RESULT=$?

mv $MUTATE_ORIGINAL.tmp $MUTATE_ORIGINAL

case $GOMUTESTING_RESULT in
0) # tests passed -> FAIL
	echo "$GOMUTESTING_DIFF"

	exit 1
	;;
1) # tests failed -> PASS
	exit 0
	;;
2) # did not compile -> SKIP
	echo "Mutation did not compile"
	echo "$GOMUTESTING_DIFF"

	exit 2
	;;
*) # Unkown exit code -> SKIP
	echo "$GOMUTESTING_DIFF"

	exit $GOMUTESTING_RESULT
	;;
esac
