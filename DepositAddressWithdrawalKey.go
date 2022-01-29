package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	exec "github.com/attestantio/go-execution-client"
	"github.com/attestantio/go-execution-client/api"
	"github.com/attestantio/go-execution-client/jsonrpc"
	"github.com/attestantio/go-execution-client/spec"
	"github.com/rs/zerolog"
)

type blsPubKey [48]byte

var depositContractAddress spec.Address
var depositEventTopics []spec.Hash

func init() {
	tmp, err := hex.DecodeString("00000000219ab540356cBB839Cbe05303d7705Fa") // Deposit contract address.
	if err != nil {
		panic(err)
	}
	copy(depositContractAddress[:], tmp)

	tmp, err = hex.DecodeString("649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c5") // DepositEvent topic.
	if err != nil {
		panic(err)
	}
	var depositEventTopic spec.Hash
	copy(depositEventTopic[:], tmp)
	depositEventTopics = []spec.Hash{depositEventTopic}
}

func main() {
	ctx := context.Background()

	// JSON-RPC connection.
	client, err := jsonrpc.New(ctx,
		jsonrpc.WithAddress("http://localhost:8545/"),
		jsonrpc.WithLogLevel(zerolog.Disabled),
	)
	if err != nil {
		panic(err)
	}

	// Parameters.
	firstBlock := uint32(11052984) // Deployment block of deposit contract.
	lastBlock, err := client.(exec.ChainHeightProvider).ChainHeight(ctx)
	if err != nil {
		panic(err)
	}
	batchSize := uint32(1000)

	processed := make(map[blsPubKey]bool)
	if err := fetchEvents(ctx, client, firstBlock, lastBlock, batchSize, processed); err != nil {
		panic(err)
	}
}

func fetchEvents(ctx context.Context,
	client exec.Service,
	firstBlock uint32,
	lastBlock uint32,
	batchSize uint32,
	processed map[blsPubKey]bool,
) error {
	curBlock := firstBlock
	for {
		toBlock := curBlock + batchSize - 1
		if toBlock > lastBlock {
			toBlock = lastBlock
		}
		events, err := client.(exec.EventsProvider).Events(ctx, &api.EventsFilter{
			FromBlock: &curBlock,
			ToBlock:   &toBlock,
			Address:   &depositContractAddress,
			Topics:    &depositEventTopics,
		})
		if err != nil {
			panic(err)
		}
		for _, event := range events {
			if err := processEvent(ctx, client, event, processed); err != nil {
				return err
			}
		}
		curBlock += batchSize
	}
}

func processEvent(ctx context.Context,
	client exec.Service,
	event *spec.TransactionEvent,
	processed map[blsPubKey]bool,
) error {
	// Obtain withdrawal credentials from event data.
	var withdrawalCredentials spec.Hash
	copy(withdrawalCredentials[:], event.Data[288:320])
	if withdrawalCredentials[0] != 0 {
		// Only dealing with BLS withdrawal credentials.
	}

	// Obtain validator public key from event data.
	var pubKey blsPubKey
	copy(pubKey[:], event.Data[192:240])

	// Return if a deposit event for this public key has already processed.
	if _, exists := processed[pubKey]; exists {
		return nil
	}

	// Obtain the receipt for the event.
	receipt, err := client.(exec.TransactionReceiptsProvider).TransactionReceipt(ctx, event.TransactionHash)
	if err != nil {
		return err
	}

	// Ensure the transaction succeeded to register it.
	if receipt.Status != 1 {
		return nil
	}

	// Ignore anything not sent directly to the deposit contract (i.e. via contract).
	if !bytes.Equal(receipt.To[:], depositContractAddress[:]) {
		return nil
	}

	// Accept.
	fmt.Printf("%#x,%#x\n", pubKey, receipt.From)
	processed[pubKey] = true

	return nil
}
