package customiddleware

import (
	"context"
	"errors"
	"github.com/Sanchir01/currency-wallet/internal/domain/contants"
	"github.com/Sanchir01/currency-wallet/internal/feature/user"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
	"time"
)

const responseWriterKey = "responseWriter"

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requesDuration)
}

var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "count-requests",
		Subsystem: "http",
		Name:      "request_total",
		Help:      "Total number of HTTP requests",
	},
	[]string{"path", "method"},
)

var requesDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "duration",
		Name:      "http_request_duration_seconds",
		Help:      "Duration of HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

func GetJWTClaimsFromCtx(ctx context.Context) (*user.Claims, error) {
	claims, ok := ctx.Value(contextkey.UserIDCtxKey).(*user.Claims)
	if !ok {
		return nil, errors.New("no JWT claims found in context")
	}
	return claims, nil
}

func AuthMiddleware(domain string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			access, err := r.Cookie("refreshToken")
			if err != nil {
				refresh, err := r.Cookie("accessToken")
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				accessToken, err := user.NewAccessToken(refresh.Value, 0, w, domain)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				token, err := user.ParseToken(accessToken)
				if err != nil {
					slog.Error("failed parse token middleware", err.Error())
					next.ServeHTTP(w, r)
					return
				}

				ctx := context.WithValue(r.Context(), contextkey.UserIDCtxKey, token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			validAccessToken, err := user.ParseToken(access.Value)
			if err != nil {
				slog.Error("failed parse token middleware", err.Error())
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), contextkey.UserIDCtxKey, validAccessToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		requestCount.WithLabelValues(r.URL.Path, r.Method).Inc()
		requesDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
