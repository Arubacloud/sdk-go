package aruba

// defaultBillingPeriod returns p when set, otherwise a pointer to the
// platform's default billing period (Hour). Centralising this here
// keeps every wrapper's toRequest() in sync with the API default.
func defaultBillingPeriod(p *BillingPeriod) *BillingPeriod {
	if p != nil {
		return p
	}
	v := BillingPeriodHour
	return &v
}
