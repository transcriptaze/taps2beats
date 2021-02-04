// +build js,wasm

// WASM interface for the taps2beats functions
package main

import (
"syscall/js"

//	"github.com/twystd/taps2beats/taps2beats"
)

const VERSION = "v0.1.0"

func main() {
	c := make(chan bool)

	js.Global().Set("goTaps", js.FuncOf(taps))

	<-c
}

func taps(this js.Value, inputs []js.Value) interface{} {
	callback := inputs[0]

	go func() {
		// if err := dispatcher.Redraw(); err != nil {
		// 	callback.Invoke(err.Error())
		// 	return
		// }

		// callback.Invoke(js.Null())
		callback.Invoke("qwerty-uiop")
	}()

	return nil
}
