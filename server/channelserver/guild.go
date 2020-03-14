package channelserver

import "fmt"

type Guild struct {
	ID          int
	Name        string
	Message     string
	CreatedAt   uint32
	MemberCount int
	Leader      GuildMember
}

type GuildMember struct {
	GuildID  uint32
	CharID   uint32
	JoinedAt uint32
	Name     string
}

func GetGuildByCharacterId(s *Session, charID uint32) (*Guild, error) {
	result, err := s.server.db.Query(`
		SELECT g.id, gc.name, created_at, (
		    SELECT count(1) FROM guild_characters gc WHERE gc.guild_id = g.id 
		) AS member_count, leader.id, leader.name
		FROM guilds g
		    JOIN (
		    	SELECT id, name from characters c WHERE c.id = g.leader_id
			) as leader
			INNER JOIN guild_characters gc
				ON g.id = gc.guild_id AND gc.character_id = $1
		LIMIT 1
	`, charID)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to retrieve guild for character '%d'", charID))
		return nil, err
	}

	hasRow := result.Next()

	if !hasRow {
		return nil, nil
	}

	guild := &Guild{
		Leader: GuildMember{},
	}

	err = result.Scan(
		&guild.ID, &guild.Name, &guild.CreatedAt, &guild.MemberCount,
		&guild.Leader.CharID, &guild.Leader.Name,
	)

	guild.Leader.GuildID = uint32(guild.ID)
	guild.Leader.JoinedAt = guild.CreatedAt

	if err != nil {
		s.logger.Error("failed to retrieve guild data from database")
		return nil, err
	}

	return guild, nil
}

func GetCharacterGuildData(s *Session, charID uint32) (*GuildMember, error) {
	result, err := s.server.db.Query(`
		SELECT guild_id, joined_at, name
			FROM guild_characters gc
				JOIN characters c on gc.character_id = c.id
			WHERE character_id=$1
		LIMIT 1
	`, charID)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to retrieve membership data for character '%d'", charID))
		return nil, err
	}

	hasRow := result.Next()

	if !hasRow {
		return nil, nil
	}

	memberData := &GuildMember{
		CharID: charID,
	}

	err = result.Scan(&memberData.GuildID, &memberData.JoinedAt, &memberData.Name)

	if err != nil {
		s.logger.Error("failed to retrieve guild data from database")
		return nil, err
	}

	return memberData, nil
}
