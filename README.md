## Chameleon

A prototype Go application that masquerades as a Habitat Builder worker

Currently, it generates and sends a "Ready" heartbeat to a running job server every 30 seconds.

The job server address can be specified as a command line parameter, or it will default to localhost:5567 if not specified (which should be fine for a test env). The heartbeat is sent as a Protobuf-encoded message over a ZMQ Pub socket.

Example:
```
./chameleon
Connecting to localhost:5567
Waiting for connection to complete
Sending heartbeat endpoint:"12345@chameleon" os:Linux state:Ready
```

On the running job server (assuming it is set to log debug output):
```
29T02:45:09Z: habitat_builder_jobsrv::server::worker_manager: Got heartbeat: endpoint: "12345@chameleon" os: Linux state: Ready
```
