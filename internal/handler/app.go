package handler

import (
	"Social-Net-Dialogs/internal/counters"
	"Social-Net-Dialogs/internal/service"
	"Social-Net-Dialogs/internal/store"
	"Social-Net-Dialogs/models"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Instance struct {
	sessionStore      store.SessionStore
	dialogueStore     store.DialogueStore
	tokenService      *service.Client
	connectToWs       chan models.ActiveWsUsers
	disconnectFromWs  chan models.ActiveWsUsers
	tracer            *trace.TracerProvider
	countersPublisher counters.Publisher
}

func NewInstance(
	sessionStore store.SessionStore,
	dialogueStore store.DialogueStore,
	tokenService *service.Client,
	connectToWs chan models.ActiveWsUsers,
	disconnectFromWs chan models.ActiveWsUsers,
	tracer *trace.TracerProvider,
	countersPublisher counters.Publisher,
) *Instance {
	return &Instance{
		sessionStore:      sessionStore,
		dialogueStore:     dialogueStore,
		tokenService:      tokenService,
		connectToWs:       connectToWs,
		disconnectFromWs:  disconnectFromWs,
		tracer:            tracer,
		countersPublisher: countersPublisher,
	}
}
