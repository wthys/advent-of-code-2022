NAME=aoc2022
BIN_DIR=./bin

build: compile-aoc
.PHONY: build

compile-aoc: bin/${NAME}
	@docker build --target bin --output ${BIN_DIR}/ . 
