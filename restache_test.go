package restache_test

import (
	"strings"
	"testing"

	"github.com/tetsuo/restache"
)

func TestTranspileFileInvalidInputFile(t *testing.T) {
	for _, tt := range []struct {
		desc      string
		inputFile string
		wantErr   string
	}{
		{
			"empty input path",
			"",
			"input file path is empty",
		},
		{
			"filename with invalid char",
			"/components/bad$name.stache",
			"contains invalid character '$'",
		},
		{
			"empty stem",
			"/components/.stache",
			"is not valid",
		},
		{
			"starts with symbol",
			"/components/_hidden.stache",
			"must start with a letter",
		},
		{
			"input is root slash",
			"/",
			"must start with a letter",
		},
		{
			"input is series of slashes",
			"///",
			"must start with a letter",
		},
		{
			"input is '.' (current directory)",
			".",
			": is a directory", // handled during parse()
		},
		{
			"input is '..' (parent directory)",
			"..",
			": is a directory", // handled during parse()
		},
		{
			"input doesn't exist",
			"/i/am/not/exist.stache",
			"no such file or directory",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			err := restache.TranspileFile(tt.inputFile, "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("got err = %v, want err containing %q", err, tt.wantErr)
			}
		})
	}
}
