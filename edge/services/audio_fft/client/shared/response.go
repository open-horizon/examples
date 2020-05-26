// Copyright (c) 2020 SoftServe Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package shared

import (
	"encoding/json"
	"time"
)

type Response struct {
	Timestamp time.Time
	Trigger   bool
}

func NewResponse(ts time.Time) *Response {
	return &Response{
		Timestamp: ts,
		Trigger:   true,
	}
}

func FromJson(data []byte) (*Response, error){
	r := new(Response)
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *Response) ToJson() ([]byte, error) {
	return json.Marshal(r)
}
