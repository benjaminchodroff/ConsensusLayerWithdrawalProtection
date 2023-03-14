package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/attestantio/go-eth2-client/spec/capella"
)

type Config struct {
	ChangeOpsEndpoint  string
	BeaconEndpoints    []string
	ForkTimestamp      int64
	PollingPeriod      int64
	SubmitPeriod       int64
	SubmitBatchSize    int64
	PriorityOperations []*capella.SignedBLSToExecutionChange
}

func (config *Config) Load(file string) error {

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return fmt.Errorf("Error when opening config file: ", err)
	}

	err = json.Unmarshal(content, config)

	if err != nil {
		return fmt.Errorf("Error when Unmarshalling config: ", err)
	}

	return nil
}

func (config *Config) Print() {
	fmt.Println("Config:")
	fmt.Println("===========================================")
	fmt.Println("Beacon URLS    :", config.BeaconEndpoints)
	fmt.Println("Change Ops URL :", config.ChangeOpsEndpoint)
	fmt.Println("Fork Timestamp :", config.ForkTimestamp)
	fmt.Println("Polling Period :", config.PollingPeriod)
	fmt.Println("Submit Period  :", config.SubmitPeriod)
	fmt.Println("Submit Size    :", config.SubmitBatchSize)

	for _, operation := range config.PriorityOperations {
		decoded, _ := operation.MarshalJSON()
		fmt.Println("Priority Op    :", string(decoded[:]))
	}

	fmt.Println("===========================================")
}
