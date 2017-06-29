RELEASE="release"
MAIN_SRC=misp-cli

all:
	$(MAKE) clean
	$(MAKE) init
	$(MAKE) compile

init:
	mkdir -p $(RELEASE)
	mkdir -p $(RELEASE)/linux
	mkdir -p $(RELEASE)/windows
	mkdir -p $(RELEASE)/darwin

compile:linux windows darwin

linux:
	GOARCH=386 GOOS=linux go build -o $(RELEASE)/linux/$(MAIN_SRC)-386 $(MAIN_SRC).go
	GOARCH=amd64 GOOS=linux go build -o $(RELEASE)/linux/$(MAIN_SRC)-amd64 $(MAIN_SRC).go

windows:
	GOARCH=386 GOOS=windows go build -o $(RELEASE)/windows/$(MAIN_SRC)-386.exe $(MAIN_SRC).go
	GOARCH=amd64 GOOS=windows go build -o $(RELEASE)/windows/$(MAIN_SRC)-amd64.exe $(MAIN_SRC).go

darwin:
	GOARCH=386 GOOS=darwin go build -o $(RELEASE)/darwin/$(MAIN_SRC)-386 $(MAIN_SRC).go
	GOARCH=amd64 GOOS=darwin go build -o $(RELEASE)/darwin/$(MAIN_SRC)-amd64 $(MAIN_SRC).go

clean:
	rm -rf $(RELEASE)/*
