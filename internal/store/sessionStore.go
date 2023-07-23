package store

import (
	"Social-Net-Dialogs/models"
)

type SessionStore interface {
	Store

	LoadSession(id string) (models.UserSession, error)
	CreateSession(m *models.UserSession) (models.UserSession, error)
}
