package services

import (
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
)

type InServices struct {
	ep out.ExperienceDbPort
	tp out.TagDbPort
	ig utils.IdGenerator
}

func NewCvServices(ep out.ExperienceDbPort, tp out.TagDbPort, ig utils.IdGenerator) *InServices {
	return &InServices{
		ep: ep,
		tp: tp,
		ig: ig,
	}
}
