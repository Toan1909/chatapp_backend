package repoimpl

import (
	"chatapp/db"
	my_err "chatapp/err"
	"chatapp/model"
	"chatapp/model/req"
	"chatapp/mylog"
	"chatapp/repo"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ConversRepoImpl struct {
	sql *db.Sql
}

// MarkMessageAsSeen implements repo.ConversRepo.
func (g ConversRepoImpl) MarkMessageAsSeen(c context.Context, messageId string, userId string) (model.Message,error) {
	var mes =model.Message{}
	_, err := g.sql.Db.Exec(`
		UPDATE messages
		SET seen_by = array_append(seen_by, $1)
		WHERE message_id = $2 AND NOT $1 = ANY(seen_by)`, userId, messageId)
	if err != nil {
		return mes,err
	}
	err = g.sql.Db.GetContext(c, &mes, `SELECT * FROM messages WHERE message_id = $1`, messageId)
	if err != nil {
		return mes, err
	}

	return mes, nil

}

// LoadListMembers implements repo.ConversRepo.
func (g ConversRepoImpl) LoadListMembers(c context.Context, conversId string) ([]model.ConversationMember, error) {
	var listMem []model.ConversationMember
	statement := `SELECT 
					users.user_id, 
					users.fullname, 
					users.phone, 
					users.url_profile_pic, 
					users.status,
					conversation_members.joined_at AS joined_at
				FROM users
				INNER JOIN conversation_members
				ON conversation_members.conversation_id = $1  AND conversation_members.user_id = users.user_id
				`
	err := g.sql.Db.SelectContext(c, &listMem, statement, conversId)
	if err != nil {
		if err == sql.ErrNoRows {
			return listMem, my_err.MemNotFound
		}
		return listMem, err
	}
	return listMem, nil
}

// LoadListConvers implements repo.ConversRepo.
func (g ConversRepoImpl) LoadListConvers(c context.Context, userId string) ([]model.Conversation, error) {
	var listConvers []model.Conversation
	statement := `SELECT 
					conversations.conversation_id, 
					conversations.conversation_name, 
					conversations.is_group, 
					conversations.created_at, 
					conversations.updated_at
				FROM conversations
				INNER JOIN conversation_members
				ON conversation_members.user_id = $1 AND conversation_members.conversation_id = conversations.conversation_id

	`
	err := g.sql.Db.SelectContext(c, &listConvers, statement, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return listConvers, my_err.ConvsersNotFound
		}
		return listConvers, err
	}
	return listConvers, nil
}

// LoadMessages implements repo.ConversRepo.
func (g ConversRepoImpl) LoadMessages(c context.Context, conversId string) ([]model.Message, error) {
	var listMessage []model.Message
	statement := `SELECT 
					*
					FROM messages
					WHERE messages.conversation_id = $1
					ORDER BY sent_at desc
					LIMIT 30`
	err := g.sql.Db.SelectContext(c, &listMessage, statement, conversId)
	if err != nil {
		if err == sql.ErrNoRows {
			return listMessage, my_err.ConvsersNotFound
		}
		return listMessage, err
	}
	return listMessage, nil
}

// SendMessage implements repo.ConversRepo.
func (g ConversRepoImpl) SendMessage(c context.Context, req req.SendMessage) (model.Message, error) {
	statement := `
	INSERT INTO 
		messages(
			message_id,
			conversation_id,
			sender_id,
			content,
			media_url,
			sent_at
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING message_id, sent_at
	`
	msg := model.Message{
		ConversationId: req.ConversationId,
		SenderId:       req.SenderId,
		Content:        req.Content,
		MediaUrl:       req.MediaUrl,
		SendAt:         time.Now(),
	}

	// Tạo UUID cho message_id
	id, _ := uuid.NewUUID()
	msg.MessageId = id.String()

	// Thực hiện câu lệnh SQL
	err := g.sql.Db.QueryRowContext(c, statement, msg.MessageId, msg.ConversationId, msg.SenderId, msg.Content, msg.MediaUrl, msg.SendAt).
		Scan(&msg.MessageId, &msg.SendAt)
	if err != nil {
		mylog.LogError(err)
		return msg, err
	}

	return msg, nil
}

// AddMember implements repo.ConversRepo.
func (g ConversRepoImpl) AddMember(c context.Context, memId string, conversId string) error {
	statement := `
	INSERT INTO 
		conversation_members(
			conversation_id,
			user_id,
			joined_at,
			last_read_at
		)
		VALUES(
			$1,
			$2,
			$3,
			$4
		)`
	_, err := g.sql.Db.ExecContext(c, statement, conversId, memId, time.Now(), time.Now())
	if err != nil {
		mylog.LogError(err)
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return my_err.MemberConflict
			}
		}
	}
	return nil
}

// CreateConvers implements repo.ConversRepo.
func (g ConversRepoImpl) CreateConvers(c context.Context, req req.ReqCreateConvers) (model.Conversation, error) {
	statement := `
	INSERT INTO 
		conversations(
			conversation_id,
			conversation_name,
			is_group,
			created_at,
			updated_at
		)
		VALUES(
			:conversation_id,
			:conversation_name,
			:is_group,
			:created_at,
			:updated_at
		)`
	convers := model.Conversation{}
	convers.ConversationName = req.ConversationName
	convers.IsGroup = false
	convers.CreatedAt = time.Now()
	convers.UpdatedAt = time.Now()
	id, _ := uuid.NewUUID()
	convers.ConversationId = id.String()
	_, err := g.sql.Db.NamedExecContext(c, statement, convers)
	if err != nil {
		mylog.LogError(err)
		return convers, err
	}
	return convers, nil
}

func NewConversRepoImpl(sql *db.Sql) repo.ConversRepo {
	return ConversRepoImpl{sql: sql}
}
