version := $(shell /bin/date "+%Y-%m-%d %H:%M")

.PHONY: build
build:
	go build -o go-ctl -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" main.go
	$(if $(shell command -v upx), upx go-ctl)


.PHONY: clean
clean:
	rm go-ctl*