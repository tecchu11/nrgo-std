package nrslog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/tecchu11/nrgo-std/nrslog"
)

var app *newrelic.Application

func TestMain(m *testing.M) {
	var err error
	app, err = newrelic.NewApplication(
		newrelic.ConfigLicense("0000000000000000000000000000000000000000"),
		newrelic.ConfigAppName("test"),
	)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestNewHandler(t *testing.T) {
	_, nrh, parent := setup(t)
	if nrh.GetOnlyForward() {
		t.Fatal("onlyForward must be true")
	}
	if nrh.GetAPP() != app {
		t.Fatal("app must be equal")
	}
	if got := nrh.GetHandler(); got != parent {
		t.Fatalf("handler must be equal to parent. Got:%v. Want:%v", got, parent)
	}
}

func TestEnabled(t *testing.T) {
	_, nrh, _ := setup(t)
	tests := map[string]struct {
		in   slog.Level
		want bool
	}{
		"debug": {in: slog.LevelDebug},
		"info":  {in: slog.LevelInfo, want: true},
		"warn":  {in: slog.LevelWarn, want: true},
		"error": {in: slog.LevelError, want: true},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			got := nrh.Enabled(context.Background(), testCase.in)
			if testCase.want != got {
				t.Fatal()
			}
		})
	}
}

func TestWithAttrs(t *testing.T) {
	buf, nhr, _ := setup(t)
	logger := slog.New(nhr.WithAttrs([]slog.Attr{slog.String("foo", "bar")}))

	logger.Info("test")

	var record struct {
		Foo string `json:"foo"`
	}
	err := json.NewDecoder(buf).Decode(&record)
	if err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if record.Foo != "bar" {
		t.Fatalf(`"foo" filed must be "bar", but got %s`, record.Foo)
	}
}

func TestWithGroup(t *testing.T) {
	buf, nhr, _ := setup(t)
	logger := slog.New(nhr.WithGroup("group"))
	logger.Info("test", slog.String("member", "foo"))

	var record struct {
		Group struct {
			Member string `json:"member"`
		} `json:"group"`
	}
	err := json.NewDecoder(buf).Decode(&record)
	if err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if record.Group.Member != "foo" {
		t.Fatalf(`group.member must be "foo", but got %s`, record.Group.Member)
	}
}

func setup(t *testing.T) (*bytes.Buffer, *nrslog.Handler, slog.Handler) {
	buf := bytes.NewBuffer(nil)
	parent := slog.NewJSONHandler(buf, nil)
	handler := nrslog.NewHandler(
		app,
		nrslog.WithHandler(parent),
	)
	nrh, ok := handler.(*nrslog.Handler)
	if !ok {
		t.Fatal("handler must be nrslog.handler")
	}
	return buf, nrh, parent
}
