package internal

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetRange(t *testing.T) {

	for _, v := range []struct {
		d    division
		want Ranges
	}{
		{
			d:    division{0, 0, 0, 0},
			want: Ranges{},
		},
		{
			d:    division{1024, 0, 0, 0},
			want: Ranges{},
		},
		{
			d:    division{0, 1024, 0, 0},
			want: Ranges{},
		},
		{
			d: division{10, 1024, 0, 0},
			want: Ranges{
				Range{0, 0, 9},
			},
		},
		{
			d: division{1024, 1024, 0, 0},
			want: Ranges{
				Range{0, 0, 1023},
			},
		},
		{
			d: division{2048, 1024, 0, 0},
			want: Ranges{
				Range{0, 0, 1023}, Range{1, 1024, 2047},
			},
		},
		{
			d: division{2049, 1024, 0, 0},
			want: Ranges{
				Range{0, 0, 1023}, Range{1, 1024, 2047}, Range{2, 2048, 2048},
			},
		},

		{
			d:    division{0, 0, 1024, 1},
			want: Ranges{},
		},
		{
			d:    division{1024, 0, 1024, 1},
			want: Ranges{},
		},
		{
			d:    division{0, 1024, 1024, 1},
			want: Ranges{},
		},
		{
			d: division{10, 1024, 8, 1},
			want: Ranges{
				Range{1, 8, 9},
			},
		},
		{
			d:    division{1024, 1024, 1024, 1},
			want: Ranges{},
		},
		{
			d: division{2048, 1024, 1024, 2},
			want: Ranges{
				Range{2, 1024, 2047},
			},
		},
		{
			d: division{2049, 1024, 1024, 2},
			want: Ranges{
				Range{2, 1024, 2047}, Range{3, 2048, 2048},
			},
		},
	} {
		t.Run(fmt.Sprintf(`%+v`, v.d), testGetRange(t, v.d, v.want))

	}

}

func testGetRange(t *testing.T, d division, want Ranges) func(t *testing.T) {

	f := func(t *testing.T) {

		if got := getRanges(d); !cmp.Equal(want, got) {
			t.Errorf(`want:%+v, got%+v`, want, got)
		}
	}
	return f
}

func TestGetHeader(t *testing.T) {

	for _, v := range []struct {
		name string
		url  string
		want header
	}{
		{
			name: "yes",
			url:  "http://i.imgur.com/z4d4kWk.jpg",
			want: header{true, 146515},
		},
		{
			name: "none",
			url:  "https://www.youtube.com/watch?v=EwTZ2xpQwpA",
			want: header{false, 0},
		},
		{
			name: "empty",
			url:  "https://www.google.com",
			want: header{false, 0},
		},
	} {
		t.Run(v.name, testGetHeader(t, v.url, v.want))
	}

}

func testGetHeader(t *testing.T, url string, want header) func(t *testing.T) {

	f := func(t *testing.T) {
		h, err := getHeader(url)
		if err != nil {
			t.Error(err)
			return
		}

		if got := h; !cmp.Equal(want, got) {
			t.Errorf(`want:%+v, got%+v`, want, got)
		}
	}
	return f
}
