package dto

type UpdateSettingRequest struct {
	Value string `json:"value" binding:"required" example:"12"`
}

type BulkUpdateSettingRequest struct {
	Settings map[string]string `json:"settings" binding:"required,min=1"`
}

type SettingResponse struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
