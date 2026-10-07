package models

type GetTimeRequest struct {
	Message string `json:"message"`
}


type TimeInput struct {
	Timezone string `json:"timezone" jsonschema:"IANA timezone, such as Asia/Kolkata"`
}