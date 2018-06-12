package ec2

import (
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/client"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/client/metadata"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/request"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/signer/v4"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/private/protocol/ec2query"
)

// EC2 provides the API operation methods for making requests to
// Amazon Elastic Compute Cloud. See this package's package overview docs
// for details on the service.
//
// EC2 methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type EC2 struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Service information constants
const (
	ServiceName = "ec2"       // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the EC2 client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a EC2 client from just a session.
//     svc := ec2.New(mySession)
//
//     // Create a EC2 client with additional configuration
//     svc := ec2.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) (*EC2, error) {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) (*EC2, error) {
	svc := &EC2{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2016-11-15",
			},
			handlers,
		),
	}

	creds, err := cfg.Credentials.Get()
	if err != nil {
		return nil, err
	}

	// Handlers
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Sign.PushBackNamed(v4.NewSignRequestHandler(creds.SecretAccessKey))
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc, nil
}

// newRequest creates a new request for a EC2 operation and runs any
// custom request initialization.
func (c *EC2) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
