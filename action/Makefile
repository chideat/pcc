
MODELS = models
BINARY = action
CLEAN  = clean

TARGETS = $(BINARY) 

.PHONY: all $(TARGETS)
all: $(TARGETS)

$(BINARY):
	go build -ldflags "-X main.BuildTimestamp=`date +%Y%m%d.%H%M%S` -X main.BuildCommit=`git rev-parse HEAD`" -o $(BINARY)

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.BuildTimestamp=`date +%Y%m%d.%H%M%S` -X main.BuildCommit=`git rev-parse HEAD`"

$(MODELS):
	make -C proto/

$(CLEAN):
	go clean
