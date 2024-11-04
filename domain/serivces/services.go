package services

import (
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/in"
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
)

type InServices struct {
	ep out.ExperienceDbPort
	tp out.TagDbPort
}

func NewCvServices(ep out.ExperienceDbPort, tp out.TagDbPort) in.UseCasesPort {
	return &InServices{
		ep: ep,
		tp: tp,
	}
}