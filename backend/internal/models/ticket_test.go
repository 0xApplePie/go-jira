package models

import (
	"testing"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Status
		wantErr bool
	}{
		{
			name:    "valid TODO status",
			input:   "TODO",
			want:    StatusTodo,
			wantErr: false,
		},
		{
			name:    "valid PROGRESS status",
			input:   "PROGRESS",
			want:    StatusProgress,
			wantErr: false,
		},
		{
			name:    "valid DONE status",
			input:   "DONE",
			want:    StatusDone,
			wantErr: false,
		},
		{
			name:    "invalid status",
			input:   "INVALID",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStatus() = %v, want %v", got, tt.want)
			}
		})
	}
} 