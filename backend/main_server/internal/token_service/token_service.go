package token_service

import (
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rlapenok/wallets/backend/main_server/internal/token_library"
	"golang.org/x/crypto/bcrypt"
)

type rowTokensService interface {
	createRowRefreshToken(string, time.Time) (*token_library.RowRefreshToken, error)
	createRowAccessToken(string, time.Time) *jwt.Token
	signToken(*jwt.Token, string) (*string, error)
	checkSignature(string, string) (bool, error)
}

type TokenServiceAuth interface {
	rowTokensService
	CreateTokens(string, time.Time) (*token_library.Tokens, error)
	CheckTokens(string, string) (bool, error)
}

type TokenService struct {
	refresh_lft int
	access_lft  int
	secret_key  string
}

func NewForTest() *TokenService {
	service := TokenService{refresh_lft: 2, access_lft: 1, secret_key: "123"}
	return &service
}

func New() *TokenService {
	r_lft := os.Getenv("refresh_lft")
	if r_lft == "" {
		//Add noraml loggining
		log.Fatal("Set lifetime refresh token in Dockerfile")
	}
	refresh_lft, err := strconv.Atoi(r_lft)
	if err != nil {
		log.Fatal("lifetime refresh token in Dockerfile != int")
	}
	a_lft := os.Getenv("access_lft")
	if a_lft == "" {
		log.Fatal("Set lifetime access token in Dockerfile")
	}
	access_lft, err := strconv.Atoi(a_lft)
	if err != nil {
		log.Fatal("lifetime access token in Dockerfile != int")

	}
	secret_key := os.Getenv("secret_key")
	if secret_key == "" {
		log.Fatal("Set secret_key in Dockerfile")
	}
	token_manager := TokenService{refresh_lft: refresh_lft, access_lft: access_lft, secret_key: secret_key}
	return &token_manager

}
func (service *TokenService) createRowRefreshToken(guid string, exp time.Time) (*token_library.RowRefreshToken, error) {
	//Encode guid into base64
	encoded_guid := base64.StdEncoding.EncodeToString([]byte(guid))
	//Create row_refresh_token
	row_token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": encoded_guid,
		"exp":  exp.Add(time.Minute * time.Duration(service.refresh_lft)).Unix(),
	})
	sign_for_access_token, err := row_token.SigningString()
	if err != nil {
		//Add loggining
		return nil, err

	}
	row_refresh_token := token_library.RowRefreshToken{Token: row_token, SignForAccessToken: sign_for_access_token}
	return &row_refresh_token, nil

}
func (service *TokenService) createRowAccessToken(guid string, exp time.Time) *jwt.Token {
	row_token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   guid,
		ExpiresAt: jwt.NewNumericDate(exp.Add(time.Minute * time.Duration(service.access_lft))),
	})
	return row_token
}

func (service *TokenService) signToken(row_token *jwt.Token, sign string) (*string, error) {
	token, err := row_token.SignedString([]byte(service.secret_key + sign))
	if err != nil {
		return nil, err
	}
	return &token, nil
}
func (service *TokenService) checkSignature(token string, signature string) (bool, error) {
	result, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrHashUnavailable
		}
		return []byte(service.secret_key + signature), nil
	})
	if err != nil {
		return false, err
	}
	if !result.Valid {
		return false, nil
	}
	return true, nil
}
func (service *TokenService) CreateTokens(guid string, exp time.Time) (*token_library.Tokens, error) {
	//Create row refresh token
	row_refresh_token, err1 := service.createRowRefreshToken(guid, exp)
	if err1 != nil {
		//add loggining
		return nil, err1
	}
	//Create row access token
	row_access_token := service.createRowAccessToken(guid, exp)
	//Create access token
	access_token, err2 := service.signToken(row_access_token, row_refresh_token.SignForAccessToken)
	if err2 != nil {
		//add loggining
		return nil, err2

	}
	//Create refresh token
	refresh_token, err3 := service.signToken(row_refresh_token.Token, *access_token)
	if err3 != nil {

		//add loggining
		return nil, err3
	}
	//Convert refresh token into base64
	b64_refresh_token := base64.StdEncoding.EncodeToString([]byte(*refresh_token))
	//Create bycrypt hash refresh token
	hashed_refresh_token, err4 := bcrypt.GenerateFromPassword([]byte(*refresh_token), bcrypt.DefaultCost)
	if err4 != nil {
		//add loggining
		return nil, err4
	}
	//bycrypt hash refresh token convert into string
	string_hashed_refresh_token := string(hashed_refresh_token)
	tokens := token_library.Tokens{
		Access:         *access_token,
		UserRefresh:    b64_refresh_token,
		StorageRefresh: string_hashed_refresh_token,
	}
	return &tokens, nil

}
func (service *TokenService) CheckTokens(access string, refresh string) (bool, error) {

	components_refresh_token := strings.Split(refresh, ".")
	signature_for_access_token := components_refresh_token[0] + "." + components_refresh_token[1]
	result_access_token, err1 := service.checkSignature(access, signature_for_access_token)
	if err1 != nil {
		return result_access_token, nil
	}
	result_refresh_token, err2 := service.checkSignature(refresh, access)
	if err2 != nil {
		return result_refresh_token, nil
	}
	switch {
	case result_access_token && result_refresh_token:
		return true, nil
	default:
		return false, nil
	}

}
