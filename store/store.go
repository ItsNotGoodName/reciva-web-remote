package store

import (
	"database/sql"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	Presets []string
	db      *sql.DB
}

func NewStore(cfg *config.Config) (*Store, error) {
	db, err := sql.Open("sqlite3", cfg.DB)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS preset (
			uri TEXT PRIMARY KEY UNIQUE,
			sid INTEGER DEFAULT 0,
			FOREIGN KEY(sid) REFERENCES stream(id)
		)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `stream` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `name` TEXT UNIQUE, `content` TEXT)")
	if err != nil {
		return nil, err
	}

	var pts []string
	for _, uri := range cfg.URIS {
		pt := Preset{URI: uri}
		err = db.QueryRow("SELECT sid FROM preset WHERE uri = ?", uri).Scan(&pt.SID)
		if err != nil {
			_, err = db.Exec("INSERT INTO preset (uri) VALUES (?)", uri)
			if err != nil {
				return nil, err
			}
			err = db.QueryRow("SELECT sid FROM preset WHERE uri = ?", uri).Scan(&pt.SID)
			if err != nil {
				return nil, err
			}
		}
		pts = append(pts, pt.URI)
	}

	return &Store{db: db, Presets: pts}, nil
}
