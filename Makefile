NAME:=aoc2022
BIN_DIR:=./bin

PROG:=$(BIN_DIR)/$(NAME)

GOFILES:=$(shell find src/ -type f -name "*.go")

NOWDATE:=$(shell TZ="EST" date +%Y%m%d)
NOWDAY:=$(shell TZ="EST" date +%d)
ENDDATE:=20221225
DOCKERRUN=docker run --rm -i --env AOC_SESSION ${AOC_RUNOPTS} aoc2022:latest $(ELAPSEDOPTS)
ifdef ELAPSED
ELAPSEDOPTS:=-e
endif

.PHONY: build run run-all clean example build-run run-bare example-bare

build: $(PROG)


$(PROG): src/go.mod $(GOFILES)
	DOCKER_BUILDKIT=1 docker build --target bin --output $(BIN_DIR)/ . 
	touch $(PROG)

build-run: $(PROG)
	docker build -f Dockerfile.run -t aoc2022:latest .

run: build-run $(PROG)
	@$(PROG) input $(DAY) | $(DOCKERRUN) $(DAY)

run-bare: $(PROG)
	@$(PROG) input $(DAY) | $(PROG) run $(DAY)

run-all: $(PROG)
	@if test "$(NOWDATE)" -lt "$(ENDDATE)"; then for day in `seq $(NOWDAY)`; do $(PROG) input $$day | $(DOCKERRUN) $$day; done; else for day in `seq 25`; do $(PROG) input $$day | $(DOCKERRUN) $$day;done;fi


clean:
	rm -f $(PROG)

example: $(PROG) build-run
	@cat examples/day$(DAY).txt | $(DOCKERRUN) $(DAY)

example-bare: $(PROG)
	@cat examples/day$(DAY).txt | $(PROG) run $(DAY)
