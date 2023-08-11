package handler

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strings"
)

func (i *Instance) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract TraceID from header
		md, _ := metadata.FromIncomingContext(r.Context())
		traceIdString := md["x-trace-id"][0]
		// Convert string to byte array
		traceId, err := trace.TraceIDFromHex(traceIdString)
		if err != nil {
			return
		}
		// Creating a span context with a predefined trace-id
		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceId,
		})
		// Embedding span config into the context
		ctx := trace.ContextWithSpanContext(r.Context(), spanContext)

		ctx, span := i.tracer.Tracer("dialog-service").Start(ctx, "BasicAuth")
		defer span.End()

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			log.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			token := authHeader[1]
			user, e := i.tokenService.ValidateToken(context.Background(), token)
			if e != nil {
				log.Println(e)
			}
			if len(user.UserId) > 0 {
				ctx := context.WithValue(r.Context(), "userId", user.UserId)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}
	})
}
