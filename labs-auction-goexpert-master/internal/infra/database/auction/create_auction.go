package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
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

	ar.triggerAuctionAutoClose(auctionEntity.Id, auctionEntity.Timestamp)

	return nil
}

func (ar *AuctionRepository) triggerAuctionAutoClose(auctionId string, startTime time.Time) {
	endTime := calculateAuctionEndTime(startTime)
	waitDuration := time.Until(endTime)
	if waitDuration < 0 {
		waitDuration = 0
	}

	go func() {
		timer := time.NewTimer(waitDuration)
		defer timer.Stop()
		<-timer.C

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ar.closeAuctionIfActive(ctx, auctionId)
	}()
}

func (ar *AuctionRepository) closeAuctionIfActive(
	ctx context.Context,
	auctionId string) *internal_error.InternalError {
	filter := bson.M{"_id": auctionId, "status": auction_entity.Active}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	if _, err := ar.Collection.UpdateOne(ctx, filter, update); err != nil {
		logger.Error("Error trying to close auction", err)
		return internal_error.NewInternalServerError("Error trying to close auction")
	}

	return nil
}

func calculateAuctionEndTime(startTime time.Time) time.Time {
	return startTime.Add(getAuctionInterval())
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}
