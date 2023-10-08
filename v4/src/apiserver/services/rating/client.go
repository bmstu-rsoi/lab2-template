package rating

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/rating"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness/httpprober"
	v1 "github.com/migregal/bmstu-iu7-ds-lab2/rating/api/http/v1"
)

var probeKey = "http-rating-client"

type Client struct {
	lg *slog.Logger

	conn *resty.Client
}

func New(lg *slog.Logger, cfg rating.Config, probe *readiness.Probe) (*Client, error) {
	client := resty.New().
		SetTransport(&http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}).
		SetBaseURL(fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port))

	c := Client{
		lg:   lg,
		conn: client,
	}

	go httpprober.New(lg, client).Ping(probeKey, probe)

	return &c, nil
}

func (c *Client) GetUserRating(
	ctx context.Context, username string,
) (rating.Rating, error) {
	resp, err := c.conn.R().
		SetHeader("X-User-Name", username).
		SetResult(&v1.RatingResponse{}).
		Get("/api/v1/rating")
	if err != nil {
		return rating.Rating{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return rating.Rating{}, fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	data := resp.Result().(*v1.RatingResponse)

	rating := rating.Rating(*data)

	return rating, nil
}

func (c *Client) UpdateUserRating(
	ctx context.Context, username string, diff int,
) error {
	resp, err := c.conn.R().
		SetHeader("X-User-Name", username).
		SetQueryParam("diff", strconv.Itoa(diff)).
		SetResult(&v1.RatingResponse{}).
		Patch("/api/v1/rating")
	if err != nil {
		return fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	return nil
}
