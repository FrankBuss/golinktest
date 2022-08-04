package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	magic = 0x56bb4200ff847612
)

type Info struct {
	Size  uint64
	Magic uint64
}

func readInfo(input *os.File) (Info, error) {
	i := Info{}
	input.Seek(-int64(binary.Size(i)), io.SeekEnd)
	err := binary.Read(input, binary.LittleEndian, &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

func hasInfo() (bool, error) {
	program, err := os.Open(os.Args[0])
	if err != nil {
		return false, err
	}
	defer program.Close()
	i, err := readInfo(program)
	if err != nil {
		return false, err
	}
	return i.Magic == magic, nil
}

func create(program_filename string, data_filename string, output_filename string) error {
	// open input files
	program, err := os.Open(program_filename)
	if err != nil {
		return err
	}
	defer program.Close()

	data, err := os.Open(data_filename)
	if err != nil {
		return err
	}
	defer data.Close()

	// create output file
	output, err := os.Create(output_filename)
	if err != nil {
		return err
	}
	defer output.Close()

	// copy program file to output
	_, err = io.Copy(output, program)
	if err != nil {
		return err
	}

	// append data to output
	_, err = io.Copy(output, data)
	if err != nil {
		return err
	}

	// append size of data and magic ID to the output
	fi, err := data.Stat()
	if err != nil {
		return err
	}
	i := Info{Size: uint64(fi.Size()), Magic: magic}
	err = binary.Write(output, binary.LittleEndian, i)
	if err != nil {
		return err
	}

	fmt.Printf("data file %v with size %v appended to program file and written to %v\n", data_filename, fi.Size(), output_filename)

	return nil
}

func extract(program_filename string, data_filename string) error {
	// open program file
	program, err := os.Open(program_filename)
	if err != nil {
		return err
	}
	defer program.Close()

	// create output file
	output, err := os.Create(data_filename)
	if err != nil {
		return err
	}
	defer output.Close()

	// read size of data from the end of the file
	i, err := readInfo(program)
	if err != nil {
		return err
	}
	if i.Magic != magic {
		return errors.New("no added data detected")
	}

	// go to start of data and copy it to the output
	program.Seek(int64(-(i.Size + uint64(binary.Size(i)))), io.SeekEnd)
	io.CopyN(output, program, int64(i.Size))
	if err != nil {
		return err
	}

	fmt.Printf("data file %v with size %v extracted\n", data_filename, i.Size)

	return nil
}

func usage(hasInfo bool) {
	if hasInfo {
		fmt.Printf("usage: %v datafile\n", os.Args[0])
		fmt.Printf("extracts the linked data\n")
		os.Exit(2)
	} else {
		fmt.Printf("usage: %v datafile outputfile\n", os.Args[0])
		fmt.Printf("links the data from 'datafile' to the program and creates a new executable named 'outputfile'\n")
		os.Exit(2)
	}
}

func main() {
	hasInfo, err := hasInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if hasInfo {
		// extract linked data
		if len(os.Args) != 2 {
			usage(hasInfo)
		}
		dataFilename := os.Args[1]
		if dataFilename == "" {
			usage(hasInfo)
		}
		err := extract(os.Args[0], dataFilename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		// link data
		if len(os.Args) != 3 {
			usage(hasInfo)
		}
		dataFilename := os.Args[1]
		outputFilename := os.Args[2]
		if dataFilename == "" || outputFilename == "" {
			usage(hasInfo)
		}
		err := create(os.Args[0], dataFilename, outputFilename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
