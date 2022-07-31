package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
)

func append_file(output *os.File, input *os.File, limit int64) error {
	bytes := make([]byte, 1024)
	for {
		n, err := input.Read(bytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if int64(n) > limit {
			n = int(limit)
		}
		limit -= int64(n)
		_, err = output.Write(bytes[0:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func create(program_filename string, data_filename string, output_filename string) int {
	// open input files
	program, err := os.Open(program_filename)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	data, err := os.Open(data_filename)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// create output file
	output, err := os.Create(output_filename)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// copy program file to output
	err = append_file(output, program, math.MaxInt64)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// append data to output
	append_file(output, data, math.MaxInt64)

	// append size of data to the output
	fi, err := data.Stat()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(fi.Size()))
	output.Write(b)

	fmt.Printf("data file %v with size %v appended to program file and written to %v\n", data_filename, fi.Size(), output_filename)

	return 0
}

func extract(program_filename string, data_filename string) int {
	// open program file
	program, err := os.Open(program_filename)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// create output file
	output, err := os.Create(data_filename)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// read size of data from the end of the file
	program.Seek(-8, 2)
	size_bytes := make([]byte, 8)
	_, err = program.Read(size_bytes)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	size := binary.LittleEndian.Uint64(size_bytes)

	// go to start of data and copy it to the output
	program.Seek(int64(-(size + 8)), 2)
	append_file(output, program, int64(size))

	fmt.Printf("data file %v with size %v extracted\n", data_filename, size)

	return 0
}

func random(output_filename string, size uint64) int {
	// create output file
	output, err := os.Create(output_filename)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// write random bytes
	bytes := make([]byte, 1024)
	for {
		n := len(bytes)
		if uint64(n) > size {
			n = int(size)
		}
		for i := 0; i < n; i++ {
			bytes[i] = byte(rand.Uint32())
		}
		_, err = output.Write(bytes[0:n])
		if err != nil {
			fmt.Println(err)
			return 1
		}
		size -= uint64(n)
		if size == 0 {
			break
		}
	}

	fmt.Printf("random file %v with size %v created\n", output_filename, size)

	return 0
}

func main() {
	if len(os.Args) > 1 {
		action := flag.String("action", "", "Action to make: {create|extract|random}.")
		data_filename := flag.String("data", "", "Filename of the data for the create action.")
		output_filename := flag.String("output", "", "The output executable or data, for the create, extract, or random command.")
		random_size := flag.Uint64("size", 0, "The size of the random output file for the random command.")
		flag.Parse()
		if *action == "" {
			goto Usage
		}
		fmt.Printf("action: %v\n", *action)
		switch *action {
		case "create":
			if *data_filename == "" || *output_filename == "" {
				goto Usage
			}
			os.Exit(create(os.Args[0], *data_filename, *output_filename))
		case "extract":
			if *output_filename == "" {
				goto Usage
			}
			os.Exit(extract(os.Args[0], *output_filename))
		case "random":
			fmt.Printf("size: %v bytes\n", *random_size)
			if *random_size == 0 {
				goto Usage
			}
			os.Exit(random(*output_filename, *random_size))
		}
	} else {
		goto Usage
	}

Usage:
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	flag.PrintDefaults()
	goto Error

Error:
	os.Exit(2)
}
