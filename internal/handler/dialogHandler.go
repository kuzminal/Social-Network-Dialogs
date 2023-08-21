package handler

import (
	"Social-Net-Dialogs/models"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"time"
)

func (i *Instance) SendMessage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not delete friend from empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId, _ := i.dialogueStore.GetChatId(context.Background(), id, userId)

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	msg.Id = uuid.Must(uuid.NewV4()).String()
	msg.FromUser = userId
	msg.ToUser = id
	msg.CreatedAt = time.Now().Format("2006-01-02")

	msg.ChatId = chatId
	err = i.dialogueStore.SaveMessage(context.Background(), msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	userSessionDTO, _ := msgpack.Marshal(msg)
	go i.countersPublisher.SendMessageInfo(context.Background(), userSessionDTO)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (i *Instance) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Extract TraceID from header
	md, _ := metadata.FromOutgoingContext(r.Context())
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

	ctx, span := i.tracer.Tracer("dialog-service").Start(ctx, "GetMessages")
	defer span.End()

	id := chi.URLParam(r, "user_id")
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not delete friend from empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg, _ := i.dialogueStore.GetMessages(context.Background(), id, userId)
	msgDTO, _ := json.Marshal(msg)
	w.Header().Add("Content-Type", "application/json")
	w.Write(msgDTO)
}

func (i *Instance) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	// Extract TraceID from header
	md, _ := metadata.FromOutgoingContext(r.Context())
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

	ctx, span := i.tracer.Tracer("dialog-service").Start(ctx, "MarkAsRead")
	defer span.End()

	messageId := chi.URLParam(r, "messageId")
	user := chi.URLParam(r, "user_id")
	//userId := r.Context().Value("userId").(string)
	if len(user) == 0 {
		log.Println("Could not mark message for empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, _ := i.dialogueStore.MarkAsRead(ctx, messageId)

	msgDTO, _ := json.Marshal(message)
	msgToKafka, _ := msgpack.Marshal(message)
	go i.countersPublisher.SendMessageInfo(context.Background(), msgToKafka)
	w.Header().Add("Content-Type", "application/json")
	w.Write(msgDTO)
}
