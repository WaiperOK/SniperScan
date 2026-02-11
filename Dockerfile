FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /bin/sniperscan ./cmd/sniperscan

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/sniperscan /usr/local/bin/sniperscan
EXPOSE 8097
ENTRYPOINT ["/usr/local/bin/sniperscan"]
CMD ["serve", "--addr", ":8097"]
