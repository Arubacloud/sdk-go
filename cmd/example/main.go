package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/compute"
	"github.com/Arubacloud/sdk-go/pkg/spec/database"
	"github.com/Arubacloud/sdk-go/pkg/spec/network"
	"github.com/Arubacloud/sdk-go/pkg/spec/project"
	"github.com/Arubacloud/sdk-go/pkg/spec/schedule"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
	"github.com/Arubacloud/sdk-go/pkg/spec/security"
	"github.com/Arubacloud/sdk-go/pkg/spec/storage"
)

func main() {
	config := &client.Config{
		ClientID:     "cmp-74603fc1-ba10-40b3-9ff9-4aef35548642",
		ClientSecret: "UZXfEZFFOz1M0t66FNTcO4c5nez76Kwf",
		Debug:        true,
	}

	// Initialize the SDK (automatically obtains JWT token)
	sdk, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use the SDK with context
	sdk = sdk.WithContext(ctx)

	fmt.Println("\n=== SDK Examples ===")
	fmt.Println("Note: SDK automatically manages OAuth2 token acquisition and refresh")

	// Initialize service clients
	projectAPI := project.NewProjectService(sdk)
	computeAPI := compute.NewCloudServerService(sdk)
	networkAPI := network.NewVPCService(sdk)
	databaseAPI := database.NewDBaaSService(sdk)
	scheduleAPI := schedule.NewJobService(sdk)
	securityAPI := security.NewKmsKeyService(sdk)
	storageAPI := storage.NewBlockStorageService(sdk)

	// Example: Project Management
	fmt.Println("--- Project Management ---")

	// Create a new project
	newProject := schema.ProjectRequest{}
	createResp, err := projectAPI.CreateProject(ctx, newProject, nil)
	if err != nil {
		log.Printf("Error creating project: %v", err)
	} else if createResp.IsSuccess() {
		fmt.Printf("✓ Created project: %+v\n", createResp.Data)

		// Update the project
		updateResp, err := projectAPI.UpdateProject(ctx, "project-id", newProject, nil)
		if err != nil {
			log.Printf("Error updating project: %v", err)
		} else if updateResp.IsSuccess() {
			fmt.Printf("✓ Updated project: %+v\n", updateResp.Data)
		}
	}

	// Example: Compute - Cloud Server
	fmt.Println("\n--- Compute: Cloud Server ---")
	projectID := "your-project-id"

	// Create a cloud server
	serverReq := schema.CloudServerRequest{
		// Add your cloud server fields here
	}
	serverResp, err := computeAPI.CreateCloudServer(ctx, projectID, serverReq, nil)
	if err != nil {
		log.Printf("Error creating cloud server: %v", err)
	} else if serverResp.IsSuccess() {
		fmt.Printf("✓ Created cloud server: %+v\n", serverResp.Data)

		// Update the cloud server
		serverUpdateResp, err := computeAPI.UpdateCloudServer(ctx, projectID, "server-id", serverReq, nil)
		if err != nil {
			log.Printf("Error updating cloud server: %v", err)
		} else if serverUpdateResp.IsSuccess() {
			fmt.Printf("✓ Updated cloud server: %+v\n", serverUpdateResp.Data)
		}
	}

	// Example: Network - VPC
	fmt.Println("\n--- Network: VPC ---")

	// Create a VPC
	vpcReq := schema.VpcRequest{
		// Add your VPC fields here
	}
	vpcResp, err := networkAPI.CreateVPC(ctx, projectID, vpcReq, nil)
	if err != nil {
		log.Printf("Error creating VPC: %v", err)
	} else if vpcResp.IsSuccess() {
		fmt.Printf("✓ Created VPC: %+v\n", vpcResp.Data)

		// Update the VPC
		vpcUpdateResp, err := networkAPI.UpdateVPC(ctx, projectID, "vpc-id", vpcReq, nil)
		if err != nil {
			log.Printf("Error updating VPC: %v", err)
		} else if vpcUpdateResp.IsSuccess() {
			fmt.Printf("✓ Updated VPC: %+v\n", vpcUpdateResp.Data)
		}
	}

	// Example: Database - DBaaS
	fmt.Println("\n--- Database: DBaaS ---")

	// Create a DBaaS instance
	dbaasReq := schema.DBaaSRequest{
		// Add your DBaaS fields here
	}
	dbaasResp, err := databaseAPI.CreateDBaaS(ctx, projectID, dbaasReq, nil)
	if err != nil {
		log.Printf("Error creating DBaaS: %v", err)
	} else if dbaasResp.IsSuccess() {
		fmt.Printf("✓ Created DBaaS: %+v\n", dbaasResp.Data)

		// Update the DBaaS instance
		dbaasUpdateResp, err := databaseAPI.UpdateDBaaS(ctx, projectID, "dbaas-id", dbaasReq, nil)
		if err != nil {
			log.Printf("Error updating DBaaS: %v", err)
		} else if dbaasUpdateResp.IsSuccess() {
			fmt.Printf("✓ Updated DBaaS: %+v\n", dbaasUpdateResp.Data)
		}
	}

	// Example: Schedule - Job
	fmt.Println("\n--- Schedule: Job ---")

	// Create a scheduled job
	jobReq := schema.JobRequest{
		// Add your job fields here
	}
	jobResp, err := scheduleAPI.CreateScheduleJob(ctx, projectID, jobReq, nil)
	if err != nil {
		log.Printf("Error creating job: %v", err)
	} else if jobResp.IsSuccess() {
		fmt.Printf("✓ Created job: %+v\n", jobResp.Data)

		// Update the job
		jobUpdateResp, err := scheduleAPI.UpdateScheduleJob(ctx, projectID, "job-id", jobReq, nil)
		if err != nil {
			log.Printf("Error updating job: %v", err)
		} else if jobUpdateResp.IsSuccess() {
			fmt.Printf("✓ Updated job: %+v\n", jobUpdateResp.Data)
		}
	}

	// Example: Security - KMS Key
	fmt.Println("\n--- Security: KMS Key ---")

	// Create a KMS key
	kmsReq := schema.KmsRequest{
		// Add your KMS key fields here
	}
	kmsResp, err := securityAPI.CreateKMSKey(ctx, projectID, kmsReq, nil)
	if err != nil {
		log.Printf("Error creating KMS key: %v", err)
	} else if kmsResp.IsSuccess() {
		fmt.Printf("✓ Created KMS key: %+v\n", kmsResp.Data)

		// Update the KMS key
		kmsUpdateResp, err := securityAPI.UpdateKMSKey(ctx, projectID, "kms-key-id", kmsReq, nil)
		if err != nil {
			log.Printf("Error updating KMS key: %v", err)
		} else if kmsUpdateResp.IsSuccess() {
			fmt.Printf("✓ Updated KMS key: %+v\n", kmsUpdateResp.Data)
		}
	}

	// Example: Storage - Block Storage (Create only, no update)
	fmt.Println("\n--- Storage: Block Storage ---")

	// Create block storage
	storageReq := schema.BlockStorageRequest{
		// Add your storage fields here
	}
	storageResp, err := storageAPI.CreateBlockStorageVolume(ctx, projectID, storageReq, nil)
	if err != nil {
		log.Printf("Error creating block storage: %v", err)
	} else if storageResp.IsSuccess() {
		fmt.Printf("✓ Created block storage: %+v\n", storageResp.Data)
	}

	fmt.Println("\n=== SDK Examples Complete ===")
	fmt.Println("\nNote: All Create operations use POST to collection paths.")
	fmt.Println("All Update operations use PUT to item paths with resource IDs.")
}
