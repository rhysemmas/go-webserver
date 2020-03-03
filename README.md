# go-webserver

A small webserver written in Go which will alternate returning HTTP 200 and 500 status codes, to be deployed to Kube

### Envoy

There is a `docker-compose.yml` and example envoy config (in the `envoy/` dir) for testing the app locally with envoyproxy. The app will bind to port 8080 on your machine so it can be reached directly. Envoy binds to port 10000 on your machine, so you can test the app via envoy as well.

`Note:` Replace the service IP address (0.0.0.0) with the private IP address of your machine, because docker networking
