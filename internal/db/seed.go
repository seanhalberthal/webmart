package db

import (
	"context"
	"fmt"
	"github.com/seanhalberthal/webmart/internal/store"
	"log"
	"math/rand"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
	"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
	"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
	"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
	"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
	"walter", "xenia", "yasmin", "zoe",
}

var titles = []string{
	"Laptop Pro 15",
	"Wireless Headphones",
	"Smartwatch X10",
	"Gaming Keyboard",
	"4K Ultra HD Monitor",
}

var descriptions = []string{
	"A high-performance laptop with a 15-inch Retina display.",
	"Noise-canceling over-ear headphones with 40-hour battery life.",
	"A sleek smartwatch with heart rate monitoring and GPS.",
	"Mechanical keyboard with RGB lighting and customizable keys.",
	"A stunning 4K resolution monitor with HDR support.",
}

var reviews = []string{
	"Amazing laptop, very fast!", "Battery life could be better.",
	"Great sound quality, but a bit pricey.", "Very comfortable to wear!",
	"Love the fitness tracking features!", "GPS accuracy could be improved.",
	"Keys feel great for typing.", "RGB effects are stunning!",
	"Crystal-clear display, perfect for work and gaming!", "Wish it had more ports.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.UserCreate(ctx, user); err != nil {
			log.Println("Failed to create user:", err)
			return
		}
	}

	products := generateProducts(200, users)
	for _, product := range products {
		if err := store.Products.ProductCreate(ctx, product); err != nil {
			log.Println("Failed to create product:", err)
			return
		}
	}

	r := generateReviews(500, users, products)
	for _, review := range r {
		if err := store.Reviews.ReviewCreate(ctx, review); err != nil {
			log.Println("Failed to create review:", err)
			return
		}
	}

	log.Println("Seeding complete")

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@mail.com",
			Password: "123123",
		}
	}

	return users
}

func generateProducts(num int, users []*store.User) []*store.Product {
	products := make([]*store.Product, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		products[i] = &store.Product{
			UserID:      user.ID,
			Title:       titles[rand.Intn(len(titles))],
			Description: descriptions[rand.Intn(len(descriptions))],
			Rating:      0,
			Price:       0,
			Stock:       0,
		}
	}

	return products
}

func generateReviews(num int, users []*store.User, products []*store.Product) []*store.Review {
	r := make([]*store.Review, num)

	for i := 0; i < num; i++ {
		r[i] = &store.Review{
			ProductID: products[i].ID,
			UserID:    users[i].ID,
			Content:   reviews[i],
		}
	}

	return r
}
