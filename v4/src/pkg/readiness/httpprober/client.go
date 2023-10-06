package httpprober

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

type Prober struct {
	lg *slog.Logger

	conn *resty.Client
}

func New(lg *slog.Logger, conn *resty.Client) Prober {
	return Prober{lg: lg, conn: conn}
}

func (c Prober) Ping(key string, probe *readiness.Probe) {
	sync.OnceFunc(func() {
		probe.Mark(key, false)
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
				probe.Mark(key, true)
				c.lg.Warn("[startup] rating client ready")
			})
		}
	}()
}
