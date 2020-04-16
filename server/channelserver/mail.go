package channelserver

import (
	"database/sql"
	"github.com/Andoryuuta/Erupe/network/binpacket"
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
	"go.uber.org/zap"
	"time"
)

type Mail struct {
	ID                   int       `db:"id"`
	SenderID             uint32    `db:"sender_id"`
	RecipientID          uint32    `db:"recipient_id"`
	Subject              string    `db:"subject"`
	Body                 string    `db:"body"`
	Read                 bool      `db:"read"`
	Deleted              bool      `db:"deleted"`
	AttachedItemReceived bool      `db:"attached_item_received"`
	AttachedItemID       *uint16   `db:"attached_item"`
	AttachedItemAmount   int16     `db:"attached_item_amount"`
	CreatedAt            time.Time `db:"created_at"`
	IsGuildInvite        bool      `db:"is_guild_invite"`
	SenderName           string    `db:"sender_name"`
}

func (m *Mail) Send(s *Session, transaction *sql.Tx) error {
	query := `
		INSERT INTO mail (sender_id, recipient_id, subject, body, attached_item, attached_item_amount, is_guild_invite)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var err error

	if transaction == nil {
		_, err = s.server.db.Exec(query, m.SenderID, m.RecipientID, m.Subject, m.Body, m.AttachedItemID, m.AttachedItemAmount, m.IsGuildInvite)
	} else {
		_, err = transaction.Exec(query, m.SenderID, m.RecipientID, m.Subject, m.Body, m.AttachedItemID, m.AttachedItemAmount, m.IsGuildInvite)
	}

	if err != nil {
		s.logger.Error(
			"failed to send mail",
			zap.Error(err),
			zap.Uint32("senderID", m.SenderID),
			zap.Uint32("recipientID", m.RecipientID),
			zap.String("subject", m.Subject),
			zap.String("body", m.Body),
			zap.Uint16p("itemID", m.AttachedItemID),
			zap.Int16("itemAmount", m.AttachedItemAmount),
			zap.Bool("isGuildInvite", m.IsGuildInvite),
		)
		return err
	}

	return nil
}

func (m *Mail) MarkRead(s *Session) error {
	_, err := s.server.db.Exec(`
		UPDATE mail SET read = true WHERE id = $1 
	`, m.ID)

	if err != nil {
		s.logger.Error(
			"failed to mark mail as read",
			zap.Error(err),
			zap.Int("mailID", m.ID),
		)
		return err
	}

	return nil
}

func (m *Mail) MarkDeleted(s *Session) error {
	_, err := s.server.db.Exec(`
		UPDATE mail SET deleted = true WHERE id = $1 
	`, m.ID)

	if err != nil {
		s.logger.Error(
			"failed to mark mail as deleted",
			zap.Error(err),
			zap.Int("mailID", m.ID),
		)
		return err
	}

	return nil
}

func GetMailListForCharacter(s *Session, charID uint32) ([]Mail, error) {
	rows, err := s.server.db.Queryx(`
		SELECT 
			m.id,
			m.sender_id,
			m.recipient_id,
			m.subject,
			m.read,
			m.attached_item,
			m.attached_item_amount,
			m.created_at,
			m.is_guild_invite,
			m.deleted,
			c.name as sender_name
		FROM mail m 
			JOIN characters c ON c.id = m.sender_id 
		WHERE recipient_id = $1 AND deleted = false
		ORDER BY m.created_at DESC, id DESC
		LIMIT 32
	`, charID)

	if err != nil {
		s.logger.Error("failed to get mail for character", zap.Error(err), zap.Uint32("charID", charID))
		return nil, err
	}

	defer rows.Close()

	allMail := make([]Mail, 0)

	for rows.Next() {
		mail := Mail{}

		err := rows.StructScan(&mail)

		if err != nil {
			return nil, err
		}

		allMail = append(allMail, mail)
	}

	return allMail, nil
}

func GetMailByID(s *Session, ID int) (*Mail, error) {
	row := s.server.db.QueryRowx(`
		SELECT 
			m.id,
			m.sender_id,
			m.recipient_id,
			m.subject,
			m.read,
			m.body,
			m.attached_item,
			m.attached_item_amount,
			m.created_at,
			m.is_guild_invite,
			m.deleted,
			c.name as sender_name
		FROM mail m 
			JOIN characters c ON c.id = m.sender_id 
		WHERE m.id = $1
		LIMIT 1
	`, ID)

	mail := &Mail{}

	err := row.StructScan(mail)

	if err != nil {
		s.logger.Error(
			"failed to retrieve mail",
			zap.Error(err),
			zap.Int("mailID", ID),
		)
		return nil, err
	}

	return mail, nil
}

func SendMailNotification(s *Session, m *Mail, recipient *Session) {
	senderName, err := getCharacterName(s, m.SenderID)

	if err != nil {
		panic(err)
	}

	bf := byteframe.NewByteFrame()

	notification := &binpacket.MsgBinMailNotify{
		SenderName: senderName,
	}

	notification.Build(bf)

	castedBinary := &mhfpacket.MsgSysCastedBinary{
		CharID:         m.SenderID,
		BroadcastType:  0x00,
		MessageType:    BinaryMessageTypeMailNotify,
		RawDataPayload: bf.Data(),
	}

	castedBinary.Build(bf)

	recipient.QueueSendMHF(castedBinary)
}

func getCharacterName(s *Session, charID uint32) (string, error) {
	row := s.server.db.QueryRow("SELECT name FROM characters WHERE id = $1", charID)

	charName := ""

	err := row.Scan(&charName)

	if err != nil {
		return "", err
	}

	return charName, nil
}
