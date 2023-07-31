package service

import (
	"Social-Net-Dialogs/internal/helper"
	"Social-Net-Dialogs/internal/store"
	"Social-Net-Dialogs/models"
	"Social-Net-Dialogs/pkg/tokenservice"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Client struct {
	Client       tokenservice.ValidateTokenClient
	SessionStore store.SessionStore
}

func (c *Client) ValidateToken(ctx context.Context, token string) (*tokenservice.ValidateTokenResponse, error) {
	var resp tokenservice.ValidateTokenResponse
	session, err := c.SessionStore.LoadSession(token)
	if err != nil {
		return nil, err
	}
	if len(session.UserId) > 0 {
		resp = tokenservice.ValidateTokenResponse{
			Token:     session.Token,
			UserId:    session.UserId,
			Id:        session.Id,
			CreatedAt: session.CreatedAt,
		}
		return &resp, nil
	}

	req := tokenservice.ValidateTokenRequest{
		Token: token,
	}

	validateToken, err := c.Client.ValidateToken(ctx, &req)
	if err != nil {
		return nil, err
	}
	go func() {
		_, err := c.SessionStore.CreateSession(&models.UserSession{
			UserId:    validateToken.UserId,
			Id:        validateToken.Id,
			Token:     validateToken.Token,
			CreatedAt: validateToken.CreatedAt,
		})
		if err != nil {
			log.Println("Cannot write session to cache")
		}
	}()

	return validateToken, nil
}

func NewTokenServiceClient(store store.SessionStore) *Client {
	host := helper.GetEnvValue("RPC_SERVER_HOST", "localhost")
	port := helper.GetEnvValue("RPC_SERVER_PORT", "50051")

	cwt, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(cwt, host+":"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err)
	}
	//defer conn.Close()

	ts := tokenservice.NewValidateTokenClient(conn)

	return &Client{
		Client:       ts,
		SessionStore: store,
	}
}
