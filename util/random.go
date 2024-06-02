package util

import (
	"github.com/go-faker/faker/v4"
	"math/rand"
	"time"
)

type RandomPost struct {
	ID        int32
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenerateRandomPost() RandomPost {
	return RandomPost{
		ID:        rand.Int31(),
		Title:     faker.Sentence(),
		Content:   faker.Paragraph(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
