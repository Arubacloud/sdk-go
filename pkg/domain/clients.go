package domain

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

/*
 * That is the proposal for the Domain
 *
 * The idea (largely used on industry) is to concentrate all the types which
 * users (developers in that case) will interact with.
 *
 * Both hierarchy and function names are structured to produce calls in a very
 * concise way such:
 *
 * - server, err := arubaClient.FromCompute().CloudServers().Create(...)
 * - keyPairs, err := arubaClient.FromCompute().KeyPairs().List(...)
 *
 * Developers can also easyly isolate only the part of the client which
 * interest, such below:
 *
 *   func DoSomethingWithAServer(id string, client arubasdk.CloudServersClient) error {
 *	     ... do stuff
 *   }
 *
 *   ... do stuff
 *
 *   serversClient := arubaClient.FromCompute().CloudServers()
 *
 *   ... do stuff
 *
 *   err := DoSomethingWithAServer("server-01", serversClient)
 *   if err != nil {
 *       return
 *   }
 *
 *   ... do stuff
 *
 */

type Client interface {
	FromAudit() AuditClient
	FromCompute() ComputeClient

	// TODO: transfer the other interfaces to here.
	//FromContainer() ContainerClient
	//FromDatabase() DatabaseClient
	//FromMetric() MetricClient
	//FromNetwork() NetworkClient
	//FromProject() ProjectClient
	//FromSchedule() ScheduleClient
	//FromSecurity() SecurityClient
	//FromStorage() StorageClient
}

type AuditClient interface {
	Events() EventsClient
}

type EventsClient interface {
	List(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.AuditEventListResponse], error)
}

type ComputeClient interface {
	CloudServers() CloudServersClient
	KeyPairs() KeyPairsClient
}

type CloudServersClient interface {
	List(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerList], error)
	Get(ctx context.Context, project string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
	Create(ctx context.Context, project string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
	Update(ctx context.Context, project string, cloudServerId string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
	Delete(ctx context.Context, projectId string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

type KeyPairsClient interface {
	List(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairListResponse], error)
	Get(ctx context.Context, project string, keyPairId string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error)
	Create(ctx context.Context, project string, body schema.KeyPairRequest, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error)
	Delete(ctx context.Context, projectId string, keyPairId string, params *schema.RequestParameters) (*schema.Response[any], error)
}
