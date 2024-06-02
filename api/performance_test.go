package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	vegeta "github.com/tsenart/vegeta/lib"
	"testing"
	"time"
)

func TestPerformance(t *testing.T) {
	createPostRequest := createPostRequest{Title: faker.Sentence(), Content: faker.Paragraph()}
	createPostRequestData, err := json.Marshal(createPostRequest)
	require.NoError(t, err)

	updatePostRequest := updatePostRequestBody{Title: faker.Sentence(), Content: faker.Paragraph()}
	updatePostRequestData, err := json.Marshal(updatePostRequest)
	require.NoError(t, err)

	targets := []vegeta.Target{
		{Method: "GET", URL: "http://0.0.0.0:8080/posts"},
		{Method: "GET", URL: "http://0.0.0.0:8080/posts/1"},
		{Method: "DELETE", URL: "http://0.0.0.0:8080/posts"},
		{Method: "POST", URL: "http://0.0.0.0:8080/posts", Body: createPostRequestData},
		{Method: "PUT", URL: "http://0.0.0.0:8080/posts/1", Body: updatePostRequestData},
	}
	attacker := vegeta.NewAttacker()
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 10 * time.Second

	for res := range attacker.Attack(vegeta.NewStaticTargeter(targets...), rate, duration, "Big Bang") {
		if res.Code == 500 {
			fmt.Println("service is out of reach")
		}
		fmt.Println(res.Attack)
	}
}
