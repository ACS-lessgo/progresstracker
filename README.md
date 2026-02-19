# ProgressTracker

A native desktop workout tracking app built with Go + Gio (gioui.org) and SQLite.

## Prerequisites

- Go 1.21+ (https://go.dev/dl/)
- GCC / C compiler (required by go-sqlite3 CGo)
  - **Windows**: Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or MinGW-w64
  - **Linux**: `sudo apt install gcc`
  - **macOS**: `xcode-select --install`

## Build Instructions

```bash
# 1. Extract the archive
tar xzf progresstracker.tar.gz
cd progresstracker

# 2. Download dependencies
go mod tidy

# 3. Run directly
go run main.go

# 4. Build binary
go build -o progresstracker        # Linux/Mac
go build -o progresstracker.exe    # Windows
```

## Features

- **5 Workout Days** from your program (Monday–Friday)
  - Monday: Chest
  - Tuesday: Back
  - Wednesday: Shoulders + Abs
  - Thursday: Arms + Abs
  - Friday: Legs
- **Log entries**: weight, reps, sets, notes, date
- **Auto-calculated volume** (weight × reps × sets)
- **Personal Best tracking** (max weight + max volume per exercise)
- **Full history view** with PB highlighted in gold
- **Charts**: weight over time + volume over time (line charts)
- **Dark theme** throughout
- **SQLite persistence** in `progress.db`

## Project Structure

```
progresstracker/
├── main.go
├── data/
│   ├── models.go       # Exercise definitions, Entry struct
│   ├── db.go           # SQLite operations
│   └── repository.go   # Repository pattern
├── logic/
│   ├── tracker.go      # Business logic for entries
│   └── analytics.go    # Chart data computation
└── ui/
    ├── app.go          # Main app state, layout, event loop
    ├── theme.go        # Dark theme colors
    ├── components.go   # Reusable UI components
    └── charts.go       # Line chart rendering
```

## Database

SQLite file `progress.db` is created automatically in the working directory.

Schema:

```sql
CREATE TABLE entries (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    exercise   TEXT NOT NULL,
    weight     REAL NOT NULL,
    reps       INTEGER NOT NULL,
    sets       INTEGER NOT NULL,
    volume     REAL NOT NULL,
    notes      TEXT,
    date       TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Screenshots:

![HomePage] (screenshots/s1.png)
![TrackerPage] (screenshots/s2.png)
