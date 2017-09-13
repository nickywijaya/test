#development
dep:
	go get -u github.com/kardianos/govendor
	govendor sync

pretty:
	# gofmt -d -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -d -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	go tool vet .
