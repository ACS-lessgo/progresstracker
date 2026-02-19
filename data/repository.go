package data

type Repository struct {
	db *DB
}

func NewRepository(db *DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(e *Entry) error {
	return r.db.InsertEntry(e)
}

func (r *Repository) HistoryFor(exercise string) ([]*Entry, error) {
	return r.db.GetEntriesByExercise(exercise)
}

func (r *Repository) All() ([]*Entry, error) {
	return r.db.GetAllEntries()
}

func (r *Repository) PersonalBest(exercise string) (*PersonalBest, error) {
	return r.db.GetPersonalBest(exercise)
}

func (r *Repository) LastEntry(exercise string) (*Entry, error) {
	return r.db.GetLastEntry(exercise)
}

func (r *Repository) GetCurrentWeek() int {
	return r.db.GetCurrentWeek()
}

func (r *Repository) SetCurrentWeek(week int) error {
	return r.db.SetCurrentWeek(week)
}
