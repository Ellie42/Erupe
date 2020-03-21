package channelserver

import (
	"database/sql"
	"encoding/binary"
	"github.com/Andoryuuta/Erupe/server/channelserver/compression/nullcomp"
	"go.uber.org/zap"
)

const CharacterSaveRPPointer = 0x22D16

type CharacterSaveData struct {
	ID           uint32
	RP           uint16
	BaseSaveData []byte
}

func GetCharacterSaveData(s *Session, charID uint32) (*CharacterSaveData, error) {
	result, err := s.server.db.Queryx(
		"SELECT id, rp, savedata FROM characters WHERE id = $1",
		charID,
	)

	if err != nil {
		s.logger.Error(
			"failed to retrieve save data for character",
			zap.Error(err),
			zap.Uint32("charID", charID),
		)
		return nil, err
	}

	saveData := &CharacterSaveData{}
	var compressedBaseSave []byte

	if !result.Next() {
		s.logger.Error(
			"no results found for character save data",
			zap.Uint32("charID", charID),
		)
		return nil, err
	}

	err = result.Scan(&saveData.ID, &saveData.RP, &compressedBaseSave)

	if err != nil {
		s.logger.Error(
			"failed to retrieve save data for character",
			zap.Error(err),
			zap.Uint32("charID", charID),
		)

		return nil, err
	}

	decompressedBaseSave, err := nullcomp.Decompress(compressedBaseSave)

	if err != nil {
		s.logger.Error("Failed to decompress savedata from db", zap.Error(err))
		return nil, err
	}

	saveData.BaseSaveData = decompressedBaseSave

	return saveData, nil
}

func (save *CharacterSaveData) Save(s *Session, transaction *sql.Tx) error {
	rpBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(rpBytes, save.RP)
	copy(save.BaseSaveData[CharacterSaveRPPointer:CharacterSaveRPPointer+2], rpBytes)

	compressedData, err := nullcomp.Compress(save.BaseSaveData)

	if err != nil {
		s.logger.Error("failed to compress saveData", zap.Error(err), zap.Uint32("charID", save.ID))
		return err
	}

	updateSQL := `
		UPDATE characters 
			SET savedata=$1, rp=$2
		WHERE id=$3
	`

	if transaction != nil {
		_, err = transaction.Exec(updateSQL, compressedData, save.RP, save.ID)
	} else {
		_, err = s.server.db.Exec(updateSQL, compressedData, save.RP, save.ID)
	}

	if err != nil {
		s.logger.Error("failed to save character data", zap.Error(err), zap.Uint32("charID", save.ID))
		return err
	}

	return nil
}
