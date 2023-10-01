package library

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/library"
	v1 "github.com/migregal/bmstu-iu7-ds-lab2/library/api/http/v1"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

var probeKey = "http-library-client"

type Client struct {
	lg *slog.Logger

	conn *http.Client

	addr string
}

func New(lg *slog.Logger, cfg library.Config, probe *readiness.Probe) (*Client, error) {
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
				c.lg.Warn("[startup] library client ready")
			})
		}
	}()
}

func (c *Client) GetLibraries(
	ctx context.Context, city string, page uint64, size uint64,
) (library.Libraries, error) {
	url := fmt.Sprintf("http://%s/api/v1/libraries", c.addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return library.Libraries{}, fmt.Errorf("failed to init http request: %w", err)
	}

	q := req.URL.Query()
	q.Add("city", city)
	q.Add("page", strconv.FormatUint(page, 10))
	if size == 0 {
		size = math.MaxUint64
	}
	q.Add("size", strconv.FormatUint(size, 10))
	req.URL.RawQuery = q.Encode()

	res, err := c.conn.Do(req)
	if err != nil {
		return library.Libraries{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return library.Libraries{}, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return library.Libraries{}, fmt.Errorf("failed to read http response")
	}

	var resp v1.LibrariesResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return library.Libraries{}, fmt.Errorf("failed to parse http ersponse")
	}

	books := library.Libraries{Total: resp.Total}
	for _, book := range resp.Items {
		books.Items = append(books.Items, library.Library(book))
	}

	return books, nil
}

func (c *Client) GetBooks(
	ctx context.Context, libraryID string, showAll bool, page uint64, size uint64,
) (library.LibraryBooks, error) {
	url := fmt.Sprintf("http://%s/api/v1/libraries/%s/books", c.addr, libraryID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return library.LibraryBooks{}, fmt.Errorf("failed to init http request: %w", err)
	}

	q := req.URL.Query()
	if showAll {
		q.Add("show_all", "1")
	}
	q.Add("page", strconv.FormatUint(page, 10))
	if size == 0 {
		size = math.MaxUint64
	}
	q.Add("size", strconv.FormatUint(size, 10))
	req.URL.RawQuery = q.Encode()

	res, err := c.conn.Do(req)
	if err != nil {
		return library.LibraryBooks{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return library.LibraryBooks{}, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return library.LibraryBooks{}, fmt.Errorf("failed to read http response")
	}

	var resp v1.BooksResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return library.LibraryBooks{}, fmt.Errorf("failed to parse http ersponse")
	}

	books := library.LibraryBooks{Total: resp.Total}
	for _, book := range resp.Items {
		books.Items = append(books.Items, library.Book(book))
	}

	return books, nil
}

func (c *Client) ObtainBook(ctx context.Context, libraryID string, bookID string) error {
	return nil
}
