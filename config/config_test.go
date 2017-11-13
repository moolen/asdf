package config

import (
	"fmt"
	"io"
	"io/ioutil"
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
				Types:         DefaultTypeMap,
			},
			err: "%!s(<nil>)",
		},
		{
			json: strings.NewReader("{\"types\":{\"foo\":\"bar\"}}"),
			conf: &Config{
				VersionFile:   DefaultVersionFile,
				ChangelogFile: DefaultChangelogFile,
				Types: map[string]string{
					"foo": "bar",
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
				Types: DefaultTypeMap,
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
				Types:     DefaultTypeMap,
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

func TestFromFile(t *testing.T) {
	table := []struct {
		content string
		config  *Config
		err     string
	}{
		{
			content: "",
			config:  nil,
			err:     "unexpected end of JSON input",
		},
		{
			content: "{}",
			config: &Config{
				VersionFile:   DefaultVersionFile,
				ChangelogFile: DefaultChangelogFile,
				Types:         DefaultTypeMap,
			},
			err: "%!s(<nil>)",
		},
	}
	for i, row := range table {
		tmp, _ := ioutil.TempFile("", "")
		tmp.WriteString(row.content)
		res, err := FromFile(tmp.Name())
		if !reflect.DeepEqual(row.config, res) {
			t.Fatalf("[%d] expected \n%#v, got \n%#v", i, row.config, res)
		}
		if fmt.Sprintf("%s", err) != row.err {
			t.Fatalf("[%d] expected \n[%s], got \n[%s]", i, row.err, err)
		}
	}
}
