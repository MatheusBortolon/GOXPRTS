package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuctionAutoClose(t *testing.T) {
	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		t.Skip("MONGODB_URL is not set; skipping integration test")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		_ = client.Disconnect(ctx)
	}()

	dbName := fmt.Sprintf("auctions_test_%d", time.Now().UnixNano())
	database := client.Database(dbName)
	defer func() {
		_ = database.Drop(ctx)
	}()

	originalInterval := os.Getenv("AUCTION_INTERVAL")
	if err := os.Setenv("AUCTION_INTERVAL", "200ms"); err != nil {
		t.Fatalf("failed to set AUCTION_INTERVAL: %v", err)
	}
	defer func() {
		if originalInterval == "" {
			_ = os.Unsetenv("AUCTION_INTERVAL")
			return
		}
		_ = os.Setenv("AUCTION_INTERVAL", originalInterval)
	}()

	repo := NewAuctionRepository(database)
	auction, createErr := auction_entity.CreateAuction(
		"Phone",
		"Electronics",
		"New phone with box and accessories",
		auction_entity.New,
	)
	if createErr != nil {
		t.Fatalf("failed to create auction entity: %v", createErr)
	}

	if err := repo.CreateAuction(ctx, auction); err != nil {
		t.Fatalf("failed to create auction: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		stored, err := repo.FindAuctionById(ctx, auction.Id)
		if err != nil {
			t.Fatalf("failed to find auction: %v", err)
		}
		if stored.Status == auction_entity.Completed {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	t.Fatalf("expected auction to be closed automatically")
}
