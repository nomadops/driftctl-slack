package driftctl

import (
	"encoding/json"
	"log"
)

// Scan struct is the json output of `driftctl scan`
type Scan struct {
	Alerts struct {
		_ []struct {
			Message string `json:"message"`
		} `json:""`
	} `json:"alerts"`
	Coverage    int64       `json:"coverage"`
	Differences interface{} `json:"differences"`
	Managed     []struct {
		ID     string `json:"id"`
		Source struct {
			InternalName string `json:"internal_name"`
			Namespace    string `json:"namespace"`
			Source       string `json:"source"`
		} `json:"source"`
		Type string `json:"type"`
	} `json:"managed"`
	Missing []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"missing"`
	ProviderName    string `json:"provider_name"`
	ProviderVersion string `json:"provider_version"`
	Summary         struct {
		TotalChanged   int64 `json:"total_changed"`
		TotalManaged   int64 `json:"total_managed"`
		TotalMissing   int64 `json:"total_missing"`
		TotalResources int64 `json:"total_resources"`
		TotalUnmanaged int64 `json:"total_unmanaged"`
	} `json:"summary"`
	Unmanaged []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"unmanaged"`
}

// ScanOutput returns a map[string]interface{} from `driftctl scan -o json://file`
func ScanOutput(driftctlJSON []byte) (map[string]interface{}, error) {
	var scan Scan
	err := json.Unmarshal(driftctlJSON, &scan)
	if err != nil {
		log.Fatal("Error when unmarshalling: ", err)
		return nil, err
	}

	var output map[string]interface{}
	scanOut, _ := json.Marshal(scan)
	json.Unmarshal(scanOut, &output)

	return output, nil
}

// ScanSummary returns a map[string]string of the Summary section from `driftctl scan -o json://file`
func ScanSummary(driftctlJSON []byte) (map[string]int, error) {
	// This all works
	scanOutput, err := ScanOutput(driftctlJSON)
	if err != nil {
		log.Fatal("Error when unmarshalling: ", err)
		return nil, err
	}

	foo, _ := json.Marshal(scanOutput["summary"])
	var summary map[string]int
	json.Unmarshal(foo, &summary)

	return summary, nil
}
