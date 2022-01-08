build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	GOARCH=amd64 GOOS=linux go build -o extFileRecovery