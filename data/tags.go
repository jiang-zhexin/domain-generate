package data

import (
	"encoding/json"
	"fmt"

	mapset "github.com/deckarep/golang-set"
)

type Tags struct {
	Value mapset.Set
}

func NewTags(tags interface{}) *Tags {
	var t Tags

	switch ot := tags.(type) {
	case nil:
		t.Value = mapset.NewThreadUnsafeSet()
	case []string:
		s := make([]interface{}, len(ot))
		for i, v := range ot {
			s[i] = v
		}
		t.Value = mapset.NewThreadUnsafeSetFromSlice(s)
	case []interface{}:
		t.Value = mapset.NewThreadUnsafeSetFromSlice(ot)
	case string:
		t.Value = mapset.NewThreadUnsafeSet()
		t.Value.Add(ot)
	default:
		fmt.Println(tags)
		panic("unacceptable tag")
	}
	return &t
}

func (t *Tags) Union(otherTags *Tags) *Tags {
	if t == nil {
		var t Tags
		t.Value = otherTags.Value
		return &t
	}
	t.Value = t.Value.Union(otherTags.Value)
	return t
}

func (t *Tags) Add(tags interface{}) {
	ot := NewTags(tags)
	t.Value.Union(ot.Value)
}

func (t *Tags) UnmarshalJSON(data []byte) error {
	var arr interface{}
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	t.Value = NewTags(arr).Value
	return nil
}

func (t *Tags) MarshalJSON() ([]byte, error) {
	if t.Value.Cardinality() == 1 {
		return json.Marshal(t.Value.ToSlice()[0])
	}
	return json.Marshal(t.Value)
}

func (t *Tags) String() string {
	line, err := json.Marshal(t.Value)
	if err != nil {
		return ""
	}
	return string(line)
}
