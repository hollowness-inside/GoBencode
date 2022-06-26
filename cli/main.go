package main

import (
	"fmt"
	gobencode "joshua/green/bencode"
	"os"
)

type Operation byte

const (
	Decode Operation = 1
	Encode Operation = 2
)

type Arguments struct {
	Op            Operation
	HasInputFile  bool
	InputFile     string
	Text          string
	OutputFile    string
	HasOutputFile bool
}

func help() {
	fmt.Println("gobencode [-e | -d [-i input_file]] [-o output_file] [text]")
	fmt.Println("\t-d - Set mode to decode")
	fmt.Println("\t-e - Set mode to encode")
	fmt.Println("\t-i filepath - input file")
	fmt.Println("\t-o filepath - redirects stdout to file")
}

func main() {
	args := Arguments{}

	i := 0
	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-d":
			args.Op = Decode
		case "-e":
			args.Op = Encode
			args.HasInputFile = true
			args.InputFile = os.Args[i+1]
			i += 1
		case "-i":
			args.HasInputFile = true
			args.InputFile = os.Args[i+1]
			i += 1
		case "-o":
			args.HasOutputFile = true
			args.OutputFile = os.Args[i+1]
			i += 1
		default:
			args.Text = os.Args[i]
			break
		}

		i += 1
	}

	stream := os.Stdout

	if args.HasOutputFile {
		var err error
		stream, err = os.Create(args.OutputFile)
		if err != nil {
			panic(err)
		}
		defer stream.Close()
	}

	if args.Op == Decode {
		var item gobencode.BencodeItem
		if args.HasInputFile {
			item = gobencode.DecodeFile(args.InputFile)
		} else {
			item = gobencode.DecodeString(args.Text)
		}
		fmt.Fprint(stream, item)
	} else if args.Op == Encode {
		if args.HasInputFile {
			EncodeJsonFile(args.InputFile)
		} else {
			EncodeJsonString(args.Text)
		}
	} else {
		help()
	}
}
