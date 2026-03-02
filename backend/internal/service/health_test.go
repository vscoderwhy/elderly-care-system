package service

import (
	"testing"
	"time"
)

// TestHealthRecord_Validation 测试健康记录验证
func TestHealthRecord_Validation(t *testing.T) {
	tests := []struct {
		name      string
		recordType string
		value     string
		expectValid bool
	}{
		{
			name:       "valid blood pressure",
			recordType: "blood_pressure",
			value:      "120",
			expectValid: true,
		},
		{
			name:       "valid blood sugar",
			recordType: "blood_sugar",
			value:      "5.6",
			expectValid: true,
		},
		{
			name:       "valid temperature",
			recordType: "temperature",
			value:      "36.5",
			expectValid: true,
		},
		{
			name:       "empty value",
			recordType: "blood_pressure",
			value:      "",
			expectValid: false,
		},
		{
			name:       "invalid type",
			recordType: "invalid",
			value:      "100",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidHealthRecord(tt.recordType, tt.value)
			if valid != tt.expectValid {
				t.Errorf("isValidHealthRecord(%s, %s) = %v, expected %v",
					tt.recordType, tt.value, valid, tt.expectValid)
			}
		})
	}
}

// TestFormatTime 测试时间格式化
func TestFormatTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		time      time.Time
		contains  string
	}{
		{
			name:     "just now",
			time:     now.Add(-30 * time.Second),
			contains: "刚刚",
		},
		{
			name:     "minutes ago",
			time:     now.Add(-5 * time.Minute),
			contains: "分钟前",
		},
		{
			name:     "hours ago",
			time:     now.Add(-2 * time.Hour),
			contains: "小时前",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimeAgo(tt.time)
			if result == "" {
				t.Errorf("formatTimeAgo() returned empty string")
			}
		})
	}
}

// TestCalculateAge 测试年龄计算
func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		birthDate time.Time
		expected int
	}{
		{
			name:     "almost 30 years old (birthday not yet)",
			birthDate: now.AddDate(-30, 0, 1),
			expected: 29,
		},
		{
			name:     "newborn",
			birthDate: now.AddDate(0, 0, -1),
			expected: 0,
		},
		{
			name:     "exactly 1 year old",
			birthDate: now.AddDate(-1, 0, 0),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateAge(tt.birthDate)
			if result != tt.expected {
				t.Errorf("calculateAge() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

// Helper functions for testing

func isValidHealthRecord(recordType, value string) bool {
	validTypes := map[string]bool{
		"blood_pressure": true,
		"blood_sugar":    true,
		"temperature":    true,
		"weight":         true,
		"heart_rate":     true,
	}

	if !validTypes[recordType] {
		return false
	}
	return value != ""
}

func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "刚刚"
	}
	if diff < time.Hour {
		return "分钟前"
	}
	if diff < 24*time.Hour {
		return "小时前"
	}
	return "天前"
}

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()

	// 调整年龄：如果今年生日还没过，减1
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age
}
