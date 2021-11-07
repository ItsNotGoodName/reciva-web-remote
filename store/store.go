package store

import (
	"database/sql"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	_ "modernc.org/sqlite"
)

type Store struct {
	Presets []string
	db      *sql.DB
}

func NewStore(cfg *config.Config) (*Store, error) {
	db, err := sql.Open("sqlite", cfg.DB)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS preset (
			url TEXT PRIMARY KEY UNIQUE,
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
	for _, url := range cfg.URLS {
		pt := Preset{URL: url}
		err = db.QueryRow("SELECT sid FROM preset WHERE url = ?", url).Scan(&pt.SID)
		if err != nil {
			_, err = db.Exec("INSERT INTO preset (url) VALUES (?)", url)
			if err != nil {
				return nil, err
			}
			err = db.QueryRow("SELECT sid FROM preset WHERE url = ?", url).Scan(&pt.SID)
			if err != nil {
				return nil, err
			}
		}
		pts = append(pts, pt.URL)
	}

	return &Store{db: db, Presets: pts}, nil
}
