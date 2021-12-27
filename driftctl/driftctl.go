package driftctl

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

// Scan struct is the json output of `driftctl scan`.
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

// scanOutput returns a map[string]interface{} from `driftctl scan -o json://file`.
func scanOutput(report []byte) (map[string]interface{}, error) {
	var scan Scan
	err := json.Unmarshal(report, &scan)
	if err != nil {
		log.Fatal().Msg("Error when unmarshalling &scan ")
		return nil, err
	}

	var output map[string]interface{}
	scanOut, _ := json.Marshal(scan)
	json.Unmarshal(scanOut, &output)

	return output, nil
}

// ScanSummary returns a map[string]string of the Summary section from `driftctl scan -o json://file`.
func ScanSummary(report []byte) (map[string]int, error) {
	// This all works
	scan, err := scanOutput(report)
	if err != nil {
		log.Fatal().Msg("Error when unmarshalling scanOutput(report)")
		return nil, err
	}

	foo, _ := json.Marshal(scan["summary"])
	var summary map[string]int
	json.Unmarshal(foo, &summary)

	return summary, nil
}

// Run executs the driftctl scan command and returns the output as a byte slice.
func Run(bucket string, driftctlJSON string) (err error) {
	tfstates := fmt.Sprintf("tfstate+s3://%v/**/*.tfstate", bucket)
	target := fmt.Sprintf("json://%v", driftctlJSON)
	cmd := exec.Command("driftctl", "scan", "--quiet", "--from", tfstates, "-o", target)

	stdout, err := cmd.Output()
	log.Info().Msg(string(stdout))
	// This might be the wrong approach. Driftctl will return an error if there are changes, which will be every change without a really well-defined driftignore.
	if err != nil {
		log.Info().Msg("Driftctl scan detected drift.")
	}

	log.Info().Msg(string(stdout))
	return nil
}
