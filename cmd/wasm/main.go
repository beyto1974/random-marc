//go:build js && wasm

// Command wasm exposes generate.Records to the browser as
// window.generateMARC(count, format, seed).
package main

import (
	"bytes"
	"syscall/js"

	"random-marc/internal/generate"
)

var formatInfo = map[string]struct {
	mime string
	ext  string
}{
	"mrc":  {"application/marc", ".mrc"},
	"json": {"application/json", ".json"},
	"xml":  {"application/xml", ".xml"},
	"text": {"text/plain", ".txt"},
}

// generateMARC(count int, format string, seed number) -> {ok, data, mime, ext, error}
func generateMARC(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return errorResult("expected (count, format, seed)")
	}
	count := args[0].Int()
	format := args[1].String()
	seed := int64(args[2].Float())

	var buf bytes.Buffer
	if err := generate.Records(count, format, seed, &buf); err != nil {
		return errorResult(err.Error())
	}

	info, ok := formatInfo[format]
	if !ok {
		return errorResult("unknown format")
	}

	data := buf.Bytes()
	array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(array, data)

	return js.ValueOf(map[string]any{
		"ok":    true,
		"data":  array,
		"mime":  info.mime,
		"ext":   info.ext,
		"error": "",
	})
}

func errorResult(msg string) js.Value {
	return js.ValueOf(map[string]any{
		"ok":    false,
		"data":  js.Null(),
		"mime":  "",
		"ext":   "",
		"error": msg,
	})
}

func main() {
	js.Global().Set("generateMARC", js.FuncOf(generateMARC))
	select {}
}
