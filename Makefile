.PHONY: test install-editor-settings

test:
	TEST_DATA_ROOT=$(shell pwd)/testdata go test -v ./...

# Copies vscode/settings.json to .vscode/ so the Go extension gets TEST_DATA_ROOT when running tests in the editor (run once per clone).
install-editor-settings:
	mkdir -p .vscode
	cp vscode/settings.json .vscode/settings.json
