FROM golang:latest AS build
WORKDIR /src
COPY go.mod main.go puzzles/ common/ .
RUN go get && go build -o /out/aoc2022

FROM scratch AS bin
COPY --from=build /out/aoc2022 /
