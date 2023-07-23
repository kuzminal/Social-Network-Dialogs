package handler

import (
	"Social-Net-Dialogs/internal/service"
	"Social-Net-Dialogs/internal/store"
	"Social-Net-Dialogs/models"
)

type Instance struct {
	sessionStore     store.SessionStore
	dialogueStore    store.DialogueStore
	tokenService     *service.Client
	connectToWs      chan models.ActiveWsUsers
	disconnectFromWs chan models.ActiveWsUsers
}

func NewInstance(
	sessionStore store.SessionStore,
	dialogueStore store.DialogueStore,
	tokenService *service.Client,
	connectToWs chan models.ActiveWsUsers,
	disconnectFromWs chan models.ActiveWsUsers,
) *Instance {
	return &Instance{
		sessionStore:     sessionStore,
		dialogueStore:    dialogueStore,
		tokenService:     tokenService,
		connectToWs:      connectToWs,
		disconnectFromWs: disconnectFromWs,
	}
}
