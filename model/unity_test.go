package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUnityParent(t *testing.T) {
	testCases := []struct {
		childMask  string
		parentMask string
	}{
		{
			childMask:  "0000XX",
			parentMask: "000XXX",
		},
		{
			childMask:  "0100XX",
			parentMask: "010XXX",
		},
		{
			childMask:  "010XXX",
			parentMask: "01XXXX",
		},
		{
			childMask:  "01XXXX",
			parentMask: "0XXXXX",
		},
		{
			childMask:  "91XXXX",
			parentMask: "9XXXXX",
		},
		{
			childMask:  "9XXXXX",
			parentMask: "",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s -> %s", tc.childMask, tc.parentMask), func(t *testing.T) {
			result := getParentMask(tc.childMask)
			assert.Equal(t, tc.parentMask, result)
		})
	}
}
