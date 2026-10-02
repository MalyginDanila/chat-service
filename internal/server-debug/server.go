package serverdebug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/MalyginDanila/chat-service/internal/buildinfo"
	"github.com/MalyginDanila/chat-service/internal/logger"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(addr string) (*Server, error) {
	if addr == "" {
		return nil, errors.New("empty addr")
	}

	lg := zap.L().Named("server-debug")

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}

	index := newIndexPage()

	e.GET("/version", s.Version)
	index.addPage("/version", "Get build information")

	e.GET("/debug/error", s.DebugError)
	index.addPage("/debug/error", "Write an error level log message")

	e.Any("/log/level", echo.WrapHandler(logger.Level))
	index.addPage("/log/level", "Get (GET) or change (PUT) log level")

	e.GET("/debug/pprof/", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	e.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	e.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	e.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	e.GET("/debug/pprof/:profile", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	index.addPage("/debug/pprof/", "Go std profiler")
	index.addPage("/debug/pprof/profile?seconds=30", "Take half-min CPU profile")

	e.GET("/", index.handler)
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(shutdownCtx) //nolint:contextcheck // родительский ctx уже отменён
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}

func (s *Server) Version(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, buildinfo.BuildInfo)
}

func (s *Server) DebugError(c echo.Context) error {
	s.lg.Error("manual error log triggered via /debug/error")
	return c.String(http.StatusOK, "Error logged\n")
}
