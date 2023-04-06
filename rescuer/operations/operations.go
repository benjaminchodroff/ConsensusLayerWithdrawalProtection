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

	err = json.Unmarshal(apiResponseBody, &operations.SigChanges)

	if err != nil {
		return fmt.Errorf("Error when Unmarshalling change operatios list:", err)
	}

	return nil
}

func (operations *Operations) Print() {

	fmt.Println("Sig Change Operations:", len(operations.SigChanges))
	fmt.Println("===========================================")
	for _, operation := range operations.SigChanges {
		decoded, _ := operation.MarshalJSON()
		fmt.Println("Operation:", string(decoded[:]))
	}
	fmt.Println("===========================================")

}
