// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package client

import (
	"fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"reflect"
	"testing"
)

type connInitializer interface {
	setConn(t *testing.T, result auto.UpResult) error
}

// CheckEnvStructValid validates an environment struct
func CheckEnvStructValid[Env any]() error {
	var env Env
	_, err := getFields(&env)
	return err
}

// CallConnInitializers validates an environment struct and initializes a connection to the testing infrastructure
func CallConnInitializers[Env any](t *testing.T, env *Env, connResult auto.UpResult) error {
	fields, err := getFields(env)

	for _, field := range fields {
		initializer := field.connInitializer
		if reflect.TypeOf(initializer).Kind() == reflect.Ptr && reflect.ValueOf(initializer).IsNil() {
			return fmt.Errorf("the field %v of %v is nil", field.name, reflect.TypeOf(env))
		}

		if err = initializer.setConn(t, connResult); err != nil {
			return err
		}
	}

	return err
}

type field struct {
	connInitializer connInitializer
	name            string
}

func getFields[Env any](env *Env) ([]field, error) {
	var fields []field
	envValue := reflect.ValueOf(*env)
	envType := reflect.TypeOf(*env)
	exportedFields := make(map[string]struct{})

	for _, f := range reflect.VisibleFields(envType) {
		if f.IsExported() {
			exportedFields[f.Name] = struct{}{}
		}
	}

	connInitializerType := reflect.TypeOf((*connInitializer)(nil)).Elem()
	for i := 0; i < envValue.NumField(); i++ {
		fieldName := envValue.Type().Field(i).Name
		if _, found := exportedFields[fieldName]; !found {
			return nil, fmt.Errorf("the field %v in %v is not exported", fieldName, envType)
		}

		initializer, ok := envValue.Field(i).Interface().(connInitializer)
		if !ok {
			return nil, fmt.Errorf("%v contains %v which doesn't implement %v",
				envType,
				fieldName,
				connInitializerType,
			)
		}
		fields = append(fields, field{
			connInitializer: initializer,
			name:            fieldName,
		})
	}
	return fields, nil
}
