package main

import (
	gobencode "joshua/green/bencode"
	"os"
)

type Package struct {
	Name  string   `bencode:"name"`
	Date  int      `bencode:"date"`
	Peers []string `bencode:"peers"`
}

func main() {
	enc := gobencode.NewEncoder(os.Stdout)
	// enc.Encode(52)
	// enc.Encode("cats are cool")
	// enc.Encode([]int{1, 2, 3, 4, 5, 6})
	// enc.Encode([]string{"cats", "are", "cool"})

	enc.Encode(Package{
		Name:  "gobencode",
		Date:  1354534,
		Peers: []string{"joshua", "codi", "mark"},
	})
}
