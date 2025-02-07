package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/peterhellberg/llm"
)

var ErrInvalidOptions = fmt.Errorf("invalid options")

type Client struct {
	base *url.URL
	http *http.Client

	apiKey               string
	contentKey           string
	collectionName       string
	collectionVectorSize int
}

func NewClient(o Options) (*Client, error) {
	if o.Base == nil {
		return nil, fmt.Errorf("%w: missing Qdrant URL", ErrInvalidOptions)
	}

	if o.CollectionName == "" {
		return nil, fmt.Errorf("%w: missing collection name", ErrInvalidOptions)
	}

	if o.ContentKey == "" {
		return nil, fmt.Errorf("%w: missing content key", ErrInvalidOptions)
	}

	if o.ContentKey == "" {
		return nil, fmt.Errorf("%w: missing content key", ErrInvalidOptions)
	}

	if o.HTTP == nil {
		o.HTTP = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}

	return &Client{
		base: o.Base,
		http: o.HTTP,

		apiKey:               o.APIKey,
		contentKey:           o.ContentKey,
		collectionName:       o.CollectionName,
		collectionVectorSize: o.CollectionVectorSize,
	}, nil
}

// UpsertPoints updates or inserts points into the Qdrant collection.
func (s *Client) UpsertPoints(ctx context.Context, vectors [][]float32, payloads []map[string]interface{}) ([]string, error) {
	ids := make([]string, len(vectors))

	for i := range ids {
		ids[i] = uuid.NewString()
	}

	payload := upsertBody{
		Batch: upsertBatch{
			IDs:      ids,
			Vectors:  vectors,
			Payloads: payloads,
		},
	}

	u := s.base.JoinPath("collections", s.collectionName, "points")

	body, status, err := s.do(ctx, http.MethodPut, u, payload)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	if status == http.StatusOK {
		return ids, nil
	}

	return nil, newAPIError("upserting vectors", body)
}

// SearchPoints queries the Qdrant collection for points based on the provided parameters.
func (s *Client) SearchPoints(ctx context.Context, vector []float32, numVectors int, scoreThreshold float32, filter any) ([]llm.Document, error) {
	payload := searchBody{
		WithPayload: true,
		WithVector:  false,
		Vector:      vector,
		Limit:       numVectors,
		Filter:      filter,
	}

	if scoreThreshold != 0 {
		payload.ScoreThreshold = scoreThreshold
	}

	u := s.base.JoinPath("collections", s.collectionName, "points", "search")

	body, statusCode, err := s.do(ctx, http.MethodPost, u, payload)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	if statusCode != http.StatusOK {
		return nil, newAPIError("querying collection", body)
	}

	var response searchResponse

	dec := json.NewDecoder(body)

	if err := dec.Decode(&response); err != nil {
		return nil, err
	}

	docs := make([]llm.Document, len(response.Result))

	for i, match := range response.Result {
		pageContent, ok := match.Payload[s.contentKey].(string)
		if !ok {
			return nil, fmt.Errorf("payload does not contain content key '%s'", s.contentKey)
		}
		delete(match.Payload, s.contentKey)

		doc := llm.Document{
			PageContent: pageContent,
			Metadata:    match.Payload,
			Score:       match.Score,
		}

		docs[i] = doc
	}

	return docs, nil
}

// do performs an HTTP request to the Qdrant API.
func (c *Client) do(ctx context.Context, method string, u *url.URL, payload any) (io.ReadCloser, int, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequestWithContext(ctx, method, u.String()+"?wait=true", body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.apiKey)

	r, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return r.Body, r.StatusCode, err
}

// newAPIError creates an error based on the Qdrant API response.
func newAPIError(task string, body io.ReadCloser) error {
	buf := new(bytes.Buffer)

	_, err := io.Copy(buf, body)
	if err != nil {
		return fmt.Errorf("failed to read body of error message: %w", err)
	}

	return fmt.Errorf("%s: %s", task, buf.String())
}

func (c *Client) ClearCollection(ctx context.Context) error {
	if c.collectionVectorSize < 1 {
		return fmt.Errorf("vector size not set, unable to clear the collection")
	}

	reqs, err := c.clearCollectionRequests(ctx)
	if err != nil {
		return err
	}

	for _, req := range reqs {
		if _, err := c.http.Do(req); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) clearCollectionRequests(ctx context.Context) ([]*http.Request, error) {
	colurl, err := newCollectionRawURL(c.base.String(), c.collectionName)
	if err != nil {
		return nil, err
	}

	deleteCollectionReq, err := newRequest(ctx, http.MethodDelete, colurl, c.apiKey, nil)
	if err != nil {
		return nil, err
	}

	createCollectionReq, err := newRequest(ctx, http.MethodPut, colurl, c.apiKey,
		strings.NewReader(fmt.Sprintf(`{
			"vectors":{"size":%d,"distance":%q},
			"optimizers_config":{"default_segment_number": 8}
		}`, c.collectionVectorSize, "Cosine")),
	)
	if err != nil {
		return nil, err
	}

	return []*http.Request{
		deleteCollectionReq,
		createCollectionReq,
	}, nil
}

func newRequest(ctx context.Context, method, rawurl, key string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawurl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Api-Key", key)

	if method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func newCollectionRawURL(rawurl, collection string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	ref := &url.URL{
		Path: "/collections/" + collection,
	}

	return u.ResolveReference(ref).String(), nil
}
