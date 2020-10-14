package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal(fmt.Sprintf("Usage: %s input.mid output.mid", os.Args[0]))
	}

	inputfile, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	input, err := ReadMidi(inputfile)

	if err != nil {
		log.Fatal(err)
	}

	result := cleanupMidi(input)

	outputFile, err := os.OpenFile(os.Args[2], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)

	err = WriteMidi(outputFile, result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Processed file %s and saved it to %s", os.Args[1], os.Args[2]))
}
