#!/bin/bash

# build the program
go build golinktest.go

# create a 1 MB test file with random data
TESTDATA=testdata.bin
./golinktest --action=random --output=$TESTDATA --size=1_000_000

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
if cmp $TESTDATA $TESTDATA2 ; then
    echo -e "\033[0;32mtest successful\033[m"
else
    echo -e "\033[0;31mtest failed\033[m"
fi
