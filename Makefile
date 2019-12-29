.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o assets/wasm/wa_output.wasm github.com/joshprzybyszewski/cribbage/wasm

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor