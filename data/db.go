package data

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func NewDB(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

// func (db *DB) migrate() error {
// 	_, err := db.conn.Exec(`
// 		CREATE TABLE IF NOT EXISTS entries (
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			exercise TEXT NOT NULL,
// 			weight REAL NOT NULL,
// 			reps INTEGER NOT NULL,
// 			sets INTEGER NOT NULL,
// 			volume REAL NOT NULL,
// 			notes TEXT,
// 			date TEXT NOT NULL,
// 			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
// 		)
// 	`)
// 	return err
// }

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) InsertEntry(e *Entry) error {
	e.Volume = e.Weight * float64(e.Reps) * float64(e.Sets)
	res, err := db.conn.Exec(
		`INSERT INTO entries (exercise, weight, reps, sets, volume, notes, date, created_at) VALUES (?,?,?,?,?,?,?,?)`,
		e.Exercise, e.Weight, e.Reps, e.Sets, e.Volume, e.Notes, e.Date, time.Now(),
	)
	if err != nil {
		return err
	}
	e.ID, _ = res.LastInsertId()
	return nil
}

func (db *DB) GetEntriesByExercise(exercise string) ([]*Entry, error) {
	rows, err := db.conn.Query(
		`SELECT id, exercise, weight, reps, sets, volume, notes, date, created_at FROM entries WHERE exercise=? ORDER BY date DESC, id DESC`,
		exercise,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []*Entry
	for rows.Next() {
		e := &Entry{}
		if err := rows.Scan(&e.ID, &e.Exercise, &e.Weight, &e.Reps, &e.Sets, &e.Volume, &e.Notes, &e.Date, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (db *DB) GetAllEntries() ([]*Entry, error) {
	rows, err := db.conn.Query(
		`SELECT id, exercise, weight, reps, sets, volume, notes, date, created_at FROM entries ORDER BY date DESC, id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []*Entry
	for rows.Next() {
		e := &Entry{}
		if err := rows.Scan(&e.ID, &e.Exercise, &e.Weight, &e.Reps, &e.Sets, &e.Volume, &e.Notes, &e.Date, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (db *DB) GetPersonalBest(exercise string) (*PersonalBest, error) {
	pb := &PersonalBest{Exercise: exercise}
	row := db.conn.QueryRow(
		`SELECT MAX(weight), MAX(volume) FROM entries WHERE exercise=?`,
		exercise,
	)
	var maxW, maxV sql.NullFloat64
	if err := row.Scan(&maxW, &maxV); err != nil {
		return nil, err
	}
	if maxW.Valid {
		pb.MaxWeight = maxW.Float64
	}
	if maxV.Valid {
		pb.MaxVolume = maxV.Float64
	}
	return pb, nil
}

func (db *DB) GetLastEntry(exercise string) (*Entry, error) {
	row := db.conn.QueryRow(
		`SELECT id, exercise, weight, reps, sets, volume, notes, date, created_at FROM entries WHERE exercise=? ORDER BY date DESC, id DESC LIMIT 1`,
		exercise,
	)
	e := &Entry{}
	err := row.Scan(&e.ID, &e.Exercise, &e.Weight, &e.Reps, &e.Sets, &e.Volume, &e.Notes, &e.Date, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			exercise TEXT NOT NULL,
			weight REAL NOT NULL,
			reps INTEGER NOT NULL,
			sets INTEGER NOT NULL,
			volume REAL NOT NULL,
			notes TEXT,
			date TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	// default to week 1
	_, err = db.conn.Exec(`
		INSERT OR IGNORE INTO settings (key, value) VALUES ('current_week', '1')
	`)
	return err
}

func (db *DB) GetCurrentWeek() int {
	row := db.conn.QueryRow(`SELECT value FROM settings WHERE key='current_week'`)
	var val string
	if err := row.Scan(&val); err != nil {
		return 1
	}
	if val == "2" {
		return 2
	}
	return 1
}

func (db *DB) SetCurrentWeek(week int) error {
	v := "1"
	if week == 2 {
		v = "2"
	}
	_, err := db.conn.Exec(`INSERT OR REPLACE INTO settings (key, value) VALUES ('current_week', ?)`, v)
	return err
}
