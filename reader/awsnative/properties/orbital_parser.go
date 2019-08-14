package properties

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

var primitiveSpecTypes = map[string]bool{
	"string":          true,
	"bool":            true,
	"int64":           true,
	"float64":         true,
	"json.RawMessage": true,
	"interface{}":     true,
}

type awsTag struct {
	Type          string
	ItemType      string
	Primitive     bool
	PrimitiveItem bool
	Required      bool
}

func (h *awsTag) Parse(data string) {
	tag := strings.Split(data, ",")

	if len(tag) == 0 {
		return
	}

	if len(tag) > 0 {
		h.Type = tag[0]
		h.Primitive = primitiveSpecTypes[h.Type]
	}

	if len(tag) > 1 {
		if tag[1] == "Required" {
			h.Required = true
			return
		}

		if h.Type == "List" || h.Type == "Map" {
			h.ItemType = tag[1]
			h.PrimitiveItem = primitiveSpecTypes[h.ItemType]
		}
	}

	if len(tag) == 3 {
		h.Required = true
	}

	return
}

//Validate the field and remove the references
func Validate(metadata reflect.Type, b interface{}, name, resourcetype string) (out interface{}, err error) {
	out = b

	if b == nil {
		log.Printf("field not provided %s", name)
		return
	}

	tag, found := metadata.FieldByName(name)

	//This should not happen
	if !found {
		err = fmt.Errorf("Struct attribute %s not found", name)
		return
	}

	// Processing scalar
	if cftag, ok := tag.Tag.Lookup("aws"); ok {
		isRef := false
		tagInfo := &awsTag{}
		tagInfo.Parse(cftag)

		if tagInfo.Primitive {
			out, _ = refFlatter(b)
			return
		}

		/* List/Maps and complex types */
		if tagInfo.Type == "List" {
			innerList, ok := b.([]interface{})
			if !ok {
				err = fmt.Errorf("Entry error wrong type in template %s %+v expected a list", name, b)
				return
			}

			// checking for references
			for i, item := range innerList {
				containsRef := false
				innerList[i], containsRef = refFlatter(item)
				// It's not a reference and not a primitive type then it's a property
				if !containsRef && !tagInfo.PrimitiveItem {
					property, propError := getProperty(resourcetype+"."+tagInfo.ItemType, item)

					if propError != nil {
						err = propError
						continue
					}

					innerList[i] = property
				}
			}

			// set out
			out = innerList
			return
		}

		if tagInfo.Type == "Map" {
			innerMap, ok := b.(map[string]interface{})
			if !ok {
				err = fmt.Errorf("Entry error wrong type in template %s %+v expected a map", name, b)
				return
			}
			for key, value := range innerMap {
				containsRef := false
				innerMap[key], containsRef = refFlatter(value)

				if !containsRef && !tagInfo.PrimitiveItem {
					property, propError := getProperty(resourcetype+"."+tagInfo.ItemType, value)
					if propError != nil {
						err = propError
						continue
					}
					innerMap[key] = property
				}
			}

			out = innerMap
			return
		}

		// if we are reaching this is is an standalone property
		out, isRef = refFlatter(b)

		if !isRef {
			property, propError := getProperty(resourcetype+"."+tagInfo.Type, b)

			if propError != nil {
				err = propError
				return
			}
			out = property
		}
	}

	return
}

//ValidateItem from a type ellison
func ValidateItem(b interface{}, rawType, resourcetype string) (out interface{}, err error) {
	isRef := false
	out = b

	if b == nil {
		log.Printf("field not provided %s", resourcetype)
		return
	}

	out, isRef = refFlatter(b)
	// not reference and not primitive type
	if !isRef && !primitiveSpecTypes[rawType] {
		// looking for the type in properties
		property, propError := getProperty(resourcetype+"."+rawType, b)

		if propError != nil {
			err = propError
		} else {
			out = property
		}
	}

	return

}

func getProperty(propertyName string, raw interface{}) (CloudFormationType, error) {
	propertyGen, ok := PropertiesIndex[propertyName]
	if !ok {
		return nil, fmt.Errorf("Invalid property not found %s", propertyName)
	}

	property := propertyGen()

	bytes, err := json.Marshal(raw)

	if err != nil {
		return nil, err
	}
	errPropertyParse := property.UnmarshalJSON(bytes)
	if errPropertyParse != nil {
		return nil, errPropertyParse
	}
	return property, nil
}

func refFlatter(raw interface{}) (flatted interface{}, containsRef bool) {
	flatted = raw
	switch v := raw.(type) {
	case map[string]interface{}:
		if ref, ok := v["Ref"]; ok {
			flatted = fmt.Sprintf("HCLRef.%s", ref)
			containsRef = true
		}
	}
	return
}
