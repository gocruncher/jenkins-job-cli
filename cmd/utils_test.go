package cmd

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFindBestChoice(t *testing.T) {
	table := []struct {
		val     string
		choice  string
		choices []string
	}{
		{"kuwa", "kuwait", []string{"albania", "samoa", "cameroun", "kuwait"}},
		{"samoa", "samoa", []string{"albania", "samoa", "cameroun", "kuwait"}},
		{"samoaabc", "samoa", []string{"albania", "samoa", "cameroun", "kuwait"}},
		{"", "albania", []string{"albania", "samoa", "cameroun", "kuwait"}},
		{"c", "cameroun", []string{"albania", "samoa", "cameroun", "kuwait"}},
		{"wwa", "albania", []string{"albania", "samoa", "cameroun", "kuwait"}},
		{"asdf", "", []string{}},
	}
	for _, tt := range table {
		t.Run(tt.val, func(t *testing.T) {
			rsp := findBestChoice(tt.val, tt.choices)
			assert.Equal(t, rsp, tt.choice, "wrong choice of "+tt.val+" case")
		})
	}

}
