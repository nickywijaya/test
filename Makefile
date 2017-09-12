#development
dep:
	go get -u github.com/kardianos/govendor
	govendor sync

pretty:
	gofmt -w **/**/*.go
	gofmt -w **/*.go
	go tool vet .
