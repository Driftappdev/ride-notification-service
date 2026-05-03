package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

type Field struct {
	Key   string
	Value interface{}
}

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct{}

func NewJSONLogger(level LogLevel) *Logger { return &Logger{} }
func LogField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
func (l *Logger) Info(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)  {}
func (l *Logger) Error(msg string, fields ...Field) {}

type Category string

const (
	CategoryValidation   Category = "validation"
	CategoryUnauthorized Category = "unauthorized"
	CategoryForbidden    Category = "forbidden"
	CategoryNotFound     Category = "not_found"
	CategoryConflict     Category = "conflict"
	CategoryRateLimit    Category = "rate_limit"
	CategoryUnavailable  Category = "unavailable"
	CategoryTimeout      Category = "timeout"
	CategoryInternal     Category = "internal"
)

type BaseError struct {
	Code     string   `json:"code"`
	Status   int      `json:"status"`
	Category Category `json:"category"`
	Message  string   `json:"message"`
}

func (e *BaseError) Error() string { return e.Message }

type APIErrorResponse struct {
	Error *BaseError `json:"error"`
}

func NewBaseError(code string, status int, message string) *BaseError {
	return &BaseError{Code: code, Status: status, Category: CategoryInternal, Message: message}
}

func HTTPStatus(err error) int {
	var be *BaseError
	if errors.As(err, &be) && be.Status > 0 {
		return be.Status
	}
	return http.StatusInternalServerError
}

func APIResponse(err *BaseError) APIErrorResponse { return APIErrorResponse{Error: err} }

type Claims struct {
	Subject string
}

type AuthConfig struct{}
type AuthManager struct{}

func NewAuthManager(cfg AuthConfig) *AuthManager { return &AuthManager{} }
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	return nil, false
}

type BoolCache struct {
	mu   sync.RWMutex
	data map[string]bool
}

func NewBoolCache(maxSize int, ttl time.Duration) *BoolCache {
	return &BoolCache{data: map[string]bool{}}
}
func CacheGetBool(c *BoolCache, key string) (bool, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}
func CacheSetBool(c *BoolCache, key string, value bool, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

type CircuitBreaker struct{}

func NewCircuitBreaker(name string, openTimeout time.Duration, minVolume int64) *CircuitBreaker {
	return &CircuitBreaker{}
}

type LogMiddlewareConfig struct{}
type HTTPLogField = Field
type HTTPLogger = Logger

func NewHTTPLogMiddleware(cfg LogMiddlewareConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { return next }
}
func NewHTTPRecovery(logger HTTPLogger, extraFields ...HTTPLogField) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { return next }
}

type Counter struct{}
type Histogram struct{}

func NewRequestCounter(name, help string) *Counter { return &Counter{} }
func NewRequestLatency(name, help string) *Histogram {
	return &Histogram{}
}
func MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }
}

type PageParams struct {
	Page int
	Size int
}
type Page[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

func ParsePage(r *http.Request) (PageParams, error) {
	return PageParams{Page: 1, Size: 20}, nil
}
func BuildPage[T any](items []T, params PageParams, r *http.Request) Page[T] {
	return Page[T]{Items: items, Page: params.Page, Size: params.Size, Total: len(items)}
}
func RespondPage[T any](w http.ResponseWriter, page Page[T]) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(page)
}

type RateLimiter interface {
	Allow(ctx context.Context, key string) (RateLimitResult, error)
}
type KeyExtractor func(*http.Request) string
type RateLimitResult struct {
	Allowed bool
}

type nopLimiter struct{}

func (n nopLimiter) Allow(ctx context.Context, key string) (RateLimitResult, error) {
	return RateLimitResult{Allowed: true}, nil
}

func NewIPSlidingWindow(limit int, window time.Duration) (RateLimiter, KeyExtractor) {
	return nopLimiter{}, func(r *http.Request) string {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return r.RemoteAddr
		}
		return host
	}
}

func Allow(limiter RateLimiter, ctx context.Context, key string) (RateLimitResult, error) {
	return limiter.Allow(ctx, key)
}
func LimitKey(extractor KeyExtractor, r *http.Request) string { return extractor(r) }

type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
}

func DefaultRetryConfig(maxAttempts int, initialDelay, maxDelay time.Duration) RetryConfig {
	return RetryConfig{MaxAttempts: maxAttempts, InitialDelay: initialDelay, MaxDelay: maxDelay}
}
func RetryDo(ctx context.Context, cfg RetryConfig, fn func(context.Context) error) error {
	attempts := cfg.MaxAttempts
	if attempts < 1 {
		attempts = 1
	}
	delay := cfg.InitialDelay
	if delay <= 0 {
		delay = 10 * time.Millisecond
	}
	for i := 0; i < attempts; i++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}
		if i == attempts-1 {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		if cfg.MaxDelay > 0 && delay < cfg.MaxDelay {
			delay *= 2
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
	}
	return nil
}

func CleanString(v string) string { return strings.Join(strings.Fields(v), " ") }
func CleanEmail(v string) string  { return strings.ToLower(strings.TrimSpace(v)) }
func CleanPhone(v string) string {
	var b strings.Builder
	for _, r := range v {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
func CleanPath(v string) string { return path.Clean(strings.TrimSpace(v)) }
func CleanURL(v string) string  { return strings.TrimSpace(v) }

func DoWithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return fn(c)
}

type Tracer struct{}
type TraceSpan struct{}
type TraceAttrType struct {
	Key   string
	Value interface{}
}

func GlobalTracer(name string) *Tracer { return &Tracer{} }
func StartTraceSpan(tracer *Tracer, ctx context.Context, name string) (context.Context, *TraceSpan) {
	return ctx, &TraceSpan{}
}
func TraceAttr(key string, value interface{}) TraceAttrType {
	return TraceAttrType{Key: key, Value: value}
}
func TraceSetAttrs(span *TraceSpan, attrs ...TraceAttrType) {}
func TraceSetError(span *TraceSpan, msg string)             {}
func TraceEnd(span *TraceSpan)                              {}

func Required(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(name + " is required")
	}
	return nil
}
func MaxLen(name, value string, max int) error {
	if len(value) > max {
		return errors.New(name + " exceeds max length")
	}
	return nil
}
func ValidateEmail(value string) error {
	if !strings.Contains(value, "@") {
		return errors.New("invalid email")
	}
	return nil
}
func OneOf(name, value string, allowed ...string) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	return errors.New(name + " invalid value")
}

type Rule func(interface{}) error
type Validator struct{}
type FieldRules map[string][]Rule
type ValidationErrors map[string]string

func NewValidator() *Validator { return &Validator{} }
func ValidateMap(fields FieldRules, data map[string]interface{}) error {
	return nil
}
func ValidationErrorsOf(err error) (ValidationErrors, bool) {
	return nil, false
}
