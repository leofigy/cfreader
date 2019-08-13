package properties

/* Global interfaces to parse this crap */

type CloudFormationType interface {
	UnmarshalJSON([]byte) error
	CloudformationType() string
}

type OrbitalType interface {
	IsProperty() bool
	OrbitalType() string
}

type FactoryFunctor func() CloudFormationType
