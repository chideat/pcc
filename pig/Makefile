DOCKER=$(shell docker --version 2>&1  > /dev/null)

MODELS = models
BINARY = pig
CLEAN  = clean

TARGETS = $(BINARY) 

.PHONY: all $(TARGETS)
all: $(TARGETS)
ifdef DOCKER
	docker build -t "pcc/pig:`date +%Y%m%d%H%M`" .
endif

$(BINARY):
	go build -ldflags "-X main.BuildTimestamp=`date +%Y%m%d.%H%M%S` -X main.BuildCommit=`git rev-parse HEAD`" -o $(BINARY)

$(CLEAN):
	go clean
