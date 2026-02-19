package main

import (
	"fmt"
	"os"

	"gioui.org/app"
	"progresstracker/data"
	"progresstracker/logic"
	"progresstracker/ui"
)

func main() {
	db, err := data.NewDB("progress.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open database:", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := data.NewRepository(db)
	tracker := logic.NewTracker(repo)
	anal := logic.NewAnalytics(repo)

	go func() {
		ui.Run(repo, tracker, anal)
		os.Exit(0)
	}()
	app.Main()
}
