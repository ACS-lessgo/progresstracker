package data

import "time"

type Entry struct {
	ID        int64
	Exercise  string
	Weight    float64
	Reps      int
	Sets      int
	Volume    float64
	Notes     string
	Date      string
	CreatedAt time.Time
}

type PersonalBest struct {
	Exercise  string
	MaxWeight float64
	MaxVolume float64
	Date      string
}

type WeekPlan struct {
	Week1 []string
	Week2 []string
}

var WorkoutPlans = map[string]WeekPlan{
	"Monday": {
		Week1: []string{
			"Flat Bench Barbell Chest Press",
			"Inclined Dumbbell Press",
			"Seated Pec Dec Flies Machine",
			"Cable Flies (Low to High)",
			"Close-Grip Dumbbell Press",
		},
		Week2: []string{
			"Incline Barbell Bench Press",
			"Decline Dumbbell Press",
			"Flat Bench Cable Flies",
			"Standing Cable Crossover (High to Low)",
			"Dumbbell Pullover",
		},
	},
	"Tuesday": {
		Week1: []string{
			"Wide Grip Lat Pulldown",
			"Seated V Bar Cable Rowing",
			"Dumbbell Rowing",
			"Close Grip Lat Pulldown",
			"Lat Pushdown",
		},
		Week2: []string{
			"Mid Grip Lat Pulldown",
			"T-Bar Row",
			"Single Arm Cable Rowing",
			"Reverse Cable Crossovers",
			"Hyperextensions",
		},
	},
	"Wednesday": {
		Week1: []string{
			"Seated Dumbbell Press",
			"Dumbbell Lateral Raises",
			"Dumbbell Alternate Front Raises",
			"Upright Rows",
			"Shrugs",
			"Crunches",
			"Russian Twists",
		},
		Week2: []string{
			"Seated Overhead Barbell Press",
			"Cable Lateral Raises",
			"Front Plate Raises",
			"Rope Face Pulls",
			"Smith Machine Shrugs",
			"Hanging Leg Raises",
			"Cable Crunches",
		},
	},
	"Thursday": {
		Week1: []string{
			"Standing Alternate Bicep Curl",
			"Seated Single Arm Tricep Extensions",
			"Standing Alternate Hammer Curls",
			"Cable Rope Pushdown",
			"Reverse Grip Barbell Curl",
			"Tricep Cable Kickbacks",
			"Leg Raises",
			"Plank",
		},
		Week2: []string{
			"Barbell Curl",
			"Overhead Dumbbell Tricep Extensions",
			"Preacher Curl",
			"Straight Bar Pushdown",
			"Zottman Curl",
			"Overhead Rope Extensions",
			"Side Plank",
			"Toe Touches",
		},
	},
	"Friday": {
		Week1: []string{
			"Smith Machine Squats",
			"Leg Extensions",
			"Walking Lunges",
			"Leg Press",
			"Hamstring Curls",
			"Standing Calf Raises",
		},
		Week2: []string{
			"Barbell Squats",
			"Bulgarian Split Squats",
			"Standing Lunges",
			"Hack Squats",
			"Romanian Deadlifts",
			"Seated Calf Raises",
		},
	},
}

var DayOrder = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

// WorkoutDays returns the exercise list for a given day and week (1 or 2)
func WorkoutDays(day string, week int) []string {
	plan, ok := WorkoutPlans[day]
	if !ok {
		return nil
	}
	if week == 2 {
		return plan.Week2
	}
	return plan.Week1
}
