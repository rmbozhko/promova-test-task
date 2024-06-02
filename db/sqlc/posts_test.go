package db

import (
	"context"
	"database/sql"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePost(t *testing.T) {
	populateDBWithValidRandomPost(t)
}

func TestGetAllPosts(t *testing.T) {
	postsTotalNumber := 5
	for i := 0; i < postsTotalNumber; i++ {
		populateDBWithValidRandomPost(t)
	}

	posts, err := testQueries.GetPosts(context.Background())

	require.NoError(t, err)
	require.NotEmpty(t, posts)
	require.True(t, postsTotalNumber <= len(posts))
}

func TestGetPostById(t *testing.T) {
	createdPost := populateDBWithValidRandomPost(t)

	post, err := testQueries.GetPostById(context.Background(), createdPost.ID)

	checkFetchedPostIsValid(t, err, post, createdPost)
}

func TestGetPostById_NotFound(t *testing.T) {
	post, err := testQueries.GetPostById(context.Background(), -1)

	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, post)
}

func TestUpdatePostById(t *testing.T) {
	createdPost := populateDBWithValidRandomPost(t)

	arg := UpdatePostByIdParams{ID: createdPost.ID, Title: faker.Sentence(), Content: faker.Paragraph()}

	alteredPost := createdPost
	alteredPost.Title = arg.Title
	alteredPost.Content = arg.Content

	post, err := testQueries.UpdatePostById(context.Background(), arg)

	checkFetchedPostIsValid(t, err, post, alteredPost)
}

func TestUpdatePostById_NotFound(t *testing.T) {
	post, err := testQueries.UpdatePostById(context.Background(), UpdatePostByIdParams{ID: -1, Title: faker.Sentence(), Content: faker.Paragraph()})

	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, post)
}

func TestDeletePostById(t *testing.T) {
	createdPost := populateDBWithValidRandomPost(t)

	err := testQueries.DeletePost(context.Background(), createdPost.ID)

	require.NoError(t, err)

	post, err := testQueries.GetPostById(context.Background(), createdPost.ID)

	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, post)
}

func populateDBWithValidRandomPost(t *testing.T) Post {
	title := faker.Sentence()
	content := faker.Paragraph()
	arg := CreatePostParams{Title: title, Content: content}

	subscription, err := testQueries.CreatePost(context.Background(), arg)

	checkInsertedPostIsValid(t, err, subscription, arg)

	return subscription
}

func checkInsertedPostIsValid(t *testing.T, err error, actual Post, expected CreatePostParams) {
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.NotZero(t, actual.ID)

	require.Equal(t, actual.Title, expected.Title)
	require.Equal(t, actual.Content, expected.Content)
}

func checkFetchedPostIsValid(t *testing.T, err error, actual Post, expected Post) {
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.NotZero(t, actual.ID)

	require.Equal(t, actual.Title, expected.Title)
	require.Equal(t, actual.Content, expected.Content)
}
