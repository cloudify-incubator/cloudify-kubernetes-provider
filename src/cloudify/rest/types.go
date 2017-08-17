/*
Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rest

type CloudifyBaseMessage struct {
	Message         string `json:"message,omitempty"`
	ErrorCode       string `json:"error_code,omitempty"`
	ServerTraceback string `json:"server_traceback,omitempty"`
}

// Common
type CloudifyPagination struct {
	Total  uint `json:"total"`
	Offset uint `json:"offset"`
	Size   uint `json:"size"`
}

type CloudifyMetadata struct {
	Pagination CloudifyPagination `json:"pagination"`
}

type CloudifyResource struct {
	Id              string `json:"id"`
	Description     string `json:"description"`
	Tenant          string `json:"tenant_name"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	CreatedBy       string `json:"created_by"`
	PrivateResource bool   `json:"private_resource"`
}
