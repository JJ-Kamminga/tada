package todo

import (
	"testing"
	"time"
)

func TestIsCompletedOlderThanDays(t *testing.T) {
	tests := []struct {
		name          string
		item          Item
		days          int
		expectedOlder bool
	}{
		{
			name: "completed item from 10 days ago should be older than 5 days",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: true,
		},
		{
			name: "completed item from 2 days ago should not be older than 5 days",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: false,
		},
		{
			name: "completed item from exactly 5 days ago should be older than 5 days (boundary)",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: true,
		},
		{
			name: "completed item from 6 days ago should be older than 5 days",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -6).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: true,
		},
		{
			name: "uncompleted item should return false",
			item: Item{
				Completed:      false,
				CompletionDate: time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			},
			days:          5,
			expectedOlder: false,
		},
		{
			name: "completed item without completion date should return false",
			item: Item{
				Completed:      true,
				CompletionDate: "",
			},
			days:          5,
			expectedOlder: false,
		},
		{
			name: "completed item with invalid date format should return false",
			item: Item{
				Completed:      true,
				CompletionDate: "invalid-date",
			},
			days:          5,
			expectedOlder: false,
		},
		{
			name: "completed item from 30 days ago should be older than 10 days",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -30).Format("2006-01-02"),
			},
			days:          10,
			expectedOlder: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsCompletedOlderThanDays(tt.days)
			if result != tt.expectedOlder {
				t.Errorf("IsCompletedOlderThanDays() = %v, want %v", result, tt.expectedOlder)
			}
		})
	}
}

func TestShouldBeVisible(t *testing.T) {
	tests := []struct {
		name            string
		item            Item
		expectedVisible bool
	}{
		{
			name: "uncompleted item should be visible",
			item: Item{
				Completed: false,
			},
			expectedVisible: true,
		},
		{
			name: "recently completed item should be visible",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			},
			expectedVisible: true,
		},
		{
			name: "completed item from 10 days ago should not be visible",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			},
			expectedVisible: false,
		},
		{
			name: "completed item from exactly 5 days ago should not be visible (boundary)",
			item: Item{
				Completed:      true,
				CompletionDate: time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
			},
			expectedVisible: false,
		},
		{
			name: "completed item without date should be visible",
			item: Item{
				Completed:      true,
				CompletionDate: "",
			},
			expectedVisible: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.ShouldBeVisible()
			if result != tt.expectedVisible {
				t.Errorf("ShouldBeVisible() = %v, want %v", result, tt.expectedVisible)
			}
		})
	}
}
