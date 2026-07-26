.PHONY: wasm serve

wasm:
	GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/wasm
	rm -f web/wasm_exec.js
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js
	chmod u+w web/wasm_exec.js

serve: wasm
	cd web && python3 -m http.server 8080
