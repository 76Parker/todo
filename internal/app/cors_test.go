package app

import "testing"

func TestIsTrustedOrigin(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{
			name:   "Localhost_1",
			origin: "http://localhost:5173",
			want:   true,
		},
		{
			name:   "Localhost_2",
			origin: "https://localhost:3000",
			want:   true,
		},
		{
			name:   "Loopback_1",
			origin: "http://127.0.0.1:8080",
			want:   true,
		},
		{
			name:   "InvalidOrigin_1",
			origin: "not-a-url",
			want:   false,
		},
		{
			name:   "UntrustedOrigin_1",
			origin: "http://evil.example.com",
			want:   false,
		},
		{
			name:   "EmptyOrigin_1",
			origin: "",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTrustedOrigin(tt.origin)
			if got != tt.want {
				t.Fatalf("isTrustedOrigin(%q) = %v, want %v", tt.origin, got, tt.want)
			}
		})
	}
}
