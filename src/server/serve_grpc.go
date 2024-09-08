package server

// func ServeGrpc(ctx context.Context, addr string) (err error) {
// 	defer log.Println("GRPC server stopped", err)

// 	lis, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		return err
// 	}
// 	defer lis.Close()

// 	server := grpc.NewServer()
// 	authenticator.RegisterAuthenticatorServiceServer(server, &authenticatorGrpc.Server{})

// 	log.Printf("Listen and Serve Auth-Grpc-Service API at: %s\n", addr)
// 	return server.Serve(lis)
// }
