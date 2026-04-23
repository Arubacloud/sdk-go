package aruba

import (
	"context"
	"testing"
)

// testItem is a minimal Wrapper implementation for list tests.
type testItem struct {
	id  string
	uri string
}

func (i testItem) ID() string  { return i.id }
func (i testItem) URI() string { return i.uri }

func makeTestList(items []testItem, prev, next, first, last string, refetch func(context.Context, string) (*List[testItem], error)) *List[testItem] {
	return newList[testItem](items, int64(len(items)), "", prev, next, first, last, nil, nil, refetch)
}

func TestList_Items(t *testing.T) {
	items := []testItem{{id: "a", uri: "/a"}, {id: "b", uri: "/b"}}
	l := makeTestList(items, "", "", "", "", nil)
	if got := l.Items(); len(got) != 2 {
		t.Errorf("Items() len = %d, want 2", len(got))
	}
	if l.Total() != 2 {
		t.Errorf("Total() = %d, want 2", l.Total())
	}
}

func TestList_HasNextHasPrev(t *testing.T) {
	l := makeTestList(nil, "", "/page2", "", "", nil)
	if !l.HasNext() {
		t.Error("HasNext() should be true")
	}
	if l.HasPrev() {
		t.Error("HasPrev() should be false when prev empty")
	}
}

func TestList_Next(t *testing.T) {
	var capturedURL string
	page2 := makeTestList([]testItem{{id: "c"}}, "", "", "", "", nil)

	refetch := func(_ context.Context, url string) (*List[testItem], error) {
		capturedURL = url
		return page2, nil
	}

	l := makeTestList([]testItem{{id: "a"}}, "", "/page2", "", "", refetch)
	got, err := l.Next(context.Background())
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}
	if capturedURL != "/page2" {
		t.Errorf("refetch called with %q, want %q", capturedURL, "/page2")
	}
	if got != page2 {
		t.Error("Next() did not return expected page")
	}
}

func TestList_NextNoLink(t *testing.T) {
	l := makeTestList(nil, "", "", "", "", nil)
	if _, err := l.Next(context.Background()); err == nil {
		t.Error("expected error when no next link")
	}
}

func TestList_Prev(t *testing.T) {
	var capturedURL string
	page0 := makeTestList(nil, "", "", "", "", nil)
	refetch := func(_ context.Context, url string) (*List[testItem], error) {
		capturedURL = url
		return page0, nil
	}
	l := makeTestList(nil, "/page0", "", "", "", refetch)
	if _, err := l.Prev(context.Background()); err != nil {
		t.Fatalf("Prev() error: %v", err)
	}
	if capturedURL != "/page0" {
		t.Errorf("refetch URL = %q", capturedURL)
	}
}

func TestList_FirstLast(t *testing.T) {
	var calls []string
	refetch := func(_ context.Context, url string) (*List[testItem], error) {
		calls = append(calls, url)
		return makeTestList(nil, "", "", "", "", nil), nil
	}
	l := makeTestList(nil, "", "", "/first", "/last", refetch)

	if _, err := l.First(context.Background()); err != nil {
		t.Fatalf("First() error: %v", err)
	}
	if _, err := l.Last(context.Background()); err != nil {
		t.Fatalf("Last() error: %v", err)
	}
	if len(calls) != 2 || calls[0] != "/first" || calls[1] != "/last" {
		t.Errorf("calls = %v", calls)
	}
}

func TestList_Cursor(t *testing.T) {
	l := makeTestList(nil, "/prev", "/next", "", "", nil)
	next, prev := l.Cursor()
	if next != "/next" || prev != "/prev" {
		t.Errorf("Cursor() = (%q, %q)", next, prev)
	}
}

func TestList_All_TwoPages(t *testing.T) {
	page2 := newList[testItem](
		[]testItem{{id: "c"}, {id: "d"}},
		4, "", "", "", "", "", nil, nil,
		nil,
	)
	refetch := func(_ context.Context, _ string) (*List[testItem], error) {
		return page2, nil
	}
	page1 := newList[testItem](
		[]testItem{{id: "a"}, {id: "b"}},
		4, "", "", "/page2", "", "", nil, nil,
		refetch,
	)

	var collected []string
	err := page1.All(context.Background(), func(item testItem) bool {
		collected = append(collected, item.id)
		return true
	})
	if err != nil {
		t.Fatalf("All() error: %v", err)
	}
	if len(collected) != 4 {
		t.Errorf("collected = %v, want [a b c d]", collected)
	}
}

func TestList_All_EarlyStop(t *testing.T) {
	l := makeTestList([]testItem{{id: "a"}, {id: "b"}, {id: "c"}}, "", "", "", "", nil)
	var count int
	_ = l.All(context.Background(), func(_ testItem) bool {
		count++
		return count < 2 // stop after second item
	})
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestList_Raw(t *testing.T) {
	raw := struct{ n int }{42}
	l := newList[testItem](nil, 0, "", "", "", "", "", raw, nil, nil)
	got, ok := l.Raw().(struct{ n int })
	if !ok || got.n != 42 {
		t.Errorf("Raw() = %v", l.Raw())
	}
}
