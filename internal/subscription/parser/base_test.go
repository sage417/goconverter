package parser

import "testing"

func TestParseSubscription(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
		wantErr bool
	}{
		{
			name:    "empty content",
			content: "",
			want:    0,
			wantErr: true,
		},
		// 添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSubscription(tt.content, "clashx")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ParseSubscription() got %v nodes, want %v", len(got), tt.want)
			}
		})
	}
}
