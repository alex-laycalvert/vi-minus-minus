package main

import (
	"fmt"
	"os"

	"github.com/alex-laycalvert/vi-minus-minus/buffer"
	"github.com/alex-laycalvert/vi-minus-minus/vimm"
)

func main() {
	buf := buffer.New()
	buf.AppendLine("")
	v, err := vimm.New()
	checkError(err)
	defer quit(v)
	v.AddBuffer(buffer.FromString(""))
	for {
		v.Show()
		if v.ProcessEvent() {
			break
		}
	}
}

func quit(view *vimm.Vimm) {
	view.End()
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
