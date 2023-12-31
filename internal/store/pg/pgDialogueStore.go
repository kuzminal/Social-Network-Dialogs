package pg

import (
	"Social-Net-Dialogs/models"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"log"
)

func (pg *Postgres) GetMessages(ctx context.Context, userFrom string, userTo string) ([]models.Message, error) {
	var messages []models.Message
	query := `select id, from_user, to_user, text from social.messages where chat_id =$1 ORDER BY created_at DESC;`
	chatId, err := pg.GetChatId(ctx, userFrom, userTo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := pg.db.Query(ctx, query, chatId)
	defer rows.Close()
	if err != nil {
		return []models.Message{}, fmt.Errorf("unable to query posts: %w", err)
	}
	message := models.Message{}
	for rows.Next() {
		err = rows.Scan(&message.Id, &message.FromUser, &message.ToUser, &message.Text)
		if err != nil {
			log.Printf("unable to scan row: %v", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (pg *Postgres) GetChatId(ctx context.Context, userFrom string, userTo string) (string, error) {
	query := `select chat_id from social.chats where (user_from=$1 and user_to=$2) or (user_from=$2 and user_to=$1);`
	row := pg.db.QueryRow(ctx, query, userFrom, userTo)
	var chatId string
	err := row.Scan(&chatId)
	if err != nil && err.Error() == "no rows in result set" {
		id, _ := pg.CreateChat(ctx, userFrom, userTo)
		return id, nil
	} else if err != nil {
		return "", err
	}
	return chatId, nil
}

func (pg *Postgres) SaveMessage(ctx context.Context, msg models.Message) error {
	query := `INSERT INTO social.messages (id, "text", to_user, from_user, created_at, chat_id, is_read) 
VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err := pg.db.Exec(ctx, query, msg.Id, msg.Text, msg.ToUser, msg.FromUser, msg.CreatedAt, msg.ChatId, msg.IsRead)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) CreateChat(ctx context.Context, fromUser string, toUser string) (string, error) {
	query := `INSERT INTO social.chats (chat_id, user_to, user_from) 
VALUES ($1, $2, $3);`
	chatId := uuid.Must(uuid.NewV4()).String()
	_, err := pg.db.Exec(ctx, query, chatId, fromUser, toUser)
	if err != nil {
		return "", err
	}
	return chatId, nil
}

func (pg *Postgres) MarkAsRead(ctx context.Context, messageId string) (models.Message, error) {
	query := `update social.messages set is_read = true where id=$1 RETURNING id, text, from_user, to_user, chat_id, created_at, is_read;`
	row := pg.db.QueryRow(ctx, query, messageId)
	var message models.Message
	var createAt pgtype.Timestamp
	err := row.Scan(&message.Id, &message.Text, &message.FromUser, &message.ToUser, &message.ChatId, &createAt, &message.IsRead)
	if err != nil && err.Error() == "no rows in result set" {
		return models.Message{}, err
	}
	message.CreatedAt = createAt.Time.Format("2006-01-02")
	return message, nil
}
