package token_service_test

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/rlapenok/wallets/backend/main_server/internal/token_service"
)

const GUID1 string = "605384e0-1b1a-47e6-8eef-79ff5cab61ca"
const GUID2 string = "605384e0-1b1a-48e6-8eef-79ff5cab61ca"

func TestCreateTokens(t *testing.T) {
	service := token_service.NewForTest()
	tokens, err := service.CreateTokens(GUID1, time.Now())
	if err != nil {
		t.Error(err)
	}
	if tokens == nil {
		t.Errorf("No tokens created")
	}

}
func TestCheckTokens1(t *testing.T) {
	service := token_service.NewForTest()
	tokens, err1 := service.CreateTokens(GUID1, time.Now())
	if err1 != nil {
		t.Error(err1)
	}
	if tokens == nil {
		t.Errorf("No tokens created")
		return
	}
	encode_refresh_token, _ := base64.StdEncoding.DecodeString(tokens.UserRefresh)

	result, err2 := service.CheckTokens(tokens.Access, string(encode_refresh_token))
	if err2 != nil {
		t.Error(err2)

	}
	if !result {
		t.Errorf("Tokens were not created by this service")
	}

}
func TestCheckTokens2(t *testing.T) {
	service := token_service.NewForTest()
	tokens1, _ := service.CreateTokens(GUID1, time.Now())
	tokens2, _ := service.CreateTokens(GUID2, time.Now())
	encode_refresh_token1, _ := base64.StdEncoding.DecodeString(tokens1.UserRefresh)
	result1, _ := service.CheckTokens(tokens2.Access, string(encode_refresh_token1))
	if result1 {
		t.Errorf("False was expected, result received %v", result1)
	}
	encode_refresh_token2, _ := base64.StdEncoding.DecodeString(tokens2.UserRefresh)
	result2, _ := service.CheckTokens(tokens1.Access, string(encode_refresh_token2))
	if result2 {
		t.Errorf("False was expected, result received %v", result1)
	}

}
