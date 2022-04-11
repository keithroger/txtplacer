PROJECTPATH=$(CURDIR)
TESTDATAPATH=$(CURDIR)/testdata
TESTIMAGES=$(TESTDATAPATH)/images

clean:
	go clean
	@if [ -d $(TESTIMAGES) ] ; then rm -r $(TESTIMAGES) ; fi

imagedir:
	@if [ ! -d $(TESTIMAGES) ] ; then mkdir -p $(TESTIMAGES) ; fi

test: imagedir
	go test ./...

coverage: imagedir
	@go test ./... -cover

lint:
	golangci-lint run --enable-all

fmt:
	gofmt -w *.go
