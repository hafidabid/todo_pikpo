package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todo_pikpo/config"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Middleware struct {
	conf config.ConfigApp
}

func (m Middleware) UnaryAuth(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authVal := md["authorization"]

	if len(authVal) == 0 {
		log.Errorf("%s please provide authorization bearer key\n", time.Now().Format("2006-01-02 15:04:05"))
		return nil, errors.New("authorization was wrong")
	}

	if len(authVal) > 0 && authVal[0] != fmt.Sprintf("Bearer %s", m.conf.EncryptKey) {
		log.Errorf("%s authorization was wrong -> %s\n", time.Now().Format("2006-01-02 15:04:05"), authVal[0])
		return nil, errors.New("authorization was wrong")
	}
	return handler(ctx, req)
}

func (m Middleware) StreamAuth(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("metadata is not provided")
	}
	authVal := md["authorization"]
	if len(authVal) == 0 || (len(authVal) > 0 && authVal[0] != fmt.Sprintf("Bearer %s", m.conf.EncryptKey)) {
		log.Errorf("%s authorization was wrong -> %s\n", time.Now().Format("2006-01-02 15:04:05"), authVal[0])
		return errors.New("authorization was wrong")
	}

	return handler(srv, stream)
}

func NewMiddleware(conf config.ConfigApp) Middleware {
	return Middleware{
		conf: conf,
	}
}
