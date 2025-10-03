package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
	stopChan   chan bool
	wg         *sync.WaitGroup
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
		stopChan:   make(chan bool),
		wg:         &sync.WaitGroup{},
	}
	
	repo.wg.Add(1)
	go repo.startAuctionCloser()
	
	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func getAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		return time.Minute * 5
	}
	return duration
}

func (ar *AuctionRepository) startAuctionCloser() {
	defer ar.wg.Done()
	
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	
	for {
		select {
		case <-ar.stopChan:
			return
		case <-ticker.C:
			ar.closeExpiredAuctions()
		}
	}
}

func (ar *AuctionRepository) closeExpiredAuctions() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	
	now := time.Now()
	auctionDuration := getAuctionDuration()
	expirationTime := now.Add(-auctionDuration)
	
	filter := bson.M{
		"status":    auction_entity.Active,
		"timestamp": bson.M{"$lt": expirationTime.Unix()},
	}
	
	update := bson.M{
		"$set": bson.M{
			"status": auction_entity.Completed,
		},
	}
	
	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error closing expired auctions", err)
		return
	}
	
	if result.ModifiedCount > 0 {
		logger.Info("Closed expired auctions", zap.Int64("count", result.ModifiedCount))
	}
}

func (ar *AuctionRepository) Stop() {
	close(ar.stopChan)
	ar.wg.Wait()
}
