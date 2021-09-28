all: build

build:
	@go build

test: build

exec: build
	@./crontab

clean:
	@rm -f *~
	@go clean
