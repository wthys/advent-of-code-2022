NAME=aoc2022
BIN_DIR=./bin

PROG=$(BIN_DIR)/$(NAME)

SOURCE=$(wildcard src/*.go)

build: $(PROG)
.PHONY: build

$(PROG): src/go.mod $(SOURCE)
	docker build --target bin --output $(BIN_DIR)/ . 
	touch $(PROG)

run: $(PROG)
	@$(PROG) input $(DAY) | $(PROG) run $(DAY)


clean:
	rm -f $(PROG)
