package dto

type SettingResponse struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type UpdateSettingRequest struct {
	Value string `json:"value" binding:"required"`
}

type BulkUpdateSettingRequest struct {
	Settings map[string]string `json:"settings" binding:"required,min=1"`
}
