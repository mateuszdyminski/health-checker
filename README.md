# Health Status Checker

Tool which periodically checks response time from any specified address.

# Setup

## Requirements

- GOLANG see: https://golang.org/doc/install

## Build

```
./dev.sh build
```

## Run

```
bin/app --log_dir="/tmp/" --alsologtostderr --stderrthreshold=INFO --dir="static" --host="localhost" --port="8090" --address="http://google.com"
```
