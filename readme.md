# simple test for linking a file to an executable

```
Usage of ./golinktest:
  -action string
        Action to make: {create|extract|random}.
  -data string
        Filename of the data for the create action.
  -output string
        The output executable or data, for the create, extract, or random command.
  -size uint
        The size of the random output file for the random command.
```

With `create` you can create a new executable, to which the data file is appended to the end. It also appends the size of the data file after this at the end. This allows `extract` to read the data file from the program file, and save it as a separate file. With the `random` action, and the `size` and `output`, you can create a file with random bytes.

See the [test script](test.sh) to see for an example and how to use it. It creates a 1 GB file for testing, creates a new executable with it, then extracts the test data and compares it with the original.

TODO: would be nice to have auto-extract for the generated new binary file, probably by adding a unique marker at the end to tag it that it contains data, and a string for the data filename, and to check if there is already a file appended, to avoid appending multiple files.
