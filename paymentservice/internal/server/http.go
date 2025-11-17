package server

import (
	v1 "paymentservice/api/paymentservice/v1"
	"paymentservice/internal/conf"
	"paymentservice/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/rs/cors"
)

// NewHTTPServer creates a new HTTP server with CORS support
func NewHTTPServer(c *conf.Server, paymentService *service.PaymentService, logger log.Logger) *http.Server {
	// Setup CORS middleware (allow frontend origin)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		// Add CORS middleware as transport filter (applied to all HTTP routes)
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

	srv := http.NewServer(opts...)

	// Register the payment service HTTP server
	v1.RegisterPaymentServiceHTTPServer(srv, paymentService)

	return srv
}
