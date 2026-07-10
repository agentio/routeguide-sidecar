FROM golang:1.26.4 AS builder
WORKDIR /app
COPY . ./
RUN apt-get update
RUN apt-get install unzip
RUN ./tools/fetch-protoc.sh
ENV PATH="/root/local/bin:${PATH}"
RUN make rpc
RUN CGO_ENABLED=0 GOOS=linux go build -v -o routeguide-sidecar .

FROM gcr.io/distroless/static-debian13
COPY --from=builder /app/routeguide-sidecar /usr/local/bin/routeguide-sidecar
CMD ["/usr/local/bin/routeguide-sidecar", "serve", "-p", "8080"]
