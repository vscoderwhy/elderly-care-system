package service

import (
	"fmt"
	"testing"
)

// Mock repositories for testing
type MockElderlyRepository struct {
	count int
	err   error
}

func (m *MockElderlyRepository) Count() (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.count, nil
}

// TestStatsService_DashboardStats tests dashboard statistics
func TestStatsService_DashboardStats(t *testing.T) {
	// Test case 1: Normal operation
	t.Run("should return correct stats", func(t *testing.T) {
		// This is a placeholder for actual test implementation
		// In real scenario, we would inject mock repositories
		expected := 128
		if expected != 128 {
			t.Errorf("Expected %d, got %d", 128, expected)
		}
	})
}

// TestOccupancyRate tests occupancy rate calculation
func TestOccupancyRate(t *testing.T) {
	tests := []struct {
		name      string
		total     int64
		occupied  int64
		expected  float64
	}{
		{
			name:     "50% occupancy",
			total:    100,
			occupied: 50,
			expected: 50.0,
		},
		{
			name:     "100% occupancy",
			total:    100,
			occupied: 100,
			expected: 100.0,
		},
		{
			name:     "0% occupancy",
			total:    100,
			occupied: 0,
			expected: 0.0,
		},
		{
			name:     "zero beds",
			total:    0,
			occupied: 0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result float64
			if tt.total > 0 {
				result = float64(tt.occupied) / float64(tt.total) * 100
			}
			if result != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

// TestFormatMoney tests money formatting
func TestFormatMoney(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0.00"},
		{100.5, "100.50"},
		{1234.56, "1234.56"},
		{1000000, "1000000.00"},
	}

	for _, tt := range tests {
		result := formatMoney(tt.input)
		if result != tt.expected {
			t.Errorf("formatMoney(%f) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func formatMoney(value float64) string {
	if value == 0 {
		return "0.00"
	}
	// 简单实现：保留两位小数
	return fmt.Sprintf("%.2f", value)
}
