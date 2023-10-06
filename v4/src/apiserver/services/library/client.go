package library

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/library"
	v1 "github.com/migregal/bmstu-iu7-ds-lab2/library/api/http/v1"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

var probeKey = "http-library-client"

type Client struct {
	lg *slog.Logger

	conn *resty.Client
}

func New(lg *slog.Logger, cfg library.Config, probe *readiness.Probe) (*Client, error) {
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

	go c.ping(probe)

	return &c, nil
}

func (c *Client) ping(probe *readiness.Probe) {
	sync.OnceFunc(func() {
		probe.Mark(probeKey, false)
	})

	func() {
		for {
			resp, err := c.conn.R().Get("/readiness")
			if err != nil {
				continue
			}

			if resp.StatusCode() != http.StatusOK {
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
	q := map[string]string{
		"city": city,
		"page": strconv.FormatUint(page, 10),
	}
	if size == 0 {
		size = math.MaxUint64
	}
	q["size"] = strconv.FormatUint(size, 10)

	resp, err := c.conn.R().
		SetQueryParams(q).
		SetResult(&v1.LibrariesResponse{}).
		Get("/api/v1/libraries")
	if err != nil {
		return library.Libraries{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.Libraries{}, fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	data := resp.Result().(*v1.LibrariesResponse)

	books := library.Libraries{Total: data.Total}
	for _, book := range data.Items {
		books.Items = append(books.Items, library.Library(book))
	}

	return books, nil
}

func (c *Client) GetBooks(
	ctx context.Context, libraryID string, showAll bool, page uint64, size uint64,
) (library.LibraryBooks, error) {
	q := map[string]string{}
	if showAll {
		q["show_all"] = "1"
	}
	q["page"] = strconv.FormatUint(page, 10)
	if size == 0 {
		size = math.MaxUint64
	}
	q["size"] = strconv.FormatUint(size, 10)

	resp, err := c.conn.R().
		SetQueryParams(q).
		SetPathParam("library_id", libraryID).
		SetResult(&v1.BooksResponse{}).
		Get("/api/v1/libraries/{library_id}/books")
	if err != nil {
		return library.LibraryBooks{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.LibraryBooks{}, fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	data := resp.Result().(*v1.BooksResponse)

	books := library.LibraryBooks{Total: data.Total}
	for _, book := range data.Items {
		books.Items = append(books.Items, library.Book(book))
	}

	return books, nil
}

func (c *Client) ObtainBook(ctx context.Context, libraryID string, bookID string) (library.ReservedBook, error) {
	body, err := json.Marshal(v1.TakeBookRequest{
		BookID:    bookID,
		LibraryID: libraryID,
	})
	if err != nil {
		return library.ReservedBook{}, fmt.Errorf("failed to format json body: %w", err)
	}

	resp, err := c.conn.R().
		SetBody(body).
		SetResult(&v1.TakeBookResponse{}).
		Post("/api/v1/books")
	if err != nil {
		return library.ReservedBook{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.ReservedBook{}, fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	data := resp.Result().(*v1.TakeBookResponse)

	return library.ReservedBook{
		Book:    library.Book(data.Book),
		Library: library.Library(data.Library),
	}, nil
}
