package logic

import (
	"progresstracker/data"
	"sort"
)

type ChartPoint struct {
	Date   string
	Value  float64
}

type Analytics struct {
	repo *data.Repository
}

func NewAnalytics(repo *data.Repository) *Analytics {
	return &Analytics{repo: repo}
}

func (a *Analytics) WeightOverTime(exercise string) ([]ChartPoint, error) {
	entries, err := a.repo.HistoryFor(exercise)
	if err != nil {
		return nil, err
	}
	// deduplicate by date, take max weight per date
	byDate := map[string]float64{}
	for _, e := range entries {
		if e.Weight > byDate[e.Date] {
			byDate[e.Date] = e.Weight
		}
	}
	return sortedPoints(byDate), nil
}

func (a *Analytics) VolumeOverTime(exercise string) ([]ChartPoint, error) {
	entries, err := a.repo.HistoryFor(exercise)
	if err != nil {
		return nil, err
	}
	byDate := map[string]float64{}
	for _, e := range entries {
		byDate[e.Date] += e.Volume
	}
	return sortedPoints(byDate), nil
}

func sortedPoints(m map[string]float64) []ChartPoint {
	var pts []ChartPoint
	for k, v := range m {
		pts = append(pts, ChartPoint{Date: k, Value: v})
	}
	sort.Slice(pts, func(i, j int) bool {
		return pts[i].Date < pts[j].Date
	})
	return pts
}
