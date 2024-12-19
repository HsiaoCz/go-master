package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/blogshare/internal/models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const (
	indexName = "blog_posts"
)

type ElasticsearchRepository struct {
	client *elasticsearch.Client
}

func NewElasticsearchRepository(client *elasticsearch.Client) *ElasticsearchRepository {
	return &ElasticsearchRepository{client: client}
}

func (r *ElasticsearchRepository) CreatePost(ctx context.Context, post *models.Post) error {
	payload, err := json.Marshal(post)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: post.ID,
		Body:       strings.NewReader(string(payload)),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("error indexing document")
	}

	return nil
}

func (r *ElasticsearchRepository) SearchPosts(ctx context.Context, query string) ([]models.Post, error) {
	// Search query
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title^2", "content", "tags"},
			},
		},
	}

	payload, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(strings.NewReader(string(payload))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("error searching documents")
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source models.Post `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	posts := make([]models.Post, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		posts[i] = hit.Source
	}

	return posts, nil
}
