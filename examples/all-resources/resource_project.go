package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

// createProject creates and updates a project
func createProject(ctx context.Context, arubaClient aruba.Client) *aruba.Project {
	fmt.Println("--- Project Management ---")

	proj := aruba.NewProject().
		WithName(resourceName(NameProject)).
		AddTag("production").
		AddTag("arubacloud-sdk").
		WithDescription("My production project")

	created, err := arubaClient.FromProject().Create(ctx, proj)
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	}
	fmt.Printf("✓ Created project with ID: %s\n", created.ID())

	// Update the project
	updated, err := arubaClient.FromProject().Update(ctx, created)
	if err != nil {
		log.Printf("Error updating project: %v", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Updated project: %s\n", updated.Name())

	return updated
}

// updateProject updates a project
func updateProject(ctx context.Context, arubaClient aruba.Client, proj *aruba.Project) {
	fmt.Println("--- Updating Project ---")

	proj.
		WithName(updatedName(proj.Name())).
		ReplaceTags("production", "arubacloud-sdk", "updated").
		WithDescription("My production project - UPDATED")

	updated, err := arubaClient.FromProject().Update(ctx, proj)
	if err != nil {
		log.Printf("Error updating project: %v", err)
		return
	}
	fmt.Printf("✓ Updated project: %s\n", updated.Name())
}

// deleteProject deletes a project
func deleteProject(ctx context.Context, arubaClient aruba.Client, proj *aruba.Project) {
	fmt.Println("--- Deleting Project ---")

	if err := arubaClient.FromProject().Delete(ctx, proj); err != nil {
		log.Printf("Error deleting project: %v", err)
		return
	}
	fmt.Printf("✓ Deleted project: %s\n", proj.ID())
}
