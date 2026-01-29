package middleware

import "net/http"

// Middleware — HTTP middleware contract
type Middleware func(http.Handler) http.Handler
