#!/bin/bash

# create a 1 GB test file with random data
TESTDATA=testdata.bin
./golinktest --action=random --output=$TESTDATA --size=1_000_000_000

# compile the Go file
go build golinktest.go

# create an executable with the testdata integrated
OUTPUT=testout
./golinktest --action=create --data=$TESTDATA --output=$OUTPUT

# exececute the new output and extract the data
chmod a+x $OUTPUT
TESTDATA2=testdata2.bin
./$OUTPUT --action=extract --output=$TESTDATA2

# compare the result
cmp $TESTDATA $TESTDATA2
