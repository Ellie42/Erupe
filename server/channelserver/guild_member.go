package channelserver

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type GuildMember struct {
	GuildID         uint32    `db:"guild_id"`
	CharID          uint32    `db:"character_id"`
	JoinedAt        time.Time `db:"joined_at"`
	Name            string    `db:"name"`
	IsApplicant     bool      `db:"is_applicant"`
	OrderIndex      uint8     `db:"order_index"`
	LastLogin       uint32    `db:"last_login"`
	AvoidLeadership bool      `db:"avoid_leadership"`
}

func (gm *GuildMember) IsSubLeader() bool {
	return gm.OrderIndex <= 3 && !gm.AvoidLeadership
}

func (gm *GuildMember) Save(s *Session) error {
	_, err := s.server.db.Exec(`
		UPDATE guild_characters
			SET avoid_leadership=$1
		WHERE character_id=$2
	`, gm.AvoidLeadership, gm.CharID)

	if err != nil {
		s.logger.Error(
			"failed to update guild member data",
			zap.Error(err),
			zap.Uint32("charID", gm.CharID),
			zap.Uint32("guildID", gm.GuildID),
		)
		return err
	}

	return nil
}

func GetGuildMembers(s *Session, guildID uint32, applicants bool) ([]*GuildMember, error) {
	rows, err := s.server.db.Queryx(`
		SELECT guild_id, joined_at, name, gc.character_id, gc.is_applicant, gc.order_index, c.last_login, gc.avoid_leadership
			FROM guild_characters gc
				JOIN characters c on gc.character_id = c.id
			WHERE guild_id = $1 AND is_applicant = $2
	`, guildID, applicants)

	if err != nil {
		s.logger.Error("failed to retrieve membership data for guild", zap.Error(err), zap.Uint32("guildID", guildID))
		return nil, err
	}

	defer rows.Close()

	members := make([]*GuildMember, 0)

	for rows.Next() {
		member, err := buildGuildMemberObjectFromDBResult(rows, err, s)

		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func GetCharacterGuildData(s *Session, charID uint32) (*GuildMember, error) {
	rows, err := s.server.db.Queryx(`
		SELECT guild_id, joined_at, name, character_id, gc.is_applicant, gc.avoid_leadership, gc.order_index, c.last_login
			FROM guild_characters gc
				JOIN characters c on gc.character_id = c.id
			WHERE character_id=$1
		LIMIT 1
	`, charID)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to retrieve membership data for character '%d'", charID))
		return nil, err
	}

	defer rows.Close()

	hasRow := rows.Next()

	if !hasRow {
		return nil, nil
	}

	return buildGuildMemberObjectFromDBResult(rows, err, s)
}

func buildGuildMemberObjectFromDBResult(rows *sqlx.Rows, err error, s *Session) (*GuildMember, error) {
	memberData := &GuildMember{}

	err = rows.StructScan(&memberData)

	if err != nil {
		s.logger.Error("failed to retrieve guild data from database", zap.Error(err))
		return nil, err
	}

	return memberData, nil
}
