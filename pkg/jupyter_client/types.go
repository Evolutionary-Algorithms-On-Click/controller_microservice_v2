// These types are generated based on the OpenAPI specification of Jupyter Kernel Gateway
// to ensure type safety and proper serialization/deserialization of JSON payloads.
package jupyterclient

import "time"

type ApiInfo struct {
	Version string `json:"version"`
}

type GetKernelSpecsResponse struct {
	Default     string                     `json:"default"`
	KernelSpecs map[string]KernelSpecEntry `json:"kernelspecs"`
}

type KernelSpecEntry struct {
	Name      string    `json:"name"`
	Spec      KernelSpecFile `json:"spec"`
	Resources map[string]string `json:"resources"`
}

type KernelSpecFile struct {
	Language        string                 `json:"language"`
	Argv            []string               `json:"argv"`
	DisplayName     string                 `json:"display_name"`
	CodemirrorMode  interface{}            `json:"codemirror_mode,omitempty"`
	Env             map[string]string      `json:"env,omitempty"`
	HelpLinks       []HelpLink             `json:"help_links,omitempty"`
}

type HelpLink struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}


type StartKernelRequest struct {
	Name string `json:"name"`
}

type Kernel struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	LastActivity   time.Time `json:"last_activity"`
	Connections    int       `json:"connections"`
	ExecutionState string    `json:"execution_state"`
}

// ErrorResponse represents a standard error response from the jupyter kernel gateway server.
type ErrorResponse struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
