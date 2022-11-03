FROM golang:1.19.2-bullseye AS dev
ENV GO111MODULE=on
WORKDIR /go/src/github.com/nubesk/binn
COPY ./ ./
RUN CGO_ENABLED=0 go build -o /go/src/github.com/nubesk/binn/build/binn main.go

FROM golang:1.19.2-alpine3.16
COPY --from=dev /go/src/github.com/nubesk/binn/build/binn /binn
CMD ["/binn"]
