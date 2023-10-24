package token_library

import "github.com/golang-jwt/jwt/v5"

type HashedRefreshToken struct {
	Hash string `bson:"hash"`
}
type RowRefreshToken struct {
	Token              *jwt.Token
	SignForAccessToken string
}
type Tokens struct {
	Access         string
	UserRefresh    string
	StorageRefresh string
}
