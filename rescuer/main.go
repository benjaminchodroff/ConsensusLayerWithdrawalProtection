package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"strings"
	"time"

	config "clwp.rescuer/config"
	operations "clwp.rescuer/operations"
	eth2Client "github.com/attestantio/go-eth2-client"
	eth2Http "github.com/attestantio/go-eth2-client/http"
	zeroLog "github.com/rs/zerolog"
	"github.com/attestantio/go-eth2-client/spec/capella"
)

func main() {

	fmt.Println("CLWP Rescue Bot Starting")
	fmt.Println("===========================================")

	context, cancel := context.WithCancel(context.Background())

	var config = loadConfig() //loads default configuation.

	loadCLIArgs(&config) //override config with cli args.

	config.Print()


	var clients = createBeaconClients(context, config.BeaconEndpoints)

	printBeaconClients(context, clients)

	fmt.Println("Loading Change Ops Please wait.")
	fmt.Println("===========================================")

	var operations = loadOperations(config.ChangeOpsEndpoint)

	operations.Print()

	//Create batches with priority ops at start

	var operationBatches = createOperationBatches(
		append(config.PriorityOperations,operations.SigChanges...), 
		config.SubmitBatchSize,
	)

	for { //forever

		if time.Now().Unix() >= config.ForkTimestamp {

			fmt.Println("CLWP Rescue Bot Active")
			fmt.Println("===========================================")

			for _, client := range clients {

				var changeSubmitter = client.(eth2Client.BLSToExecutionChangesSubmitter)

				for i, batch := range operationBatches {

					fmt.Println("Submitting Change Ops To Client: ", client.Address(), "Batch Number:", i)
					fmt.Println("===========================================")

					var err = changeSubmitter.SubmitBLSToExecutionChanges(context, batch)

					if err != nil {
	
						fmt.Println("Error Submitting Change Ops", err)
	
					} else {
	
						fmt.Println("Successfully Submited Change Ops")
					}

				}


			}

			time.Sleep(time.Duration(config.SubmitPeriod) * time.Millisecond)

		} else {

			var sleepTime = time.Duration(
				math.Min(
					float64((config.ForkTimestamp-time.Now().Unix())*1000),
					float64(config.PollingPeriod),
				),
			) * time.Millisecond;


			fmt.Println(
				"CLWP Polling:",
				"Current Time Is:", time.Now().Unix(),
				"Activation Time Is:", config.ForkTimestamp,
				"Sleeping For:", sleepTime,
			)

			time.Sleep(sleepTime)

		}
	}

	cancel()
}

func loadConfig() config.Config {

	var config = config.Config{}

	err := config.Load("./config.json")

	if err != nil {
		fmt.Println("Error loading config:", err)
		panic("Error loading config")
	}

	return config
}

func loadCLIArgs(config *config.Config) {

	beaconEndpoints := strings.Join(config.BeaconEndpoints, ",")

	flag.StringVar(&beaconEndpoints, "beaconEndpoints", beaconEndpoints, "List of Beacon node endpoint Urls seperated by a ','")
	flag.StringVar(&config.ChangeOpsEndpoint, "changeOpsSource", config.ChangeOpsEndpoint, "Source for change operations files.")
	flag.Int64Var(&config.ForkTimestamp, "forkTimestamp", config.ForkTimestamp, "Timestamp for Capella Hardfork bot activation.")
	flag.Int64Var(&config.PollingPeriod, "pollingPeriod", config.PollingPeriod, "Polling period for bot to check for activation.")
	flag.Int64Var(&config.SubmitPeriod, "submitPeriod", config.SubmitPeriod, "Time between submissions once bot is active.")
	flag.Int64Var(&config.SubmitBatchSize, "submitBatchSize", config.SubmitBatchSize, "Size of batches of change ops to submit at a time.")

	flag.Usage()
	flag.Parse()

	config.BeaconEndpoints = strings.Split(beaconEndpoints, ",")

}

func loadOperations(endpoint string) operations.Operations {

	var operations = operations.Operations{}

	err := operations.Load(endpoint)

	if err != nil {
		fmt.Println("Error loading operations:", err)
		panic("Error loading operations")
	}

	return operations
}

func createOperationBatches(operations []*capella.SignedBLSToExecutionChange, batchSize int64) [][]*capella.SignedBLSToExecutionChange {

	batchCount := int64(math.Ceil( float64(len(operations)) / float64(batchSize)))

	batches := make([][]*capella.SignedBLSToExecutionChange, batchCount)

	fmt.Println("Creating", batchCount, "Batches")

	for i := range batches { 
		var batchStartIndex = batchSize * int64(i);
		var batchEndIndex = int64(math.Min(float64(len(operations)), float64((batchStartIndex + batchSize)) ));

		batches[i] = operations[batchStartIndex:batchEndIndex];

		fmt.Println("Loaded Batch From From", batchStartIndex, "To:", batchEndIndex)
	}

	fmt.Println("===========================================")

	return batches;
}

func createBeaconClients(context context.Context, endpoints []string) []eth2Client.Service {

	clients := make([]eth2Client.Service, len(endpoints))

	for i, endpoint := range endpoints {

		client, err := eth2Http.New(context,
			eth2Http.WithAddress(endpoint),
			eth2Http.WithLogLevel(zeroLog.InfoLevel),
			eth2Http.WithTimeout(time.Second*2),
		)

		if err != nil {
			fmt.Println("Error Creating Client:" , endpoint)
			panic("Error creating beacon node client")
		}

		clients[i] = client
	}

	return clients
}

func printBeaconClients(context context.Context, clients []eth2Client.Service) {

	for _, client := range clients {

		genesisProvider := client.(eth2Client.GenesisProvider)

		gensisInfo, err := genesisProvider.Genesis(context)

		fmt.Println("Genesis Information:")
		fmt.Println("===========================================")

		if err != nil {

			fmt.Println(
				"Client:", client.Address(),
				"Error fetching genesis information from beacon node",
			)

		} else {

			fmt.Println(
				"Client:", client.Address(),
				"Genesis Root Is:", gensisInfo.GenesisValidatorsRoot.String(),
				"Genesis Time Is:", gensisInfo.GenesisTime,
			)
		}

		fmt.Println("===========================================")
	}

}
