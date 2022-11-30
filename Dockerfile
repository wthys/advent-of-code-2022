FROM golang:latest AS build
WORKDIR /src
COPY go.mod go.sum puzzles/ common/ .
RUN go build -o /out/aoc2022

FROM scratch AS bin
COPY --from=build /out/aoc2022 /
