package tokens_storage

import (
	"context"
	"fmt"
	"log"

	config "github.com/rlapenok/wallets/backend/main_server/internal/configs/storages/tokens_storage_config"
	"github.com/rlapenok/wallets/backend/main_server/internal/token_library"
	mongo_driver "go.mongodb.org/mongo-driver/mongo"
	mongo_opt "go.mongodb.org/mongo-driver/mongo/options"
)

type Crud interface {
	SaveToken(string, context.Context) error
	DeleteToken(string, context.Context) error
	CheckToken(string, context.Context) error
}

type TokensStorage struct {
	coll *mongo_driver.Collection
}

// Create new Storage for saving refresh tokens
func New() *TokensStorage {
	//Create config for DataBase
	config := config.New()
	//Get uri for connection
	uri := config.GetUri()
	//Set auth creditional
	auth := mongo_opt.Credential{Username: config.Login, Password: config.Password, AuthSource: config.DbName}
	//Set client opt
	opt := mongo_opt.Client().ApplyURI(uri).SetAuth(auth)
	//Connect to MongoDb
	client, err := mongo_driver.Connect(context.Background(), opt)
	if err != nil {
		//Add normal loggining
		log.Fatal(err)
	}
	ping_err := client.Ping(context.Background(), nil)
	if ping_err != nil {
		//Add normal loggining
		log.Fatal(ping_err)
	}
	fmt.Println("Connected to TokenStorage")
	collection := client.Database(config.DbName).Collection(config.CollName)
	return &TokensStorage{coll: collection}

}
func (storage *TokensStorage) SaveToken(hash string, ctx context.Context) error {
	token := token_library.HashedRefreshToken{Hash: hash}
	_, err := storage.coll.InsertOne(ctx, token)
	if err != nil {
		return err
	}
	return nil

}
func (storage *TokensStorage) DeleteToken(hash string, ctx context.Context) error {
	token := token_library.HashedRefreshToken{Hash: hash}
	_, err := storage.coll.DeleteOne(ctx, token)
	if err != nil {
		return err
	}
	return nil

}
func (storage *TokensStorage) CheckToken(hash string, ctx context.Context) error {
	var value *token_library.HashedRefreshToken
	token := token_library.HashedRefreshToken{Hash: hash}
	err := storage.coll.FindOne(ctx, token).Decode(value)
	if err != nil {
		return err
	}
	return nil

}
