package restache

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateTagName(t *testing.T) {
	tests := []struct {
		desc      string
		input     string
		wantUpper bool
		wantErr   string
	}{
		{
			"valid with uppercase letter",
			"tagName",
			true,
			"",
		},
		{
			"valid with multiple uppercase letters",
			"TagNAME123",
			true,
			"",
		},
		{
			"valid lowercase only",
			"tagname",
			false,
			"",
		},
		{
			"starts with number",
			"1tag",
			false,
			"must start with a letter",
		},
		{
			"starts with symbol",
			"_tag",
			false,
			"must start with a letter",
		},
		{
			"contains dash",
			"tag-name",
			false,
			"",
		},
		{
			"contains symbol",
			"tag$name",
			false,
			"contains invalid character '$'",
		},
		{
			"empty string (panic case)",
			"",
			false,
			"index out of range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.wantErr == "" || !strings.Contains(fmt.Sprint(r), tt.wantErr) {
						t.Errorf("panic = %v, wantErr = %q", r, tt.wantErr)
					}
				}
			}()

			ok, err := validateTagName(tt.input)

			if ok != tt.wantUpper {
				t.Errorf("got ok = %v, want %v", ok, tt.wantUpper)
			}
			if err != nil {
				if tt.wantErr == "" || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("got err = %v, wantErr to contain %q", err, tt.wantErr)
				}
			} else if tt.wantErr != "" {
				t.Errorf("got err = nil, wantErr %q", tt.wantErr)
			}
		})
	}
}
