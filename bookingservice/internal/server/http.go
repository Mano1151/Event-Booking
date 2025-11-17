package server

import (
	v1 "bookingservice/api/bookingservice/v1"
	"bookingservice/internal/conf"
	"bookingservice/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/rs/cors"
)

// NewHTTPServer creates a new HTTP server with CORS support
func NewHTTPServer(c *conf.Server, bookingService *service.BookingService, logger log.Logger) *http.Server {
	// ✅ Setup CORS middleware (works for React frontend, Postman, and other origins)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // your React dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		// 👇 Add CORS as a transport filter (applied to all routes)
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

	// ✅ Create server with all options
	srv := http.NewServer(opts...)

	// ✅ Register your BookingService routes
	v1.RegisterBookingServiceHTTPServer(srv, bookingService)

	return srv
}
