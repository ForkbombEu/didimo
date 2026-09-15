// SPDX-FileCopyrightText: 2026 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package validators

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// JOSEJWEProtectedHeaderValidator verifies a field in a compact JWE protected
// header. Nested fields use dot notation, for example epk.crv.
type JOSEJWEProtectedHeaderValidator struct{}

func (JOSEJWEProtectedHeaderValidator) ID() string { return "jose.jwe_protected_header" }

func (JOSEJWEProtectedHeaderValidator) Validate(_ context.Context, input Input) Result {
	params, err := DecodeParams[struct {
		Field   string `json:"field"`
		Value   any    `json:"value"`
		Present *bool  `json:"present"`
	}](input.Params)
	if err != nil {
		return Result{Status: StatusError, Message: err.Error()}
	}
	if params.Field == "" {
		return Result{Status: StatusError, Message: fieldParamRequired}
	}
	compact, ok := input.Value.(string)
	if !ok {
		return Result{Status: StatusFail, Message: fmt.Sprintf("input is %T, expected compact JWE", input.Value)}
	}
	parts := strings.Split(compact, ".")
	if len(parts) != 5 {
		return Result{Status: StatusFail, Message: "input is not a compact JWE"}
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Result{Status: StatusFail, Message: "JWE protected header is not valid base64url"}
	}
	header := map[string]any{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Result{Status: StatusFail, Message: "JWE protected header is not valid JSON"}
	}
	value, found := nestedJWEHeaderValue(header, params.Field)
	if params.Present != nil {
		if found != *params.Present {
			return Result{Status: StatusFail, Message: fmt.Sprintf("JWE protected header field %q presence is %t, expected %t", params.Field, found, *params.Present)}
		}
		return Result{Status: StatusPass, Message: fmt.Sprintf("JWE protected header field %q presence matches", params.Field)}
	}
	if !found {
		return Result{Status: StatusFail, Message: fmt.Sprintf("JWE protected header field %q is missing", params.Field)}
	}
	if !reflect.DeepEqual(value, params.Value) {
		return Result{Status: StatusFail, Message: fmt.Sprintf("JWE protected header field %q is %v, expected %v", params.Field, value, params.Value)}
	}
	return Result{Status: StatusPass, Message: fmt.Sprintf("JWE protected header field %q matches", params.Field)}
}

func nestedJWEHeaderValue(header map[string]any, field string) (any, bool) {
	var current any = header
	for _, component := range strings.Split(field, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = object[component]
		if !ok {
			return nil, false
		}
	}
	return current, true
}
