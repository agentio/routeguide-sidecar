# routeguide-sidecar

This directory contains a command-line tool that includes a server and client for a sample gRPC service that demonstrates all four streaming modes of gRPC. The service implements the same Route Guide API that is used in the grpc-go [route_guide](https://github.com/grpc/grpc-go/tree/master/examples/route_guide) example.

The server and client can communicate over TCP ports or Linux abstract sockets.

To use TCP ports, start the server with the `--port` flag:
```sh
routeguide-sidecar serve --port 8088
```

Call the client with the `--address` flag as below:
```sh
routeguide-sidecar client --address localhost:8088
```

To use Linux abstract sockets, run the server with the `--socket` flag:
```sh
routeguide-sidecar serve --socket @routeguide
```

Call the client with the `--address` flag as below:
```sh
routeguide-sidecar client --address unix:@routeguide
```
