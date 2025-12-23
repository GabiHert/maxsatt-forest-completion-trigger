package dto

type ForestData struct {
	ForestId string `json:"forest_id" validate:"required,uuid"`
}
