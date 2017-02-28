package main

import (
	"bytes"
	"flag"
	"fmt"

	"./src"
)

func main() {
	var err error
	var fptr, qptr string
	var fmr, qmr *wapiq.MapResult

	flag.StringVar(&fptr, "f", "", "-f=\"<filename>\"")
	flag.StringVar(&qptr, "q", "", "-q=\"<string>\"")
	flag.Parse()
	if len(fptr) == 0 && len(qptr) == 0 {
		fmt.Println("WAPIQ CLI\nusage:\n",
			"\t-f <filename>      Load from file any query results.\n",
			"\t-q  <string>        Execute query from passed string.\n",
		)
		return
	}
	w := (&wapiq.WAPIQ{}).New()

	if len(fptr) > 0 {
		fmr, err = w.Load(fptr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if len(qptr) > 0 {
		qmr, err = w.Query(qptr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	jfmr := ""
	jqmr := ""
	r := bytes.Buffer{}
	if qmr != nil && fmr != nil {
		r.WriteRune('[')
	}
	if fmr != nil {
		jfmr, err = fmr.JSON()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(jfmr) > 0 {
			r.WriteString(jfmr)
		}
	}
	if qmr != nil {
		jqmr, err = qmr.JSON()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	if len(jfmr) > 0 && len(jqmr) > 0 {
		r.WriteRune(',')
	}
	if len(jqmr) > 0 {
		r.WriteString(jqmr)
	}

	if qmr != nil && fmr != nil {
		r.WriteRune(']')
	}

	fmt.Println(r.String())
}
