package logic

import (
	"fmt"
	"progresstracker/data"
	"strconv"
	"time"
)

type Tracker struct {
	repo *data.Repository
}

func NewTracker(repo *data.Repository) *Tracker {
	return &Tracker{repo: repo}
}

func (t *Tracker) AddEntry(exercise, weightStr, repsStr, setsStr, notes, date string) (*data.Entry, error) {
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil || weight < 0 {
		return nil, fmt.Errorf("invalid weight")
	}
	reps, err := strconv.Atoi(repsStr)
	if err != nil || reps <= 0 {
		return nil, fmt.Errorf("invalid reps")
	}
	sets, err := strconv.Atoi(setsStr)
	if err != nil || sets <= 0 {
		return nil, fmt.Errorf("invalid sets")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	e := &data.Entry{
		Exercise: exercise,
		Weight:   weight,
		Reps:     reps,
		Sets:     sets,
		Notes:    notes,
		Date:     date,
	}
	if err := t.repo.Save(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (t *Tracker) GetHistory(exercise string) ([]*data.Entry, error) {
	return t.repo.HistoryFor(exercise)
}

func (t *Tracker) GetPersonalBest(exercise string) (*data.PersonalBest, error) {
	return t.repo.PersonalBest(exercise)
}

func (t *Tracker) GetLastEntry(exercise string) (*data.Entry, error) {
	return t.repo.LastEntry(exercise)
}
