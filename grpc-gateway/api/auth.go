package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Kyejoon-Lee/grpc-gateway/config"
	"github.com/Kyejoon-Lee/grpc-gateway/ent/proto/entpb"
	"github.com/Kyejoon-Lee/grpc-gateway/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

type testModel struct {
	IdToken string `json:"id_token"`
}

func Login(c *gin.Context) {
	cfg := config.GetConfig()
	ex := c.Query("code")
	query := url.Values{}
	query.Add("grant_type", "authorization_code")
	query.Add("code", ex)
	query.Add("redirect_uri", "http://127.0.0.1:9091/login")
	query.Add("client_id", cfg.ClientID)
	query.Add("client_secret", cfg.ClientSecret)
	base, err := url.Parse("https://kauth.kakao.com/oauth/token")
	if err != nil {
		log.Println(err)
	}
	base.RawQuery = query.Encode()
	b := bytes.NewBufferString("")
	resp, err := http.Post(base.String(), "application/x-www-form-urlencoded", b)
	if err != nil {
		log.Println(err)
	}
	mod := testModel{}
	sendByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	json.Unmarshal(sendByte, &mod)
	parsed, err := jwt.Parse(mod.IdToken, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	log.Println(parsed.Claims.(jwt.MapClaims)["email"])
	temp := entpb.User{Email: "ll"}
	grpcClient := service.GrpcCLI{}
	grpcClient.StartGrpcConnection()

	resp1, err := grpcClient.UserServiceClientCLI.Create(context.Background(), &entpb.CreateUserRequest{User: &temp})
	if err != nil {
		log.Error(err)
	}
	log.Info(resp1)
	defer grpcClient.ShutdownGrpcConnection()
	c.Next()

}
