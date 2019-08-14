package awsnative

import (
	"encoding/json"
	"errors"

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

//UnmarshalJSON basic parsing for object
func (r *Resource) UnmarshalJSON(b []byte) error{
	var holder rawResource
	err := json.Unmarshal(b, &holder)
	if err != nil{
		return err
	}

	// setup
	r.Type = holder.Type
	r.Metadata = holder.Metadata

	generator, ok := AwsnativeIndex[holder.Type]
	if !ok {
		return ErrUnsupportedResource
	}
	// Some resources have no properties 
	if holder.Properties == nil {
		return nil
	}

	nativeResource := generator()

	// inner unmarshal
	err = json.Unmarshal(*holder.Properties, nativeResource)
	if err != nil {
		return err
	}

	// final conversion
	r.Properties = nativeResource
	return nil
}
