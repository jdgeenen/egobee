package egobee

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type testJSON struct {
	Duration TokenDuration `json:"duration"`
}

func TestUnmarshalTokenDuration(t *testing.T) {
	for _, tt := range []struct {
		name string
		json string
		want *testJSON
	}{
		{
			name: "unmarshal string duration with no units",
			json: `{"duration":"12345"}`,
			want: &testJSON{Duration: TokenDuration{Duration: time.Second * 12345}},
		},
		{
			name: "unmarshal float duration",
			json: `{"duration":12345}`,
			want: &testJSON{Duration: TokenDuration{Duration: time.Second * 12345}},
		},
		{
			name: "unmarshal string duration with units",
			json: `{"duration":"3h25m45s"}`,
			want: &testJSON{Duration: TokenDuration{Duration: time.Second * 12345}},
		},
	} {
		got := &testJSON{}
		if err := json.Unmarshal([]byte(tt.json), &got); err != nil {
			t.Errorf("%v: got unexpected error: %v", tt.name, err)
		} else if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%v: got: %v, wanted: %v", tt.name, got, tt.want)
		}
	}
}

func TestMarshalTokenDuration(t *testing.T) {
	for _, tt := range []struct {
		name string
		val  *testJSON
		want string
	}{
		{
			name: "marshal",
			val:  &testJSON{Duration: TokenDuration{Duration: time.Second * 12345}},
			want: `{"duration":"3h25m45s"}`,
		},
	} {
		if got, err := json.Marshal(tt.val); err != nil {
			t.Errorf("%v: got unexpected error: %v", tt.name, err)
		} else if string(got) != tt.want {
			t.Errorf("%v: got: %q, wanted: %q", tt.name, got, tt.want)
		}
	}
}
