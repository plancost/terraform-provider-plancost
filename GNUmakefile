default: fmt lint install

tools:
	@echo "==> installing required tooling..."
	go install github.com/katbyte/terrafmt@latest
	go install golang.org/x/tools/cmd/goimports@latest
build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

fmt:
	goimports -w .

terrafmt:
	@echo "==> Fixing acceptance test terraform blocks code with terrafmt..."
	@find internal -name "*_test.go" | grep -v "testdata" | sort | while read f; do terrafmt fmt -f $$f; done
	@echo "==> Fixing website terraform blocks code with terrafmt..."
	@find . \( -name "*.html.markdown" -o -name "*.md" \) | grep -v "testdata" | sort | while read f; do terrafmt fmt $$f; done

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

docs:
	go run tools/generate_resources/main.go

.PHONY: fmt lint test testacc build install docs
