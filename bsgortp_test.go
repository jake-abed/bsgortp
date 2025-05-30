package bsgortp

import (
	"testing"
)

func TestEmptyTextPost(t *testing.T) {
	_, err := GenPost("", []string{"en"})
	if err == nil {
		t.Errorf("empty string should cause an error")
	}
}

func TestSimplePosts(t *testing.T) {
	tests := []string{
		"hi!",
		"I'm in that mode?",
	}

	for _, tt := range tests {
		post, err := GenPost(tt, []string{"en"})

		if err != nil {
			t.Errorf("input=%s failed: %s", tt, err)
		}

		if len(post.Facets) != 0 {
			t.Errorf("expected 0 facets from input=%s", tt)
		}
	}
}

func TestPostsWithLinks(t *testing.T) {
	tests := []struct {
		Input              string
		ExpectedUrl        string
		ExpectedFacetCount int
		ExpectedByteStart  int64
		ExpectedByteEnd    int64
	}{
		{
			"go visit https://cats.cool",
			"https://cats.cool",
			1,
			9,
			26,
		},
		{
			"https://lucky.me is a copy of dog.dev",
			"https://lucky.me",
			2,
			0,
			16,
		},
		{
			"my website is jakeabed.dev",
			"https://jakeabed.dev",
			1,
			14,
			26,
		},
		{
			"http://scooby.doo redirects to jakeabed.dev",
			"http://scooby.doo",
			2,
			0,
			17,
		},
	}

	for _, tt := range tests {
		post, err := GenPost(tt.Input, []string{"en"})

		if err != nil {
			t.Errorf("input=%s failed: %s", tt.Input, err)
		}

		if len(post.Facets) != tt.ExpectedFacetCount {
			t.Errorf("got %d facets, expected %d",
				len(post.Facets), tt.ExpectedFacetCount)
		}

		facet := post.Facets[0]
		feature := facet.Features[0]

		if feature.RichtextFacet_Link.Uri != tt.ExpectedUrl {
			t.Errorf(
				"expected url=%s, got=%s",
				tt.ExpectedUrl,
				feature.RichtextFacet_Link.Uri,
			)
		}

		idx := facet.Index

		if idx.ByteStart != tt.ExpectedByteStart {
			t.Errorf("incorrect byte start: got=%d - expected %d",
				idx.ByteStart, tt.ExpectedByteStart)
		}

		if idx.ByteEnd != tt.ExpectedByteEnd {
			t.Errorf("incorrect byte end: got=%d - expected %d",
				idx.ByteEnd, tt.ExpectedByteEnd)
		}
	}
}

func TestAllInOne(t *testing.T) {
	post, err := GenPost(
		"Hey @jakeabed.dev, jakeabed.dev is a #buggy site",
		[]string{"en"},
	)
	if err != nil {
		t.Errorf("error generating post: %s", err.Error())
	}

	if len(post.Facets) != 3 {
		t.Errorf("expected 3 facets, got=%d", len(post.Facets))
	}

	third := post.Facets[0].Features[0]
	second := post.Facets[1].Features[0]
	first := post.Facets[2].Features[0]

	if first.RichtextFacet_Mention == nil {
		t.Errorf("expected first facet to be mention!")
	}

	if second.RichtextFacet_Link == nil {
		t.Errorf("expected second facet to be link!")
	}

	if third.RichtextFacet_Tag == nil {
		t.Errorf("expected third facet to be tag!")
	}

	if second.RichtextFacet_Link.Uri != "https://jakeabed.dev" {
		t.Errorf("second facet uri wrong, got=%s expected=%s",
			second.RichtextFacet_Link.Uri, "https://jakeabed.dev")
	}

	if third.RichtextFacet_Tag.Tag != "buggy" {
		t.Errorf("third facet tag wrong, got=%s expected=%s",
			third.RichtextFacet_Tag.Tag, "buggy")
	}

}
