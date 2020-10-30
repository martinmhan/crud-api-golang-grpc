# crud-api-golang-grpc

Message-based RPC API built with Go, gRPC, rabbitMQ, and MongoDB

# Summary
Super simple test API I built to do a couple things:
  - Experiment with event-driven backend architectures
  - Familiarize myself with new technologies (GoLang, gRPC, rabbitMQ)

# Notes
  - I used gRPC for the flexibility to handle tasks either synchronously or asynchronously
    - Reads are done synchronously and the requested data is returned with the gRPC response
    - Writes are done asynchronously via the message queue. The Initial gRPC call returns immediately and the write event is fired off from the producer to the queue, which the consumer then picks up to complete the task

# Pre Requisites
   - GoLang
   - gRPC
   - rabbitMQ
   - MongoDB
 
# Setup
 TBD
