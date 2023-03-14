package operations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/attestantio/go-eth2-client/spec/capella"
)

type Operations struct {
	SigChanges []*capella.SignedBLSToExecutionChange
}

type opsFile struct {
	Name         string
	Download_url string
}

func (operations *Operations) Load(url string) error {

	apiClient := http.Client{Timeout: time.Second * 2}

	apiRequest, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return fmt.Errorf("Error initialising change operations list:", err)
	}

	apiResponse, err := apiClient.Do(apiRequest)

	if err != nil {
		return fmt.Errorf("Error fetching change operations list:", err)
	}

	apiResponseBody, err := ioutil.ReadAll(apiResponse.Body)

	var opsFilePaths []opsFile

	err = json.Unmarshal(apiResponseBody, &opsFilePaths)

	if err != nil {
		return fmt.Errorf("Error when Unmarshalling change operatios list:", err)
	}

	operations.SigChanges = make([]*capella.SignedBLSToExecutionChange, (len(opsFilePaths)))

	for i, file := range opsFilePaths {

		operations.SigChanges[i], err = fetchChangeOperation(&apiClient, file.Download_url)
	}

	return nil
}

func (operations *Operations) Print() {

	fmt.Println("Sig Change Operations:" , len(operations.SigChanges))
	fmt.Println("===========================================")
	for _, operation := range operations.SigChanges {
		decoded, _ := operation.MarshalJSON()
		fmt.Println("Operation:", string(decoded[:]))
	}
	fmt.Println("===========================================")

}

func fetchChangeOperation(apiClient *http.Client, url string) (*capella.SignedBLSToExecutionChange, error) {

	changeOp := capella.SignedBLSToExecutionChange{}

	apiRequest, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return &changeOp, fmt.Errorf("Error creating change operation request:", url, err)
	}

	apiResponse, err := apiClient.Do(apiRequest)

	if err != nil {
		return &changeOp, fmt.Errorf("Error fetching change operations file:", err)
	}

	apiResponseBody, err := ioutil.ReadAll(apiResponse.Body)

	var changeOpsList = make([]capella.SignedBLSToExecutionChange, 1)

	err = json.Unmarshal(apiResponseBody, &changeOpsList)

	if err != nil {
		return &changeOp, fmt.Errorf("Error when Unmarshalling change operatios item: ", err)
	}

	changeOp = changeOpsList[0]

	return &changeOp, nil
}
