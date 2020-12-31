# GRPC-Backend

This project is the point where all grpc connections arrive !

## Build
A docker file is present to build it !

Just run:
```
docker build . -t registry.zouzland.com/grpc-backend:snapshot
```

## Configuration

The different configuration are the following:

```bash
PORT=9000
K_SINK_ENROLL=#no default but required
K_SINK_AUTHENTICATE=#no default but required
```