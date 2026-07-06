BIN   := ai-chat
GO    := go
DEST  ?= /usr/local/bin

.PHONY: build install uninstall clean

build:
	$(GO) build -o $(BIN) .

install: build
	install -d "$(DEST)"
	install -m 755 "$(BIN)" "$(DEST)/$(BIN)"
	@echo "Installed $(BIN) to $(DEST)/$(BIN)"

uninstall:
	rm -f "$(DEST)/$(BIN)"

clean:
	rm -f "$(BIN)"
