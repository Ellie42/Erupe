package channelserver

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Andoryuuta/Erupe/common/stringsupport"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type FestivalColour string

const (
	FestivalColourNone FestivalColour = "none"
	FestivalColourRed  FestivalColour = "red"
	FestivalColourBlue FestivalColour = "blue"
)

var FestivalColourCodes = map[FestivalColour]uint8{
	FestivalColourBlue: 0x00,
	FestivalColourRed:  0x01,
	FestivalColourNone: 0xFF,
}

type GuildApplicationType string

const (
	GuildApplicationTypeApplied GuildApplicationType = "applied"
	GuildApplicationTypeInvited GuildApplicationType = "invited"
)

type Guild struct {
	ID             uint32         `db:"id"`
	Name           string         `db:"name"`
	MainMotto      uint8          `db:"main_motto"`
	SubMotto       uint8          `db:"sub_motto"`
	CreatedAt      time.Time      `db:"created_at"`
	MemberCount    uint16         `db:"member_count"`
	RP             uint32         `db:"rp"`
	Comment        string         `db:"comment"`
	FestivalColour FestivalColour `db:"festival_colour"`
	GuildHallType  uint16         `db:"guild_hall"`
	Icon           *GuildIcon     `db:"icon"`

	GuildLeader
}

type GuildLeader struct {
	LeaderCharID uint32 `db:"leader_id"`
	LeaderName   string `db:"leader_name"`
}

type GuildIconPart struct {
	Index    uint16
	ID       uint16
	Size     uint8
	Rotation uint8
	PosX     uint16
	PosY     uint16
}

type GuildApplication struct {
	ID              int                  `db:"id"`
	GuildID         uint32               `db:"guild_id"`
	CharID          uint32               `db:"character_id"`
	ActorID         uint32               `db:"actor_id"`
	ApplicationType GuildApplicationType `db:"application_type"`
	CreatedAt       time.Time            `db:"created_at"`
}

type GuildIcon struct {
	Parts []GuildIconPart
}

func (gi *GuildIcon) Scan(val interface{}) (err error) {
	switch v := val.(type) {
	case []byte:
		err = json.Unmarshal(v, &gi)
	case string:
		err = json.Unmarshal([]byte(v), &gi)
	}

	return
}

func (gi *GuildIcon) Value() (valuer driver.Value, err error) {
	return json.Marshal(gi)
}

const guildInfoSelectQuery = `
SELECT g.id,
       g.name,
       g.rp,
	   g.main_motto,
	   g.sub_motto,
       created_at,
       leader_id,
       lc.name as leader_name,
       comment,
       festival_colour,
	   guild_hall,
       icon,
       (
           SELECT count(1) FROM guild_characters gc WHERE gc.guild_id = g.id
       )             AS member_count
FROM guilds g
         JOIN guild_characters lgc ON lgc.character_id = leader_id
         JOIN characters lc on leader_id = lc.id
`

func (guild *Guild) Save(s *Session) error {
	_, err := s.server.db.Exec(`
		UPDATE guilds SET main_motto=$1, sub_motto=$6, comment=$3, festival_colour=$4, icon=$5 WHERE id=$2
	`, guild.MainMotto, guild.ID, guild.Comment, guild.FestivalColour, guild.Icon, guild.SubMotto)

	if err != nil {
		s.logger.Error("failed to update guild data", zap.Error(err), zap.Uint32("guildID", guild.ID))
		return err
	}

	return nil
}

func (guild *Guild) CreateApplication(s *Session, charID uint32, applicationType GuildApplicationType, transaction *sql.Tx) error {

	sql := `
		INSERT INTO guild_applications (guild_id, character_id, actor_id, application_type) 
		VALUES ($1, $2, $3, $4)
	`

	var err error

	if transaction == nil {
		_, err = s.server.db.Exec(sql, guild.ID, charID, s.charID, applicationType)
	} else {
		_, err = transaction.Exec(sql, guild.ID, charID, s.charID, applicationType)
	}

	if err != nil {
		s.logger.Error(
			"failed to add guild application",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", charID),
		)
		return err
	}

	return nil
}

