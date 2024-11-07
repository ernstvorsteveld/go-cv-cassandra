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
	"EXP0000004": "Internal Server error: error while creating experience",
	"EXP0000005": "Internal Server error: error while marshalling experience",
	"EXP0000006": "Experience with the provided id does not exist",
	"EXP0000007": "Other error wjile retrieving experience",
	"TAG0000001": "Error while retrieving tags",
	"TAG0000002": "error while marshalling input payload Tag",
	"TAG0000003": "Internal Server error: error while creating Tag",

	"ALG0000001": "Internal Server error",
}
