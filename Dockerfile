FROM scratch

ADD ./bin/go-webserver-dist /go-webserver

#ENV ADDR="8080"

ENTRYPOINT ["/go-webserver"]

