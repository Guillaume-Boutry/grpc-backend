FROM registry.zouzland.com/face-authenticator-builder AS builder
RUN go get google.golang.org/grpc
COPY . $GOPATH/src/github.com/Guillaume-Boutry/grpc-backend
WORKDIR $GOPATH/src/github.com/Guillaume-Boutry/grpc-backend
RUN go build ./cmd/tmp-backend

FROM registry.zouzland.com/face-authenticator-runer
COPY models.txt models.txt
RUN wget -i models.txt --directory-prefix=/opt/grpc-backend && bzip2 -d $(ls /opt/grpc-backend/*.bz2)
COPY --from=builder /go/src/github.com/Guillaume-Boutry/grpc-backend/tmp-backend /opt/grpc-backend/backend

CMD ["/opt/grpc-backend/backend"]