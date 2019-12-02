package main

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

type connector struct {
	grpcConnexion  *grpc.ClientConn
	dgraphClient   *dgo.Dgraph
	context        *context.Context
	awsSession     *session.Session
	awsRegion      string
	awsAccountName string
	awsAccountID   string
	ressources     map[string]map[string]string
	waitGroup      sync.WaitGroup
}

func newConnector(profile, region, dgraph *string) *connector {
	// Create AWS session (credentials from ~/.aws/config)
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,  //enable use of ~/.aws/config
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider, //ask for MFA if needed
		Profile:                 *profile,
		Config:                  aws.Config{Region: aws.String(*region)},
	}))

	resultSts, err := sts.New(awsSession).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err.Error())
	}

	connexion, err := grpc.Dial(*dgraph, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(connexion))
	ctx := context.Background()

	return &connector{
		dgraphClient:   dgraphClient,
		grpcConnexion:  connexion,
		context:        &ctx,
		awsSession:     awsSession,
		awsRegion:      *region,
		awsAccountName: *profile,
		awsAccountID:   *resultSts.Account,
		ressources:     map[string]map[string]string{},
		waitGroup:      sync.WaitGroup{},
	}
}
