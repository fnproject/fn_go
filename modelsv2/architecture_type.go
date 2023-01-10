// Code generated by go-swagger; DO NOT EDIT.

package modelsv2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// ArchitectureType Supported Architecture type enum
//
// swagger:model ArchitectureType
type ArchitectureType string

const (

	// ArchitectureTypeX86 captures enum value "x86"
	ArchitectureTypeX86 ArchitectureType = "x86"

	// ArchitectureTypeArm captures enum value "arm"
	ArchitectureTypeArm ArchitectureType = "arm"
)

// for schema
var architectureTypeEnum []interface{}

func init() {
	var res []ArchitectureType
	if err := json.Unmarshal([]byte(`["x86","arm"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		architectureTypeEnum = append(architectureTypeEnum, v)
	}
}

func (m ArchitectureType) validateArchitectureTypeEnum(path, location string, value ArchitectureType) error {
	if err := validate.EnumCase(path, location, value, architectureTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this architecture type
func (m ArchitectureType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateArchitectureTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
