MAIN=cmd/app/main.go
BUILD_DIR =build
O_FILE=server

.PHONY: all build run clear

all: build

build: 
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(O_FILE) $(MAIN)

run: build
	./$(BUILD_DIR)/$(O_FILE)

clear:
	rm -rf $(BUILD_DIR)
# 

update_mod:
	go mod tidy