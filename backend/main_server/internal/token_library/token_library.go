package token_library

type HashedRefreshToken struct {
	Hash string `bson:"hash"`
}
