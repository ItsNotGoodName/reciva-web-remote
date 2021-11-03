package store

import (
	"database/sql"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	Presets map[string]Preset // Presets is a map of preset uri to Preset
	db      *sql.DB
}

func NewStore(cfg *config.Config) (*Store, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `preset` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `uri` TEXT UNIQUE, `sid` INTEGER DEFAULT 0)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `stream` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `name` TEXT UNIQUE, `content` TEXT)")
	if err != nil {
		return nil, err
	}

	pts := make(map[string]Preset)
	for _, uri := range cfg.URIS {
		pt := Preset{URI: uri}
		err = db.QueryRow("SELECT id, sid FROM preset WHERE uri = ?", uri).Scan(&pt.ID, &pt.SID)
		if err != nil {
			_, err = db.Exec("INSERT INTO preset (uri) VALUES (?)", uri)
			if err != nil {
				return nil, err
			}
			err = db.QueryRow("SELECT id, sid FROM preset WHERE uri = ?", uri).Scan(&pt.ID, &pt.SID)
			if err != nil {
				return nil, err
			}
		}
		pts[pt.URI] = pt
	}

	return &Store{db: db, Presets: pts}, nil
}
