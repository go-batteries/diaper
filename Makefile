fmt:
	go fmt ./...
vet:
	go vet ./...

test: vet
	go test ./...

release: fmt test
ifndef TAG
        $(error TAG is not set)
endif
	COMMIT=$(shell git log -n 1 origin/main --pretty=format:"%H")
	git tag -a $(TAG) $(COMMIT) -m "Release $(TAG)"
	git push origin $(TAG)
