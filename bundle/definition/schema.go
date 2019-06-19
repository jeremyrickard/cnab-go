package definition

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/qri-io/jsonschema"
)

type Definitions map[string]*Schema

// Schema represents a JSON Schema compatible CNAB Definition
type Schema struct {
	Comment              string                 `json:"$comment,omitempty" mapstructure:"$comment,omitempty"`
	ID                   string                 `json:"$id,omitempty" mapstructure:"$ref,omitempty"`
	Ref                  string                 `json:"$ref,omitempty" mapstructure:"$ref,omitempty"`
	AdditionalItems      interface{}            `json:"additionalItems,omitempty" mapstructure:"additionalProperties,omitempty"`
	AdditionalProperties interface{}            `json:"additionalProperties,omitempty" mapstructure:"additionalProperties,omitempty"`
	AllOf                []*Schema              `json:"allOf,omitempty" mapstructure:"allOf,omitempty"`
	Const                interface{}            `json:"const,omitempty" mapstructure:"const,omitempty"`
	Contains             *Schema                `json:"contains,omitempty" mapstructure:"contains,omitempty"`
	ContentEncoding      string                 `json:"contentEncoding,omitempty" mapstructure:"contentEncoding,omitempty"`
	ContentMediaType     string                 `json:"contentMediaType,omitempty" mapstructure:"contentMediaType,omitempty"`
	Default              interface{}            `json:"default,omitempty" mapstructure:"default,omitempty"`
	Definitions          Definitions            `json:"definitions,omitempty" mapstructure:"definitions,omitempty"`
	Dependencies         map[string]interface{} `json:"dependencies,omitempty" mapstructure:"dependencies,omitempty"`
	Description          string                 `json:"description,omitempty" mapstructure:"description,omitempty"`
	Else                 *Schema                `json:"else,omitempty" mapstructure:"else,omitempty"`
	Enum                 []interface{}          `json:"enum,omitempty" mapstructure:"enum,omitempty"`
	Examples             []interface{}          `json:"examples,omitempty" mapstructure:"examples,omitempty"`
	ExclusiveMaximum     *float64               `json:"exclusiveMaximum,omitempty" mapstructure:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum     *float64               `json:"exclusiveMinimum,omitempty" mapstructure:"exclusiveMinimum,omitempty"`
	Format               string                 `json:"format,omitempty" mapstructure:"format,omitempty"`
	If                   *Schema                `json:"if,omitempty" mapstructure:"if,omitempty"`
	//Items can be a Schema or an Array of Schema :(
	Items         interface{} `json:"items,omitempty" mapstructure:"items,omitempty"`
	Maximum       *float64    `json:"maximum,omitempty" mapstructure:"maximum,omitempty"`
	MaxLength     *float64    `json:"maxLength,omitempty" mapstructure:"maxLength,omitempty"`
	MinItems      *float64    `json:"minItems,omitempty" mapstructure:"minItems,omitempty"`
	MinLength     *float64    `json:"minLength,omitempty" mapstructure:"minLength,omitempty"`
	MinProperties *float64    `json:"minProperties,omitempty" mapstructure:"minProperties,omitempty"`
	Minimum       *float64    `json:"minimum,omitempty" mapstructure:"minimum,omitempty"`
	MultipleOf    *float64    `json:"multipleOf,omitempty" mapstructure:"multipleOf,omitempty"`
	Not           *Schema     `json:"not,omitempty" mapstructure:"not,omitempty"`
	OneOf         *Schema     `json:"oneOf,omitempty" mapstructure:"oneOf,omitempty"`

	PatternProperties map[string]*Schema `json:"patternProperties,omitempty" mapstructure:"patternProperties,omitempty"`

	Properties    map[string]*Schema `json:"properties,omitempty" mapstructure:"properties,omitempty"`
	PropertyNames *Schema            `json:"propertyNames,omitempty" mapstructure:"propertyNames,omitempty"`
	ReadOnly      *bool              `json:"readOnly,omitempty" mapstructure:"readOnly,omitempty"`
	Required      []string           `json:"required,omitempty" mapstructure:"required,omitempty"`
	Then          *Schema            `json:"then,omitempty" mapstructure:"then,omitempty"`
	Title         string             `json:"title,omitempty" mapstructure:"title,omitempty"`
	Type          interface{}        `json:"type,omitempty" mapstructure:"type,omitempty"`
	UniqueItems   *bool              `json:"uniqueItems,omitempty" mapstructure:"uniqueItems,omitempty"`
	WriteOnly     *bool              `json:"writeOnly,omitempty" mapstructure:"writeOnly,omitempty"`
}

// SingleType returns true if a schema has a single type declaration
func (s *Schema) SingleType() bool {
	_, ok := s.Type.(string)
	return ok
}

// GetType will return the singular type for a given schema and a success boolean. If the
// schema does not have a single type, it will return the false boolean and an error.
func (s *Schema) GetType() (string, bool, error) {
	if !s.SingleType() {
		return "", false, errors.New("this schema has multiple types")
	}
	typeString, ok := s.Type.(string)
	if !ok {
		return "", false, errors.New("couldn't get type")
	}
	return typeString, ok, nil
}

// GetType will return the types for a given schema and a success boolean. If the
// schema has a single type, it will return the false boolean and an error.
func (s *Schema) GetTypes() ([]string, bool, error) {
	if s.SingleType() {
		return nil, false, errors.New("this schema a single type")
	}
	data, ok := s.Type.([]interface{})
	if !ok {
		return nil, false, errors.New("unknown type values")
	}
	typeStrings := []string{}
	for _, val := range data {
		typeString, ok := val.(string)
		if !ok {
			return nil, false, errors.New("unknown type value")
		}
		typeStrings = append(typeStrings, typeString)
	}
	return typeStrings, ok, nil
}

// UnmarshalJSON provides an implementation of a JSON unmarshaler that uses the
// github.com/qri-io/jsonschema to load and validate a given schema. If it is valid,
// then the json is unmarshaled.
func (s *Schema) UnmarshalJSON(data []byte) error {

	// Before we unmarshal into the cnab-go bundle/definition/Schema type, unmarshal into
	// the library struct so we can handle any validation errors in the schema. If there
	// are any errors, return those.
	js := new(jsonschema.RootSchema)
	err := js.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	// The schema is valid at this point, so now use an indirect wrapper type to actually
	// unmarshal into our type.
	type wrapperType Schema
	wrapper := struct {
		*wrapperType
	}{
		wrapperType: (*wrapperType)(s),
	}
	err = json.Unmarshal(data, &wrapper)
	if err != nil {
		return err
	}
	return nil
}
