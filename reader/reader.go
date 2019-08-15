package reader

import (
	"fmt"
	"encoding/json"

	"github.com/leofigy/cfreader/reader/awsnative"
)

// see: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-anatomy.html
type Template struct {
	AWSTemplateFormatVersion string                 `json:"AWSTemplateFormatVersion,omitempty"`
	Transform                *awsnative.Transform   `json:"Transform,omitempty"`
	Description              string                 `json:"Description,omitempty"`
	Metadata                 map[string]interface{} `json:"Metadata,omitempty"`
	Parameters               map[string]interface{} `json:"Parameters,omitempty"`
	Mappings                 map[string]interface{} `json:"Mappings,omitempty"`
	Conditions               map[string]interface{} `json:"Conditions,omitempty"`
	Resources                Resources              `json:"Resources,omitempty"`
	Outputs                  map[string]interface{} `json:"Outputs,omitempty"`
}

type Resources map[string]awsnative.Resource

type SyntaxTransformer struct {
	Tmpl *Template
}

//UnmarshalJSON default
func (reader *SyntaxTransformer) UnmarshalJSON(src []byte) error {
	template := &Template{}
	if err := json.Unmarshal(src, template); err != nil {
		return err
	}
	reader.Tmpl = template
	return nil
}

func (resources *Resources) UnmarshalJSON(b []byte) error {
	var rawResources map[string]*json.RawMessage
	nativeResources := Resources{}

	err := json.Unmarshal(b, &rawResources)

	if err != nil {
		return err
	}

	for name, raw := range rawResources{
		nativeResource := awsnative.Resource{}
		err := json.Unmarshal(*raw, &nativeResource)
		if err == awsnative.ErrUnsupportedResource {
			fmt.Printf("Removing %s not supported", name)
			continue
		}

		// leaving as it is for debugging will be updated
		if err != nil {
			fmt.Printf("Resource %s translation issues %+v\n", name, err)
		}

		nativeResources[name] = nativeResource
	}

	*resources = nativeResources

	return nil
}
