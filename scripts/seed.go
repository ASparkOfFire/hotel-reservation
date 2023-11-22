package main

import (
	"context"
	"fmt"
	"github.com/asparkoffire/hotel-reservation/api"
	"github.com/asparkoffire/hotel-reservation/db"
	"github.com/asparkoffire/hotel-reservation/db/fixtures"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx             = context.Background()
		err             error
		mongodbEndpoint = os.Getenv("MONGO_DB_URL")
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbEndpoint))

	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)

	store := db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(&store, "James", "Foo", false)
	fmt.Println("user -> ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(&store, "Elliot", "Alderson", true)
	fmt.Println("admin -> ", api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(&store, "Some Hotel", "New York City", 5, nil)
	room := fixtures.AddRoom(&store, "large", true, 99.99, hotel.ID)
	booking := fixtures.AddBooking(&store, user.Id, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println()
	fmt.Println("booking -> ", booking.ID)
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(&store, name, location, rand.Intn(5)+1, nil)
	}
}
