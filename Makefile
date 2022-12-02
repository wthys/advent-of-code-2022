NAME:=aoc2022
BIN_DIR:=./bin

PROG:=$(BIN_DIR)/$(NAME)

SOURCE:=$(wildcard src/*.go)

NOWDATE:=$(shell TZ="EST" date +%Y%m%d)
NOWDAY:=$(shell TZ="EST" date +%d)
ENDDATE:=20221225

build: $(PROG)
.PHONY: build

$(PROG): src/go.mod $(SOURCE)
	docker build --target bin --output $(BIN_DIR)/ . 
	touch $(PROG)

run: $(PROG)
	@$(PROG) input $(DAY) | $(PROG) run $(DAY)

run-all: $(PROG)
	@if test "$(NOWDATE)" -lt "$(ENDDATE)"; then for day in `seq $(NOWDAY)`; do $(PROG) input $$day | $(PROG) run $$day; done; else for day in `seq 25`; do $(PROG) input $$day | $(PROG) run $$day;done;fi

clean:
	rm -f $(PROG)
