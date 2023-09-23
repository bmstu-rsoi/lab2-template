package apiutils

import (
	"log/slog"
	"os"
)

type Serveable interface {
	ListenAndServe(adr string) error
}

type Callable func(errs chan<- error)

func NewCallable(addr string, api Serveable) Callable {
	return func(errs chan<- error) {
		errs <- api.ListenAndServe(addr)
	}
}

func Serve(lg *slog.Logger, apis ...Callable) {
	errs := make(chan error, len(apis))

	for _, api := range apis {
		go api(errs)
	}

	err := <-errs
	if err != nil {
		lg.Error("[shutdown] terminating application", "error", err.Error())
		os.Exit(1)
	}
}
