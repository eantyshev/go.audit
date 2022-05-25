package usecase

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.audit/internal/entity"
)

//nolint:gomnd
func TestFilterStats(t *testing.T) {
	cases := []struct {
		description string
		ids         []int64
		stats       []*entity.Stat
		expected    []*entity.Stat
	}{
		{
			"no stat for tailing banners",
			[]int64{0, 1, 2},
			[]*entity.Stat{{BannerID: 0, Clicks: 2, Shows: 3}},
			[]*entity.Stat{
				{BannerID: 0, Clicks: 2, Shows: 3},
				{BannerID: 1, Clicks: 0, Shows: 0},
				{BannerID: 2, Clicks: 0, Shows: 0}},
		},
		{
			"no stat for heading banners",
			[]int64{1, 2, 3},
			[]*entity.Stat{{BannerID: 3, Clicks: 2, Shows: 3}},
			[]*entity.Stat{
				{BannerID: 1, Clicks: 0, Shows: 0},
				{BannerID: 2, Clicks: 0, Shows: 0},
				{BannerID: 3, Clicks: 2, Shows: 3}},
		},
		{
			"unsorted, no stat for banners in the middle",
			[]int64{2, 4, 1},
			[]*entity.Stat{
				{BannerID: 4, Clicks: 0, Shows: 9},
				{BannerID: 2, Clicks: 1, Shows: 5},
				{BannerID: 3, Clicks: 2, Shows: 3}},
			[]*entity.Stat{
				{BannerID: 1, Clicks: 0, Shows: 0},
				{BannerID: 2, Clicks: 1, Shows: 5},
				{BannerID: 4, Clicks: 0, Shows: 9}},
		},
	}
	for i, c := range cases {
		log.Printf("case %d: %s\n", i, c.description)
		stats := FilterStats(c.ids, c.stats)
		assert.Equal(t, c.expected, stats)
	}
}
