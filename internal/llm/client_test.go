package llm

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "openai provider",
			config: Config{
				Provider: "openai",
				APIKey:   "test-key",
			},
			wantErr: false,
		},
		{
			name: "default provider (empty)",
			config: Config{
				Provider: "",
				APIKey:   "test-key",
			},
			wantErr: false,
		},
		{
			name: "unsupported provider",
			config: Config{
				Provider: "unsupported",
				APIKey:   "test-key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() client = nil, want non-nil")
			}
		})
	}
}


