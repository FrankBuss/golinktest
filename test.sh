#!/bin/bash

# build the program
go build golinktest.go

# create an executable with the testdata integrated
TESTDATA=testdata.bin
OUTPUT=testout
./golinktest $TESTDATA $OUTPUT

# exececute the new output and extract the data
chmod a+x $OUTPUT
TESTDATA2=testdata2.bin
./$OUTPUT $TESTDATA2

# compare the result
if cmp $TESTDATA $TESTDATA2 ; then
    echo -e "\033[0;32mtest successful\033[m"
else
    echo -e "\033[0;31mtest failed\033[m"
fi
