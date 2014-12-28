#!/bin/bash

export GOMUTESTING_DIFF=$(diff -u $MUTATE_ORIGINAL $MUTATE_CHANGED)

mv $MUTATE_ORIGINAL $MUTATE_ORIGINAL.tmp
cp $MUTATE_CHANGED $MUTATE_ORIGINAL

go test ./... > /dev/null

export GOMUTESTING_RESULT=$?

mv $MUTATE_ORIGINAL.tmp $MUTATE_ORIGINAL

case $GOMUTESTING_RESULT in
0)
	echo "$GOMUTESTING_DIFF"

	exit 1
	;;
1)
	exit 0
	;;
*) # Unkown exit code
	echo "$GOMUTESTING_DIFF"

	exit $GOMUTESTING_RESULT
	;;
esac
