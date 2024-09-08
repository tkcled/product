package grpc

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"sync"

// 	service_account "hshelby-tkcled-product/src/service/account"
// 	service_shared "hshelby-tkcled-product/src/service/shared"

// 	"github.com/mex-gf-tienhung/proto/golang/authenticator"
// )

// type Server struct {
// 	authenticator.AuthenticatorServiceServer
// }

// func (s *Server) AccountCreate(ctx context.Context, args *authenticator.AccountCreateRequest) (*authenticator.AccountCreateResponse, error) {
// 	if args == nil {
// 		return nil, fmt.Errorf("bad request: nil data")
// 	}

// 	numberThread := 100
// 	var channel = make(chan service_account.AccountAddCommand, numberThread+1)
// 	var waitGroup sync.WaitGroup

// 	// Number of concurency goroutine at a time.
// 	waitGroup.Add(numberThread)

// 	for thread := 0; thread < numberThread; thread++ {
// 		go func() {
// 			for {
// 				// Close goroutine and minus wait group current thread
// 				input, isStill := <-channel
// 				if !isStill {
// 					waitGroup.Done()
// 					return
// 				}

// 				_, err := service_account.AccountAdd(ctx, &input)
// 				if err != nil {
// 					log.Println(err)
// 				}
// 			}
// 		}()
// 	}

// 	for _, ele := range args.ListAccount {
// 		data := service_account.AccountAddCommand{
// 			Email:    ele.Email,
// 			Password: ele.Password,
// 			Role:     ele.Role,
// 			StaffID:  ele.StaffId,
// 		}

// 		channel <- data
// 	}
// 	return &authenticator.AccountCreateResponse{}, nil
// }

// func (s *Server) TokenVerify(ctx context.Context, args *authenticator.TokenVerifyRequest) (*authenticator.TokenVerifyResponse, error) {
// 	if args == nil {
// 		return nil, fmt.Errorf("bad request: nil data")
// 	}

// 	data := service_shared.TokenVerifyCommand{
// 		Token:      args.JwtToken,
// 		TypeVerify: service_shared.TypeVerifyBasic,
// 	}

// 	info, err := service_shared.TokenVerify(ctx, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &authenticator.TokenVerifyResponse{
// 		Email:     info.Email,
// 		AccountId: info.AccountID,
// 		Role:      info.Role,
// 		Status:    int32(info.Status),
// 	}, nil
// }

// func (s *Server) AccountBlock(ctx context.Context, args *authenticator.AccountBlockRequest) (*authenticator.AccountBlockResponse, error) {
// 	if args == nil {
// 		return nil, fmt.Errorf("bad request: nil data")
// 	}

// 	data := service_account.AccountBlockByStaffIDCommand{
// 		StaffID: args.StaffId,
// 	}

// 	_, err := service_account.AccountBlockByStaffID(ctx, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &authenticator.AccountBlockResponse{}, nil
// }

// func (s *Server) AccountUnblock(ctx context.Context, args *authenticator.AccountUnblockRequest) (*authenticator.AccountUnblockRespone, error) {
// 	if args == nil {
// 		return nil, fmt.Errorf("bad request: nil data")
// 	}

// 	data := service_account.AccountUnblockByStaffIDCommand{
// 		StaffID: args.StaffId,
// 	}

// 	_, err := service_account.AccountUnblockByStaffID(ctx, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &authenticator.AccountUnblockRespone{}, nil
// }
