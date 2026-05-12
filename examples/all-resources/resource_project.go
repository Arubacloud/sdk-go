package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createProject provisions a new project. Separating Create from Update makes
// the sequence explicit in the orchestrator.
func createProject(ctx context.Context, arubaClient aruba.Client) *aruba.Project {
	fmt.Println("--- Project ---")

	proj := aruba.NewProject().
		WithName(resourceName(NameProject)).
		AddTag("production").
		AddTag("arubacloud-sdk").
		WithDescription("My production project")

	created, err := arubaClient.FromProject().Create(ctx, proj)
	if err != nil {
		log.Fatalf("✗ Failed to create Project: %s", formatErr(err))
	}
	printCreated("Project", created.Name(), created.ID())

	return created
}

// updateProject updates a project name, tags, and description.
func updateProject(ctx context.Context, arubaClient aruba.Client, proj *aruba.Project) {
	fmt.Println("--- Updating Project ---")

	proj.
		WithName(updatedName(proj.Name())).
		ReplaceTags("production", "arubacloud-sdk", "updated").
		WithDescription("My production project - UPDATED")

	updated, err := arubaClient.FromProject().Update(ctx, proj)
	if err != nil {
		log.Printf("Error updating project: %s", formatErr(err))
		return
	}
	fmt.Printf("✓ Updated project: %s\n", updated.Name())
}

// deleteProject removes the project and waits until it is fully gone.
func deleteProject(ctx context.Context, arubaClient aruba.Client, proj *aruba.Project) {
	printDeleteBanner("Project")
	if err := arubaClient.FromProject().Delete(ctx, proj); err != nil {
		printDeleteError("Project", err)
		return
	}
	printDeleteSubmitted("Project", proj.Name())
	waitUntilGone(ctx, "Project "+proj.Name(), func(ctx context.Context) error {
		_, err := arubaClient.FromProject().Get(ctx, proj)
		return err
	})
}
