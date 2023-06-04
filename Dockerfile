FROM golang AS build

WORKDIR /go/src/app

COPY . .


RUN go mod download

RUN apt-get update
RUN apt install -y protobuf-compiler
RUN apt install -y golang-goprotobuf-dev

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

RUN protoc --version

RUN chmod +x generate_proto.sh
RUN ./generate_proto.sh

RUN GOOS=linux go build -ldflags="-s -w" -o ./server

FROM golang
WORKDIR /usr/bin
COPY --from=build /go/src/app/server .
COPY --from=build /go/src/app/.env .

RUN chmod +x server

EXPOSE 1000-15000

ENTRYPOINT ["./server"]


