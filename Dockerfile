FROM golang:alpine as build
COPY . /go/src/kohlbau.de/x/go-urlredirect-gitlab 
WORKDIR /go/src/kohlbau.de/x/go-urlredirect-gitlab 
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o go-urlredirect-gitlab main.go

FROM scratch
COPY --from=build /go/src/kohlbau.de/x/go-urlredirect-gitlab/go-urlredirect-gitlab /bin/
ENTRYPOINT ["/bin/go-urlredirect-gitlab"]
