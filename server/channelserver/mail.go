package channelserver

import (
	"database/sql"
	"github.com/Andoryuuta/Erupe/network/binpacket"
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
	"go.uber.org/zap"
)

type Mail struct {
	SenderID    uint32
	RecipientID uint32
	Subject     string
	Body        string
	Read        bool
}

func (m *Mail) Send(s *Session, transaction *sql.Tx) error {
	query := `
		INSERT INTO mail (sender_id, recipient_id, subject, body)
		VALUES ($1, $2, $3, $4)
	`

	var err error

	if transaction == nil {
		_, err = s.server.db.Exec(query, m.SenderID, m.RecipientID, m.Subject, m.Body)
	} else {
		_, err = transaction.Exec(query, m.SenderID, m.RecipientID, m.Subject, m.Body)
	}

	if err != nil {
		s.logger.Error(
			"failed to send mail",
			zap.Error(err),
			zap.Uint32("senderID", m.SenderID),
			zap.Uint32("recipientID", m.RecipientID),
			zap.String("subject", m.Subject),
			zap.String("body", m.Body),
		)
		return err
	}

	return nil
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
