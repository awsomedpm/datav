package storage

import (
	"errors"

	"github.com/xObserve/xObserve/query/pkg/db"
	"github.com/xObserve/xObserve/query/pkg/e"
)

/* update table structure to current xobserve version */
func update() error {
	var visibleTo string
	err := db.Conn.QueryRow("SELECT visible_to FROM dashboard limit 1").Scan(&visibleTo)
	if err != nil && e.IsErrNoColumn(err) {
		_, err = db.Conn.Exec("ALTER TABLE dashboard ADD COLUMN visible_to VARCHAR(32) DEFAULT 'team'")
		if err != nil {
			return errors.New("update storage error:" + err.Error())
		}
	}

	var isPublic bool
	err = db.Conn.QueryRow("SELECT is_public FROM team limit 1").Scan(&isPublic)
	if err != nil && e.IsErrNoColumn(err) {
		_, err = db.Conn.Exec("ALTER TABLE team ADD COLUMN is_public BOOL DEFAULT false")
		if err != nil {
			return errors.New("update storage error:" + err.Error())
		}
	}

	// default_selected VARCHAR(255),
	var defaultSelected bool
	err = db.Conn.QueryRow("SELECT default_selected FROM variable limit 1").Scan(&defaultSelected)
	if err != nil && e.IsErrNoColumn(err) {
		_, err = db.Conn.Exec("ALTER TABLE variable ADD COLUMN default_selected VARCHAR(255)")
		if err != nil {
			return errors.New("update storage error:" + err.Error())
		}
	}

	// team_id INTEGER NOT NULL
	var teamId string
	err = db.Conn.QueryRow("SELECT team_id FROM variable limit 1").Scan(&teamId)
	if err != nil && e.IsErrNoColumn(err) {
		_, err = db.Conn.Exec("ALTER TABLE variable ADD COLUMN team_id INTEGER DEFAULT 1")
		if err != nil {
			return errors.New("update storage error:" + err.Error())
		}
	}

	err = db.Conn.QueryRow("SELECT team_id FROM datasource limit 1").Scan(&teamId)
	if err != nil && e.IsErrNoColumn(err) {
		_, err = db.Conn.Exec("ALTER TABLE datasource ADD COLUMN team_id INTEGER DEFAULT 1")
		if err != nil {
			return errors.New("update storage error:" + err.Error())
		}
	}

	// team allow_global
	var allowGlobal bool
	err = db.Conn.QueryRow("SELECT allow_global FROM team limit 1").Scan(&allowGlobal)
	if err != nil && e.IsErrNoColumn(err) {
		_, err = db.Conn.Exec("ALTER TABLE team ADD COLUMN allow_global BOOL DEFAULT true")
		if err != nil {
			return errors.New("update storage error:" + err.Error())
		}
	}
	return nil
}