APP=thing-greenhouse-iot
BIN=$(PWD)/$(APP)

GO ?= go

pi: clean
	@echo "[pi] Building..."
	@cd cmd && @GOOS=linux GOARM=7 GOARCH=arm $(GO) build -o $(BIN)\

run r: build
	@echo "[run] Running..."
	@$(BIN)

clean:
	@echo "[clean] Removing $(BIN)..."
	@rm $(BIN)

.PHONY: pi run clean 
