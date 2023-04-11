package main

import (
	"context"
	server "github.com/adsrkey/grpc-company-info-parser/internal/delivery/grpc"
	"github.com/adsrkey/grpc-company-info-parser/parser/parserpb"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	srv := &server.Server{}
	go runGatewayServer(srv)
	runGrpcServer(srv)
}

func runGrpcServer(srv *server.Server) {
	s := grpc.NewServer()
	parserpb.RegisterParserServiceServer(s, srv)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal("cannot create gRPC listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGatewayServer(srv *server.Server) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux := runtime.NewServeMux()

	err := parserpb.RegisterParserServiceHandlerServer(ctx, grpcMux, srv)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	r := mux.NewRouter()
	r.Handle("/api/v1/company/inn/{inn}/info", grpcMux)

	swaggerHandler := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./doc/swagger")))
	r.Handle("/swaggerui/", swaggerHandler)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, r)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server")
	}
}
