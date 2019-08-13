package awsnative

import (
    orbs "github.com/leofigy/cfreader/reader/awsnative/properties"
)

//Resource struct to read the cloudformation template
type Resource struct {
	Type       string             `json:"Type"`
	Properties orbs.CloudFormationType `json:"Properties,omitempty"`
	Metadata   interface{}        // not used now
}

