package aruba

import "testing"

func TestDefaultBillingPeriod(t *testing.T) {
	tests := []struct {
		name  string
		input *BillingPeriod
		want  BillingPeriod
	}{
		{name: "nil defaults to Hour", input: nil, want: BillingPeriodHour},
		{name: "explicit value echoed back", input: func() *BillingPeriod { v := BillingPeriod("Month"); return &v }(), want: "Month"},
		{name: "Hour is echoed, not re-allocated", input: func() *BillingPeriod { v := BillingPeriodHour; return &v }(), want: BillingPeriodHour},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := defaultBillingPeriod(tc.input)
			if got == nil {
				t.Fatal("returned nil")
			}
			if *got != tc.want {
				t.Errorf("got %q, want %q", *got, tc.want)
			}
		})
	}
}
