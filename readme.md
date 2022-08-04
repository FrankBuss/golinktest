# simple test for linking a file to an executable

Call the `golinktest` program with 2 arguments: the data file which you want to link, and the new executable program. Then you can call the new executable with one argument, which specifies the name of the data file to which the linked data is extracted.

See the [test script](test.sh) to see for an example and how to use it. It compiles the program, creates a new executable with it, then extracts the test data and compares it with the original.
