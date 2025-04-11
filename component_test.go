package restache

import (
	"fmt"
	"path/filepath"
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

func TestNewComponentEntry(t *testing.T) {
	tests := []struct {
		desc     string
		dir      string
		input    string
		wantTag  string
		wantStem string
		wantExt  string
		wantErr  string
	}{
		{
			"valid filename with uppercase",
			"/components",
			"Foo.stache",
			"foo",
			"Foo",
			".stache",
			"",
		},
		{
			"valid filename lowercase",
			"/components",
			"bar.stache",
			"bar",
			"bar",
			".stache",
			"",
		},
		{
			"no extension",
			"/components",
			"Widget",
			"widget",
			"Widget",
			"",
			"",
		},
		{
			"filename with invalid char",
			"/components",
			"bad$name.stache",
			"",
			"",
			"",
			"contains invalid character '$'",
		},
		{
			"filename with slash",
			"/components",
			"dir/file.stache",
			"",
			"",
			"",
			"must be a filename",
		},
		{
			"empty stem",
			"/components",
			".stache",
			"",
			"",
			"",
			"is not valid",
		},
		{
			"starts with symbol",
			"/components",
			"_hidden.stache",
			"",
			"",
			"",
			"must start with a letter",
		},
		{
			"empty input",
			"/components",
			"",
			"",
			"",
			"",
			"must be a filename",
		},
		{
			"input is '.'",
			"/components",
			".",
			"",
			"",
			"",
			"is not valid",
		},
		{
			"input is '..'",
			"/components",
			"..",
			"",
			"",
			"",
			"must start with a letter",
		},
		{
			"input is root slash",
			"/components",
			"/",
			"",
			"",
			"",
			"must start with a letter",
		},
		{
			"input is empty",
			"/components",
			"",
			"",
			"",
			"",
			"must be a filename",
		},
		{
			"input is series of slashes",
			"/components",
			"///",
			"",
			"",
			"",
			"must be a filename",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			entry, err := newComponentEntry(tt.dir, tt.input)

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("got err = %v, want err containing %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			wantPath := filepath.Join(tt.dir, tt.input)
			if entry.path != wantPath {
				t.Errorf("got path = %q, want %q", entry.path, wantPath)
			}
			if entry.tag != tt.wantTag {
				t.Errorf("got tag = %q, want %q", entry.tag, tt.wantTag)
			}
			if entry.stem != tt.wantStem {
				t.Errorf("got stem = %q, want %q", entry.stem, tt.wantStem)
			}
			if entry.ext != tt.wantExt {
				t.Errorf("got ext = %q, want %q", entry.ext, tt.wantExt)
			}
		})
	}
}