func (guild *Guild) Disband(s *Session) error {
	transaction, err := s.server.db.Begin()

	if err != nil {
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	_, err = transaction.Exec("DELETE FROM guild_characters WHERE guild_id = $1", guild.ID)

	if err != nil {
		s.logger.Error("failed to remove guild characters", zap.Error(err), zap.Uint32("guildId", guild.ID))
		rollbackTransaction(s, transaction)
		return err
	}

	_, err = transaction.Exec("DELETE FROM guilds WHERE id = $1", guild.ID)

	if err != nil {
		s.logger.Error("failed to remove guild", zap.Error(err), zap.Uint32("guildID", guild.ID))
		rollbackTransaction(s, transaction)
		return err
	}

	err = transaction.Commit()

	if err != nil {
		s.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	s.logger.Info("Character disbanded guild", zap.Uint32("charID", s.charID), zap.Uint32("guildID", guild.ID))

	return nil
}

func (guild *Guild) RemoveCharacter(s *Session, charID uint32) error {
	_, err := s.server.db.Exec("DELETE FROM guild_characters WHERE character_id=$1", charID)

	if err != nil {
		s.logger.Error(
			"failed to remove character from guild",
			zap.Error(err),
			zap.Uint32("charID", charID),
			zap.Uint32("guildID", guild.ID),
		)

		return err
	}

	return nil
}

func (guild *Guild) AcceptApplication(s *Session, charID uint32) error {
	transaction, err := s.server.db.Begin()

	if err != nil {
		s.logger.Error("failed to start db transaction", zap.Error(err))
		return err
	}

	_, err = transaction.Exec(`DELETE FROM guild_applications WHERE character_id = $1`, charID)

	if err != nil {
		s.logger.Error("failed to accept character's guild application", zap.Error(err))
		rollbackTransaction(s, transaction)
		return err
	}

	_, err = transaction.Exec(`
		INSERT INTO guild_characters (guild_id, character_id, order_index)
		VALUES ($1, $2, (SELECT MAX(order_index) + 1 FROM guild_characters WHERE guild_id = $1))
	`, guild.ID, charID)

	if err != nil {
		s.logger.Error(
			"failed to add applicant to guild",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", charID),
		)
		rollbackTransaction(s, transaction)
		return err
	}

	err = transaction.Commit()

	if err != nil {
		s.logger.Error("failed to commit db transaction", zap.Error(err))
		rollbackTransaction(s, transaction)
		return err
	}

	return nil
}

// This is relying on the fact that invitation ID is also character ID right now
// if invitation ID changes, this will break.
func (guild *Guild) CancelInvitation(s *Session, charID uint32) error {
	_, err := s.server.db.Exec(
		`DELETE FROM guild_applications WHERE character_id = $1 AND guild_id = $2 AND application_type = 'invited'`,
		charID, guild.ID,
	)

	if err != nil {
		s.logger.Error(
			"failed to cancel guild invitation",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", charID),
		)
		return err
	}

	return nil
}

func (guild *Guild) RejectApplication(s *Session, charID uint32) error {
	_, err := s.server.db.Exec(
		`DELETE FROM guild_applications WHERE character_id = $1 AND guild_id = $2 AND application_type = 'applied'`,
		charID, guild.ID,
	)

	if err != nil {
		s.logger.Error(
			"failed to reject guild application",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", charID),
		)
		return err
	}

	return nil
}

func (guild *Guild) ArrangeCharacters(s *Session, charIDs []uint32) error {
	transaction, err := s.server.db.Begin()

	if err != nil {
		s.logger.Error("failed to start db transaction", zap.Error(err))
		return err
	}

	for i, id := range charIDs {
		_, err := transaction.Exec("UPDATE guild_characters SET order_index = $1 WHERE character_id = $2", 2+i, id)

		if err != nil {
			err = transaction.Rollback()

			if err != nil {
				s.logger.Error("failed to rollback db transaction", zap.Error(err))
			}

			return err
		}
	}

	err = transaction.Commit()

	if err != nil {
		s.logger.Error("failed to commit db transaction", zap.Error(err))
		return err
	}

	return nil
}

func (guild *Guild) DonateRP(s *Session, rp uint16, transaction *sql.Tx) (err error) {
	updateSQL := "UPDATE guilds SET rp = rp + $1 WHERE id = $2"

	if transaction != nil {
		_, err = transaction.Exec(updateSQL, rp, guild.ID)
	} else {
		_, err = s.server.db.Exec(updateSQL, rp, guild.ID)
	}

	if err != nil {
		s.logger.Error(
			"failed to donate RP to guild",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
		)

		return err
	}

	return nil
}

func (guild *Guild) GetApplicationForCharID(s *Session, charID uint32, applicationType GuildApplicationType) (*GuildApplication, error) {
	row := s.server.db.QueryRowx(`
		SELECT * from guild_applications WHERE character_id = $1 AND guild_id = $2 AND application_type = $3
	`, charID, guild.ID, applicationType)

	application := &GuildApplication{}

	err := row.StructScan(application)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		s.logger.Error(
			"failed to retrieve guild application for character",
			zap.Error(err),
			zap.Uint32("charID", charID),
			zap.Uint32("guildID", guild.ID),
		)
		return nil, err
	}

	return application, nil
}

func (guild *Guild) HasApplicationForCharID(s *Session, charID uint32) (bool, error) {
	row := s.server.db.QueryRowx(`
		SELECT 1 from guild_applications WHERE character_id = $1 AND guild_id = $2
	`, charID, guild.ID)

	num := 0

	err := row.Scan(&num)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		s.logger.Error(
			"failed to retrieve guild applications for character",
			zap.Error(err),
			zap.Uint32("charID", charID),
			zap.Uint32("guildID", guild.ID),
		)
		return false, err
	}

	return true, nil
}

