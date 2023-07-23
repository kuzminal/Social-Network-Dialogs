package tarantool

import (
	"Social-Net-Dialogs/models"
	"github.com/pkg/errors"
)

func (t *Tarantool) LoadSession(token string) (models.UserSession, error) {
	var session []models.UserSession
	err := t.conn.CallTyped("get_session_by_user_id", []interface{}{token}, &session)
	if err != nil {
		return models.UserSession{}, err
	}
	if len(session) != 1 {
		return models.UserSession{}, errors.Errorf("Cannot find user with id: %s", session)
	} else {
		return session[0], nil
	}
}

func (t *Tarantool) CreateSession(m *models.UserSession) (models.UserSession, error) {
	var session []models.UserSession
	err := t.conn.CallTyped("create_session", []interface{}{m.Id, m.UserId, m.Token, m.CreatedAt}, &session)
	if err != nil {
		return models.UserSession{}, err
	}
	if len(session) != 1 {
		return models.UserSession{}, errors.Errorf("Cannot create session user with id: %s", m.Id)
	} else {
		return session[0], nil
	}
}
