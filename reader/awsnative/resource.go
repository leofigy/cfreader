package awsnative

//Concrete types
type Resource struct {
	Type string `json:"Type"`
	Properties CloudFormationType `json:"Properties,omitempty"` 
	Metadata interface{} // not used now
}