package main

import (
	"context"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createRecurringJob schedules a recurring cron job targeting the given resource.
func createRecurringJob(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, target aruba.Ref) *aruba.Job {
	printBanner("Recurring Job", "")

	j := aruba.NewJob().
		InProject(proj).
		Named(resourceName(NameJobRecurring)).
		Tagged("schedule").
		Tagged("recurring").
		InRegion(aruba.RegionITBGBergamo).
		Enabled().
		WithCron("0 10 * * *").
		RecurringUntil(time.Now().AddDate(0, 2, 0)).
		WithSteps(aruba.NewJobStep().
			Named("poweroff-step").
			Targeting(target).
			WithAction("poweroff").
			WithVerb(aruba.HTTPVerbPOST))

	res, err := arubaClient.FromSchedule().Jobs().Create(ctx, j)
	if err != nil {
		printCreateError("Recurring Job", err)
		return nil
	}
	printCreated("Recurring Job", res.Name(), res.JobID())
	return res
}

// createOneShotJob schedules a one-shot job to fire 24 hours from now on the given resource.
func createOneShotJob(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, target aruba.Ref) *aruba.Job {
	printBanner("One-Shot Job", "")

	j := aruba.NewJob().
		InProject(proj).
		Named(resourceName(NameJobOneShot)).
		Tagged("schedule").
		Tagged("oneshot").
		InRegion(aruba.RegionITBGBergamo).
		Enabled().
		OneShotAt(time.Now().Add(24 * time.Hour)).
		WithSteps(aruba.NewJobStep().
			Named("poweroff-step").
			Targeting(target).
			WithAction("poweroff").
			WithVerb(aruba.HTTPVerbPOST))

	res, err := arubaClient.FromSchedule().Jobs().Create(ctx, j)
	if err != nil {
		printCreateError("One-Shot Job", err)
		return nil
	}
	printCreated("One-Shot Job", res.Name(), res.JobID())
	return res
}

// deleteJob removes the scheduled job. Jobs persist as historical records on
// the platform after Delete (no "Deleted"/"Cancelled" state is enumerated in
// the SDK's jobTerminalStates), so polling for HTTP 404 always exhausts the
// wait budget without any signal. Submit the delete and move on.
func deleteJob(ctx context.Context, arubaClient aruba.Client, j *aruba.Job, label string) {
	pretty := label + " Job"
	printDeleteBanner(pretty)
	if err := arubaClient.FromSchedule().Jobs().Delete(ctx, j); err != nil {
		printDeleteError(pretty, err)
		return
	}
	printDeleteSubmitted(pretty, j.Name())
}
