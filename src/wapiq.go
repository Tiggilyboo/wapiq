package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"

	"../src/parser"
)

func main() {
	var fptr, sptr string
	flag.StringVar(&fptr, "f", "", "-f <filename>")
	flag.StringVar(&sptr, "s", "", "-s <string>")
	flag.Parse()

	if len(fptr) == 0 && len(sptr) == 0 {
		fmt.Println("WAPIQ CLI\nusage: \n\t-f <filename>\n\t-s <string>")
		return
	}
	var err error
	var b []byte
	if len(fptr) > 0 {
		fmt.Println(fptr)
		b, err = ioutil.ReadFile(fptr)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		b = []byte(sptr)
	}
	r := bytes.NewReader(b)
	p := parser.NewParser(r)
	a, err := p.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%v", a)
}
