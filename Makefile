.PHONY: wasm serve

wasm:
	GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/wasm
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js

serve: wasm
	cd web && python3 -m http.server 8080
