package main

import (
	"fmt"
	gobencode "joshua/green/bencode"
	"os"
)

type Package struct {
	Name  string
	Date  int
	Peers []string
}

func main() {
	FileDecodeExample()
	StringDecodeExample()

	EncodeExample1()
	EncodeExample2()
}

func FileDecodeExample() {
	fmt.Println("\n\n\t\tFileDecodeExample")
	res := gobencode.DecodeFile("sample.torrent")
	fmt.Println(res)
}

func StringDecodeExample() {
	fmt.Println("\n\n\t\tStringDecodeExample")
	res := gobencode.DecodeString("l5:hello5:worldi12345ed6:animal3:cat4:typei1eee")
	fmt.Println(res)
}

func EncodeExample1() {
	fmt.Println("\n\n\t\tEncodeExample1")
	e := gobencode.NewEncoder(os.Stdout)

	list := gobencode.DecodeString("l5:hello5:worldi12345ed6:animal3:cat4:typei1eee")
	e.Encode(list)
}

func EncodeExample2() {
	fmt.Println("\n\n\t\tEncodeExample2")

	e := gobencode.NewEncoder(os.Stdout)

	e.Encode(1234)

	fmt.Println()
	e.Encode([]string{"hello", "world"})

	fmt.Println()
	e.Encode(map[string]int{"a": 1, "b": 2, "c": 3})
}
