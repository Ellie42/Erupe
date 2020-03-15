package channelserver

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Guild struct {
	ID          uint32
	Name        string
	Message     string
	CreatedAt   time.Time
	MemberCount uint16
	Leader      *GuildMember
}

type GuildMember struct {
	GuildID  uint32
	CharID   uint32
	JoinedAt time.Time
	Name     string
}

func GetGuildInfoByID(s *Session, guildID uint32) (*Guild, error) {
	rows, err := s.server.db.Query(`
		SELECT g.id, g.name, created_at, (
			SELECT count(1) FROM guild_characters gc WHERE gc.guild_id = g.id
		) AS member_count, leader_id, lc.name as leader_name, lgc.joined_at as leader_joined
		FROM guilds g
				 JOIN guild_characters lgc ON lgc.character_id = leader_id
				 JOIN characters lc on leader_id = lc.id
		WHERE g.id == $1 
		LIMIT 1
	`, guildID)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to retrieve guild '%d'", guildID))
		return nil, err
	}

	return buildGuildObjectFromDbResult(rows, err, s)
}

func GetGuildInfoByCharacterId(s *Session, charID uint32) (*Guild, error) {
	rows, err := s.server.db.Query(`
		SELECT g.id, g.name, created_at, (
			SELECT count(1) FROM guild_characters gc WHERE gc.guild_id = g.id
		) AS member_count, leader_id, lc.name as leader_name, lgc.joined_at as leader_joined
		FROM guilds g
				 JOIN guild_characters lgc ON lgc.character_id = leader_id
				 JOIN characters lc on leader_id = lc.id
				 JOIN guild_characters gc
					  ON g.id = gc.guild_id AND gc.character_id = $1
		LIMIT 1
	`, charID)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to retrieve guild for character '%d'", charID))
		return nil, err
	}

	defer rows.Close()

	return buildGuildObjectFromDbResult(rows, err, s)
}

func GetGuildMembers(s *Session, guildID uint32, memberCount uint16) ([]*GuildMember, error) {
	rows, err := s.server.db.Queryx(`
		SELECT guild_id, joined_at, name, gc.character_id
			FROM guild_characters gc
				JOIN characters c on gc.character_id = c.id
			WHERE guild_id = $1
	`, guildID)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to retrieve membership data for guild '%d'", guildID))
		return nil, err
	}

	defer rows.Close()

	members := make([]*GuildMember, memberCount)

	i := 0

	for rows.Next() {
		member, err := buildGuildMemberObjectFromDBResult(rows, err, s)

		if err != nil {
			return nil, err
		}

		members[i] = member
		i++
	}

	return members, nil
}

func buildGuildObjectFromDbResult(result *sql.Rows, err error, s *Session) (*Guild, error) {
	hasRow := result.Next()

	if !hasRow {
		return nil, nil
	}

	guild := &Guild{
		Leader: &GuildMember{},
	}

	err = result.Scan(
		&guild.ID, &guild.Name, &guild.CreatedAt, &guild.MemberCount,
		&guild.Leader.CharID, &guild.Leader.Name, &guild.Leader.JoinedAt,
	)

	guild.Leader.GuildID = guild.ID
	guild.Leader.JoinedAt = guild.CreatedAt

	if err != nil {
		s.logger.Error("failed to retrieve guild data from database")
		return nil, err
	}

	return guild, nil
}

func GetCharacterGuildData(s *Session, charID uint32) (*GuildMember, error) {
	rows, err := s.server.db.Queryx(`
		SELECT guild_id, joined_at, name, character_id
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

	err = rows.Scan(&memberData.GuildID, &memberData.JoinedAt, &memberData.Name, &memberData.CharID)

	if err != nil {
		s.logger.Error("failed to retrieve guild data from database")
		return nil, err
	}

	return memberData, nil
}
