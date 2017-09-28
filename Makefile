.PHONY: all

REGISTRY  = registry.bukalapak.io/bukalapak
DDIR      = deploy
ODIR      = $(DDIR)/_output
NOCACHE   = --no-cache
VERSION   = $(shell git show -q --format=%h)
SERVICES ?= http-go-xample background-go-xample cron-go-xample
ENV      ?= default

dep:
	go get -u github.com/kardianos/govendor
	govendor sync

pretty:
	# gofmt -d -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -d -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	go tool vet .

coverage:
	@./coverage.sh
	go tool cover -html=coverage.out -o coverage.html

all:
	consul compile build push deployment

test:
	govendor test -v +local,^program

compile:
	@$(foreach var, $(SERVICES), GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(ODIR)/$(var)/bin/$(var) app/$(var)/main.go;)

$(ODIR):
	mkdir -p $(ODIR)

consul: $(ODIR)
	@wget https://releases.hashicorp.com/envconsul/0.6.2/envconsul_0.6.2_linux_amd64.tgz
	@tar -xf envconsul_0.6.2_linux_amd64.tgz -C $(ODIR)/
	@rm envconsul_0.6.2_linux_amd64.tgz

build:
	@$(foreach var, $(SERVICES), docker build $(NOCACHE) -t $(REGISTRY)/go-xample/$(var):$(VERSION) -f ./deploy/$(var)/Dockerfile .;)

push:
	@$(foreach var, $(SERVICES), docker push $(REGISTRY)/go-xample/$(var):$(VERSION);)

deployment: $(ODIR)
ifeq ($(ENV),default)
	kubelize deployment -v $(VERSION) $(SERVICES)
else
	kubelize deployment -e $(ENV) -v $(VERSION) $(SERVICES)
endif

$(ENV):
	@$(foreach var, $(SERVICES), kubectl replace -f $(ODIR)/$(var)/$@/deployment.yml;)

setup:
	docker run --rm -it --network host -v $PWD/db:/app/db -v $PWD/.env:/app/.env registry.bukalapak.io/sre/migration:0.0.1 db:create
	docker run --rm -it --network host -v $PWD/db:/app/db -v $PWD/.env:/app/.env registry.bukalapak.io/sre/migration:0.0.1 db:migrate

migrate:
	docker run --rm -it --network host -v $PWD/db:/app/db -v $PWD/.env:/app/.env registry.bukalapak.io/sre/migration:0.0.1 db:migrate
