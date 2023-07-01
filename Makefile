BUILDDIR := ./build

build:  clean prep
	go build -o $(BUILDDIR)/ ./...

clean:
	@rm -rf $(BUILDDIR)
	@rm -rf ./vendor

prep:
	@mkdir -p $(BUILDDIR)
	@go mod tidy
	@go mod vendor

run: build
	@$(BUILDDIR)/lilly
