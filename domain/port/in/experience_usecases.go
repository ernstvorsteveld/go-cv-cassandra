package in

type ListExperienceCommand struct {
	Limit int32
	Page  string
	Tag   string
	Name  string
}

type CreateExperienceCommand struct {
	Name string
	Tags []string
}

type GetExperienceCommand struct {
	Id string
}

type ListExperienceParameters struct {
	Limit int32
	Page  string
	Tag   string
	Name  string
}

func NewListExperienceCommand(parms *ListExperienceParameters) *ListExperienceCommand {
	return &ListExperienceCommand{Limit: parms.Limit, Page: parms.Page, Tag: parms.Tag, Name: parms.Name}
}

func NewCreateExperienceCommand(name string, tags []string) *CreateExperienceCommand {
	return &CreateExperienceCommand{Name: name, Tags: tags}
}

func NewGetExperienceCommand(id string) *GetExperienceCommand {
	return &GetExperienceCommand{Id: id}
}
