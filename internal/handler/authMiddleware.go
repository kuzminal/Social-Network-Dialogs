package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strings"
)

func (i *Instance) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := i.tracer.Tracer("dialog-service").Start(r.Context(), "BasicAuth")
		defer span.End()

		traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
		ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			log.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			token := authHeader[1]
			user, e := i.tokenService.ValidateToken(ctx, token)
			if e != nil {
				log.Println(e)
			}
			if len(user.UserId) > 0 {
				ctx := context.WithValue(ctx, "userId", user.UserId)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}
	})
}
