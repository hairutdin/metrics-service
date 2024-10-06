SHELL=/bin/bash

MODULES = agent server

PLATFORMS = linux windows darwin
ARCHITECTURES = amd64 arm64

all: tidy build test

tidy:
	go mod tidy

build:
	$(foreach MODULE,$(MODULES), \
		$(foreach OS,$(PLATFORMS), \
			$(foreach ARCH,$(ARCHITECTURES), \
				GOOS=$(OS) GOARCH=$(ARCH) go build -o bin/$(MODULE)-$(OS)-$(ARCH) ./cmd/$(MODULE) ; \
			) \
		) \
	)

test:
	$(foreach MODULE,$(MODULES), \
		go test -v ./cmd/$(MODULE)/... ; \
	)

clean:
	rm -rf ./bin

perm:
	chmod -R +x bin

