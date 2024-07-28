package post

import (
	"encoding/json"
	"os"
)

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Posts []Post

func GetPostsByOffsetAndLimit(offset int, limit int) (Posts, error) {
	bytes, err := getResultsBytes()
	if err != nil {
		return Posts{}, err
	}
	var posts Posts
	err = json.Unmarshal(bytes, &posts)
	if err != nil {
		return Posts{}, err
	}
	var filteredPosts Posts
	for i, post := range posts {
		if i >= offset && i < offset+limit {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts, nil
}

func GetFirstCursor() (int, error) {
	var posts Posts
	bytes, err := getResultsBytes()
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(bytes, &posts)
	if err != nil {
		return 0, err
	}
	var post Post
	for i, p := range posts {
		if i == 0 {
			post = p
			break
		}
	}
	return post.ID, nil
}

func GetPostsByCursorAndLimit(cursor int, limit int) (Posts, int, error) {
	var posts Posts
	bytes, err := getResultsBytes()
	if err != nil {
		return Posts{}, 0, err
	}
	err = json.Unmarshal(bytes, &posts)
	if err != nil {
		return Posts{}, 0, err
	}
	var filteredPosts Posts
	var cursorIndex int
	for i, post := range posts {
		if post.ID == cursor {
			cursorIndex = i
		}
		if post.ID >= cursor {
			if i >= cursorIndex && i < cursorIndex+limit {
				filteredPosts = append(filteredPosts, post)
			}
		}
	}
	if len(filteredPosts) == 0 {
		return Posts{}, len(posts) + 1, nil
	}
	lastPost := filteredPosts[len(filteredPosts)-1]
	var nextCursor int
	for _, post := range posts {
		if post.ID == lastPost.ID {
			nextCursor = post.ID + 1
			break
		}
	}
	return filteredPosts, nextCursor, nil
}

func getResultsBytes() ([]byte, error) {
	bytes, err := os.ReadFile("./post/results.json")
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
