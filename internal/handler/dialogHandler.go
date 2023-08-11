package handler

import (
	"Social-Net-Dialogs/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
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

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	msg.ChatId = chatId
	err = i.dialogueStore.SaveMessage(context.Background(), msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (i *Instance) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx, span := i.tracer.Tracer("dialog-service").Start(r.Context(), "GetMessages")
	defer span.End()

	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

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
