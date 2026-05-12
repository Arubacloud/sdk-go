package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createRecurringJob schedules a recurring cron job targeting the given resource.
func createRecurringJob(ctx context.Context, arubaClient aruba.Client, proj aruba.Ref, target aruba.Ref) *aruba.Job {
	printBanner("Recurring Job", "")

	j := aruba.NewJob().
		IntoProject(proj).
		WithName(resourceName(NameJobRecurring)).
		AddTag("schedule").
		AddTag("recurring").
		InRegion(aruba.RegionITBGBergamo).
		WithEnabled(true).
		WithCron("0 10 * * *").
		RecurringUntil(time.Now().AddDate(0, 2, 0)).
		AddStep(aruba.NewJobStep().
			Named("poweroff-step").
			OfResource(target).
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
		IntoProject(proj).
		WithName(resourceName(NameJobOneShot)).
		AddTag("schedule").
		AddTag("oneshot").
		InRegion(aruba.RegionITBGBergamo).
		WithEnabled(true).
		OneShotAt(time.Now().Add(24 * time.Hour)).
		AddStep(aruba.NewJobStep().
			Named("poweroff-step").
			OfResource(target).
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

// deleteJob removes the scheduled job identified by label.
func deleteJob(ctx context.Context, arubaClient aruba.Client, j *aruba.Job, label string) {
	fmt.Printf("--- Deleting %s Job ---\n", label)
	if err := arubaClient.FromSchedule().Jobs().Delete(ctx, j); err != nil {
		log.Printf("Error deleting %s Job: %s", label, formatErr(err))
		return
	}
	fmt.Printf("✓ Deleted %s Job: %s\n", label, j.Name())
}
