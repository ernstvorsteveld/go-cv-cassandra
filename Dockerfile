FROM golang:1.23 AS build
WORKDIR /go/src
ADD . .
ENV CGO_ENABLED=0
RUN go build -o /go/src/cmd ./...


FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/cmd/cmd ./
COPY --from=build /go/src/cmd/config.yml ./
EXPOSE 8091/tcp
ENTRYPOINT ["./cmd"]
