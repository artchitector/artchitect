package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArt_GetUnityMask(t *testing.T) {
	testCases := []struct {
		ID        uint
		UnityType int
		Result    string
	}{
		{
			ID:        1,
			UnityType: Unity100K,
			Result:    "0XXXXX",
		},
		{
			ID:        1,
			UnityType: Unity10K,
			Result:    "00XXXX",
		},
		{
			ID:        1,
			UnityType: Unity1K,
			Result:    "000XXX",
		},
		{
			ID:        1,
			UnityType: Unity100,
			Result:    "0000XX",
		},
		{
			ID:        654321,
			UnityType: Unity100K,
			Result:    "6XXXXX",
		},
		{
			ID:        654321,
			UnityType: Unity10K,
			Result:    "65XXXX",
		},
		{
			ID:        654321,
			UnityType: Unity1K,
			Result:    "654XXX",
		},
		{
			ID:        654321,
			UnityType: Unity100,
			Result:    "6543XX",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%d", tc.ID, tc.UnityType), func(t *testing.T) {
			art := Art{ID: tc.ID}
			result := art.GetUnityMask(tc.UnityType)
			assert.Equal(t, tc.Result, result)
		})
	}
}
