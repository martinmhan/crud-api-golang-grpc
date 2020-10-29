# crud-api-golang-grpc

Message-based RPC API built with Go, gRPC, rabbitMQ, and MongoDB

# Summary
Simple test API I built to do a couple things:
  - Experiment with event-driven backend architectures
  - Learn and implement technologies new to me (GoLang, gRPC, rabbitMQ)

# Notes
  - I used a request/response-style of communication for direct User-API interactions for the flexibility to handle tasks either synchronously or asynchronously
    - Reads are done synchronously and completed with a gRPC response containing the requested data - the message queue is not used.
    - Writes are done asynchronously via the message queue. The Initial gRPC call returns immediately with a response indicating that the event (i.e., message) has been fired off from the producer to the queue. The consumer then picks up the message from the queue and completes the task when available.
  - This test API is intentially over-simplified. More robust approaches may include:
    - An API gateway service that handles authentication and passes messages to a separate producer service
    - For read-heavy applications, a read-optimized view service that handles Reads without having to query the database directly as often
    - Separate queues for reads/replies if we want reads to follow event driven architecture in addition to writes
  - Also, for simplicity's sake, I built both producer/consumer services in one monorepo. This did prove useful in consolidating database code though, which both services use.

# Pre Requisites
   - GoLang
   - gRPC
   - rabbitMQ
   - MongoDB
 
# Setup
 TBD
