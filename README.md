# Logistics Engine API

**What is Our  Example Project**

This is example is of a repository that hosts the source code for an API server capable of handling requests, as specified by the client's operations, encoded using protocol buffers.  It is intentionally divided client and server source-code between two different repositories, to give a sense of interaction two real distributed systems.
You can find the complete client  example [here](https://github.com/ivanbulyk/clients_logistics_engine_api).

So to see everything in action, you can test it by running our server first:

```text
$ go run ./cmd/logistics/ main.go
```
or just

```text
$ make 
```

Then in other terminal, we run our client, in project root:

```text
$ go run ./cmd/logistics/ main.go
```
or just

```text
$ make 
```

The server should wait infinitely, emitting logs on calls, and the client should be returning without any error on the terminal. Then you want to hit the localhost:50051 LogisticsEngineAPI/MetricsReport, with any gRPC client to see the calculations result.