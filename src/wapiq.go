package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"../src/parser"
)

func main() {
	b, err := ioutil.ReadFile("test.wapiq")
	if err != nil {
		fmt.Println(err.Error())
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
