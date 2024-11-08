package middleware

type Error struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}

var ExperienceErrors = map[string]string{
	"EXP0000001": "error while retrieving experiences",
	"EXP0000002": "error while marshalling experiences",
	"EXP0000003": "error while marshalling input payload Experience",
	"EXP0000004": "internal server error: error while creating experience",
	"EXP0000005": "internal server error: error while marshalling experience",
	"EXP0000006": "experience with the provided id does not exist",
	"EXP0000007": "other error while retrieving experience",
	"TAG0000001": "internal server error while retrieving tags",
	"TAG0000002": "error while marshalling input payload tag",
	"TAG0000003": "internal Server error: error while creating tag",
	"TAG0000004": "not tags found",
	"TAG0000005": "error while marshalling input payload tag",

	"ALG0000001": "internal Server error",
}
