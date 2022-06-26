package main

import (
	"fmt"
	gobencode "joshua/green/bencode"
)

type Package struct {
	Name  string
	Date  int
	Peers []string
}

func main() {
	res := gobencode.DecodeFile("sample.torrent")
	fmt.Println(res.String())
}

// func example_encode(e *gobencode.Encoder) {
// 	e.Encode(52)
// 	e.Encode("cats are cool")
// 	e.Encode([]int{1, 2, 3, 4, 5, 6})
// 	e.Encode([]string{"cats", "are", "cool"})
// 	fmt.Println()
// }

// func example_struct_encode(e *gobencode.Encoder) {
// 	e.Encode(Package{
// 		Name:  "gobencode",
// 		Date:  1354534,
// 		Peers: []string{"joshua", "codi", "mark"},
// 	})
// 	fmt.Println()
// }
