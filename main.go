package main

import (
	"fmt"
	"os"

	"github.com/alex-laycalvert/vimm/app"
	"github.com/alex-laycalvert/vimm/buffer"
)

func main() {
	buf := buffer.New()
	buf.AppendLine("")
	v, err := app.New()
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

func quit(view *app.App) {
	view.End()
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
