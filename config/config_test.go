package config

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	table := []struct {
		json io.Reader
		conf *Config
		err  string
	}{
		{
			json: strings.NewReader(""),
			err:  "unexpected end of JSON input",
		},
		{
			json: strings.NewReader("{}"),
			conf: &Config{
				VersionFile:   DefaultVersionFile,
				ChangelogFile: DefaultChangelogFile,
				Types:         defaultTypes(),
			},
			err: "%!s(<nil>)",
		},
		{
			json: strings.NewReader("{\"types\":[{\"key\":\"foo\",\"label\":\"bar\",\"major\":true}]}"),
			conf: &Config{
				VersionFile:   DefaultVersionFile,
				ChangelogFile: DefaultChangelogFile,
				Types: TypeConstraints{
					&TypeConstraint{
						Key:   "foo",
						Label: "bar",
						Major: true,
					},
				},
			},
			err: "%!s(<nil>)",
		},
		{
			json: strings.NewReader("{ \"branch_suffix\": {\"foo\":\"bar\"} }"),
			conf: &Config{
				VersionFile:   DefaultVersionFile,
				ChangelogFile: DefaultChangelogFile,
				BranchSuffix: map[string]string{
					"foo": "bar",
				},
				Types: defaultTypes(),
			},
			err: "%!s(<nil>)",
		},
		{
			json: strings.NewReader("{\"version_file\":\"myversion\",\"changelog_file\":\"myfile\",\"branch_suffix\":{\"foo\":\"bar\"},\"ticket_url\":\"htp\"}"),
			conf: &Config{
				VersionFile:   "myversion",
				ChangelogFile: "myfile",
				BranchSuffix: map[string]string{
					"foo": "bar",
				},
				Types:     defaultTypes(),
				TicketURL: "htp",
			},
			err: "%!s(<nil>)",
		},
	}
	for i, row := range table {
		result, err := FromJSON(row.json)
		if fmt.Sprintf("%s", err) != row.err {
			t.Fatalf("[%d] expected \n[%s], got \n[%s]", i, row.err, err)
		}
		if !reflect.DeepEqual(result, row.conf) {
			t.Fatalf("[%d]\n[%#v]\n[%#v]\n", i, row, result)
		}
	}
}
