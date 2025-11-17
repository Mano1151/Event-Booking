package server

import (
	v1 "eventservice/api/eventservice/v1"
	"eventservice/internal/conf"
	"eventservice/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/rs/cors"
)

// NewHTTPServer creates a new HTTP server with CORS support
func NewHTTPServer(c *conf.Server, greeter *service.ShowEventService, logger log.Logger) *http.Server {
	// Setup CORS middleware (works both for frontend + Postman + backend)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // frontend dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		// 👇 Add CORS as a transport filter (applied to all handlers)
		http.Filter(corsMiddleware.Handler),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	// Create server with options
	srv := http.NewServer(opts...)

	// Register routes
	v1.RegisterEventServiceHTTPServer(srv, greeter)

	return srv
}