func CreateGuild(s *Session, guildName string) (int32, error) {
	transaction, err := s.server.db.Begin()

	if err != nil {
		s.logger.Error("failed to start db transaction", zap.Error(err))
		return 0, err
	}

	guildNameSafe, err := stringsupport.ConvertShiftJISToUTF8(guildName)

	if err != nil {
		panic(err)
	}

	guildResult, err := transaction.Query(
		"INSERT INTO guilds (name, leader_id) VALUES ($1, $2) RETURNING id",
		guildNameSafe, s.charID,
	)

	if err != nil {
		s.logger.Error("failed to create guild", zap.Error(err))
		rollbackTransaction(s, transaction)
		return 0, err
	}

	var guildId int32

	guildResult.Next()

	err = guildResult.Scan(&guildId)

	if err != nil {
		s.logger.Error("failed to retrieve guild ID", zap.Error(err))
		rollbackTransaction(s, transaction)
		return 0, err
	}

	err = guildResult.Close()

	if err != nil {
		s.logger.Error("failed to finalise query", zap.Error(err))
		rollbackTransaction(s, transaction)
		return 0, err
	}

	_, err = transaction.Exec(`
		INSERT INTO guild_characters (guild_id, character_id)
		VALUES ($1, $2)
	`, guildId, s.charID)

	if err != nil {
		s.logger.Error("failed to add character to guild", zap.Error(err))
		rollbackTransaction(s, transaction)
		return 0, err
	}

	err = transaction.Commit()

	if err != nil {
		s.logger.Error("failed to commit guild creation", zap.Error(err))
		return 0, err
	}

	return guildId, nil
}

func rollbackTransaction(s *Session, transaction *sql.Tx) {
	err := transaction.Rollback()

	if err != nil {
		s.logger.Error("failed to rollback transaction", zap.Error(err))
	}
}

func FindGuildsByName(s *Session, name string) ([]*Guild, error) {
	searchTerm := fmt.Sprintf("%%%s%%", name)

	rows, err := s.server.db.Queryx(fmt.Sprintf(`
		%s
		WHERE g.name ILIKE $1
	`, guildInfoSelectQuery), searchTerm)

	if err != nil {
		s.logger.Error("failed to find guilds for search term", zap.Error(err), zap.String("searchTerm", name))
		return nil, err
	}

	defer rows.Close()

	guilds := make([]*Guild, 0)

	for rows.Next() {
		guild, err := buildGuildObjectFromDbResult(rows, err, s)

		if err != nil {
			return nil, err
		}

		guilds = append(guilds, guild)
	}

	return guilds, nil
}

func GetGuildInfoByID(s *Session, guildID uint32) (*Guild, error) {
	rows, err := s.server.db.Queryx(fmt.Sprintf(`
		%s
		WHERE g.id = $1 
		LIMIT 1
	`, guildInfoSelectQuery), guildID)

	if err != nil {
		s.logger.Error("failed to retrieve guild", zap.Error(err), zap.Uint32("guildID", guildID))
		return nil, err
	}

	defer rows.Close()

	hasRow := rows.Next()

	if !hasRow {
		return nil, nil
	}

	return buildGuildObjectFromDbResult(rows, err, s)
}

func GetGuildInfoByCharacterId(s *Session, charID uint32) (*Guild, error) {
	rows, err := s.server.db.Queryx(fmt.Sprintf(`
		%s
		WHERE EXISTS(
				SELECT 1
				FROM guild_characters gc1
				WHERE gc1.character_id = $1
				  AND gc1.guild_id = g.id
			)
		   OR EXISTS(
				SELECT 1
				FROM guild_applications ga
				WHERE ga.character_id = $1
				  AND ga.guild_id = g.id
				  AND ga.application_type = 'applied'
			)
		LIMIT 1
	`, guildInfoSelectQuery), charID)

	if err != nil {
		s.logger.Error("failed to retrieve guild for character", zap.Error(err), zap.Uint32("charID", charID))
		return nil, err
	}

	defer rows.Close()

	hasRow := rows.Next()

	if !hasRow {
		return nil, nil
	}

	return buildGuildObjectFromDbResult(rows, err, s)
}

func buildGuildObjectFromDbResult(result *sqlx.Rows, err error, s *Session) (*Guild, error) {
	guild := &Guild{}

	err = result.StructScan(guild)

	if err != nil {
		s.logger.Error("failed to retrieve guild data from database", zap.Error(err))
		return nil, err
	}

	return guild, nil
}
