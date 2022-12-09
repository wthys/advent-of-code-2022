FROM golang:latest AS build
WORKDIR /src
COPY src /src
RUN go get
RUN go test ./... && go build -o /out/aoc2022

FROM scratch AS bin
COPY --from=build /out/aoc2022 /
