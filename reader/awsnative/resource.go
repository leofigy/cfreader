package awsnative

import (
	"errors"
	"encoding/json"

    orbs "github.com/leofigy/cfreader/reader/awsnative/properties"
)

//ErrUnsupportedResource fatal error
var ErrUnsupportedResource = errors.New("Resource not supported")

//Resource struct to read the cloudformation template
type Resource struct {
	Type       string             `json:"Type"`
	Properties orbs.CloudFormationType `json:"Properties,omitempty"`
	Metadata   interface{}        // not used now
}

type rawResource struct {
	Type string
	Properties *json.RawMessage
	Metadata interface{}
}

func (r *Resource) UnmarshalJSON(b []byte) error{
	var holder rawResource
	err := json.Unmarshal(b, &holder)
	if err != nil{
		return err
	}

	nativeResource, ok := AwsnativeIndex[holder.Type]

	if !ok {
		return ErrUnsupportedResource
	}

	// Some resources have no properties 
	if holder.Properties == nil {
		r.Type = holder.Type
		r.Metadata = holder.Metadata
		return
	}

	


	return nil
}
