SOURCES = main.go
BINARY = letsencrypt-cert-days

BUILD_OPTIONS = -modcacherw
#BUILD_OPTIONS += -race
BUILD_OPTIONS += -ldflags="-w -X 'main.Version=$$(git tag -l --sort taggerdate | tail -1)' -extldflags '-static'"
BUILD_OPTIONS += -buildmode=pie
BUILD_OPTIONS += -tags netgo,osusergo
BUILD_OPTIONS += -trimpath

all: $(SOURCES) dependencies $(BINARY)

test:
	go test -v ./...

dependencies:
	go mod tidy

$(BINARY): $(SOURCES)
	env GOAMD64=v2 \
		CGO_ENABLED=1 \
	go build -o $(BINARY) $(BUILD_OPTIONS) .

tagpush: all
	./bin/stepup_tag.sh
	git push origin HEAD
	git push origin HEAD --tags

run:
	go run $(BUILD_OPTIONS) .

clean:
	rm -f $(BINARY)
