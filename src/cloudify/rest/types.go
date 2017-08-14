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
