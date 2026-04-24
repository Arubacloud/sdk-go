package aruba

import (
	"context"
	"fmt"
)

// Wrapper is the constraint for List items: every resource wrapper satisfies it.
type Wrapper interface {
	URI() string
	ID() string
}

// List is a paginated collection of resource wrappers. Per-resource clients
// construct it after a server List call; callers use Next/Prev/All to iterate.
type List[T Wrapper] struct {
	items      []T
	total      int64
	self       string
	prev       string
	next       string
	first      string
	last       string
	callerOpts []CallOption
	refetch    func(ctx context.Context, url string) (*List[T], error)
	raw        any
}

// newList constructs a List from a server reply. lr holds the pagination links,
// raw is the full *types.Response[XxxList] stored as any for caller inspection,
// and refetch is provided by the per-resource client to fetch adjacent pages.
func newList[T Wrapper](
	items []T,
	total int64,
	self, prev, next, first, last string,
	raw any,
	opts []CallOption,
	refetch func(ctx context.Context, url string) (*List[T], error),
) *List[T] {
	return &List[T]{
		items:      items,
		total:      total,
		self:       self,
		prev:       prev,
		next:       next,
		first:      first,
		last:       last,
		callerOpts: opts,
		refetch:    refetch,
		raw:        raw,
	}
}

// Items returns the items on the current page.
func (l *List[T]) Items() []T { return l.items }

// Total returns the server-reported total item count across all pages.
func (l *List[T]) Total() int64 { return l.total }

// HasNext reports whether a next page is available.
func (l *List[T]) HasNext() bool { return l.next != "" }

// HasPrev reports whether a previous page is available.
func (l *List[T]) HasPrev() bool { return l.prev != "" }

// Next fetches the next page.
func (l *List[T]) Next(ctx context.Context) (*List[T], error) {
	if !l.HasNext() {
		return nil, fmt.Errorf("no next page")
	}
	return l.refetch(ctx, l.next)
}

// Prev fetches the previous page.
func (l *List[T]) Prev(ctx context.Context) (*List[T], error) {
	if !l.HasPrev() {
		return nil, fmt.Errorf("no previous page")
	}
	return l.refetch(ctx, l.prev)
}

// First fetches the first page.
func (l *List[T]) First(ctx context.Context) (*List[T], error) {
	if l.first == "" {
		return nil, fmt.Errorf("no first page link")
	}
	return l.refetch(ctx, l.first)
}

// Last fetches the last page.
func (l *List[T]) Last(ctx context.Context) (*List[T], error) {
	if l.last == "" {
		return nil, fmt.Errorf("no last page link")
	}
	return l.refetch(ctx, l.last)
}

// Cursor returns the next and previous page URLs.
func (l *List[T]) Cursor() (next, prev string) {
	return l.next, l.prev
}

// Raw returns the raw server response as any. Cast to the concrete
// *types.Response[XxxList] type to inspect response metadata.
func (l *List[T]) Raw() any { return l.raw }

// All iterates all pages, calling yield for each item. Iteration stops early
// if yield returns false. Returns the first error encountered while fetching.
func (l *List[T]) All(ctx context.Context, yield func(T) bool) error {
	current := l
	for {
		for _, item := range current.items {
			if !yield(item) {
				return nil
			}
		}
		if !current.HasNext() {
			return nil
		}
		next, err := current.Next(ctx)
		if err != nil {
			return err
		}
		current = next
	}
}
