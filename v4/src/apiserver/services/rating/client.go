package rating

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/rating"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
	v1 "github.com/migregal/bmstu-iu7-ds-lab2/rating/api/http/v1"
)

var probeKey = "http-rating-client"

type Client struct {
	lg *slog.Logger

	conn *http.Client

	addr string
}

func New(lg *slog.Logger, cfg rating.Config, probe *readiness.Probe) (*Client, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	conn := http.Client{
		Transport: tr,
	}

	c := Client{
		lg:   lg,
		conn: &conn,
		addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	go c.ping(probe)

	return &c, nil
}

func (c *Client) ping(probe *readiness.Probe) {
	sync.OnceFunc(func() {
		probe.Mark(probeKey, false)
	})

	func() {
		for {
			resp, err := c.conn.Get(c.addr + "/readiness")
			if err != nil {
				continue
			}

			if resp.StatusCode != http.StatusOK {
				continue
			}

			sync.OnceFunc(func() {
				probe.Mark(probeKey, true)
				c.lg.Warn("[startup] rating client ready")
			})
		}
	}()
}

func (c *Client) GetUserRating(
	ctx context.Context, username string,
) (rating.Rating, error) {
	url := fmt.Sprintf("http://%s/api/v1/rating", c.addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return rating.Rating{}, fmt.Errorf("failed to init http request: %w", err)
	}

	req.Header.Add("X-User-Name", username)

	res, err := c.conn.Do(req)
	if err != nil {
		return rating.Rating{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return rating.Rating{}, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return rating.Rating{}, fmt.Errorf("failed to read http response")
	}

	var resp v1.RatingResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return rating.Rating{}, fmt.Errorf("failed to parse http ersponse")
	}

	rating := rating.Rating(resp)

	return rating, nil
}
