package api

import (
	"testing"
	"time"
)

func TestNextDate(t *testing.T) {
	now := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		start   string
		repeat  string
		want    string
		wantErr bool
	}{
		{
			name:   "days",
			start:  "20240113",
			repeat: "d 7",
			want:   "20240127",
		},
		{
			name:   "year and leap day",
			start:  "20240229",
			repeat: "y",
			want:   "20250301",
		},
		{
			name:   "month days",
			start:  "20240116",
			repeat: "m 16,5",
			want:   "20240205",
		},
		{
			name:   "last and eighteenth day",
			start:  "20240201",
			repeat: "m -1,18",
			want:   "20240218",
		},
		{
			name:   "weekday",
			start:  "20240126",
			repeat: "w 1",
			want:   "20240129",
		},
		{
			name:    "empty rule",
			start:   "20240126",
			repeat:  "",
			wantErr: true,
		},
		{
			name:    "invalid interval",
			start:   "20240126",
			repeat:  "d 405",
			wantErr: true,
		},
		{
			name:    "invalid weekday",
			start:   "20240126",
			repeat:  "w 8",
			wantErr: true,
		},
		{
			name:    "invalid month day",
			start:   "20240126",
			repeat:  "m 4,40",
			wantErr: true,
		},
		{
			name:    "invalid month",
			start:   "20240126",
			repeat:  "m 5 1,13",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NextDate(now, test.start, test.repeat)

			if test.wantErr {
				if err == nil {
					t.Fatalf("ожидалась ошибка, получена дата %q", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			if got != test.want {
				t.Fatalf("получено %q, ожидалось %q", got, test.want)
			}
		})
	}
}
