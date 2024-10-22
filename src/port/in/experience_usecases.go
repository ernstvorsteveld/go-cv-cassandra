package in

type ListExperienceCommand struct {
	Page int
	Size int
}

type CreateExperienceCommand struct {
	Name string
	Tags []string
}

type GetExperienceCommand struct {
	Id string
}

func NewListExperienceCommand(page int, size int) *ListExperienceCommand {
	return &ListExperienceCommand{Page: page, Size: size}
}

func NewCreateExperienceCommand(name string, tags []string) *CreateExperienceCommand {
	return &CreateExperienceCommand{Name: name, Tags: tags}
}

func NewGetExperienceCommand(id string) *GetExperienceCommand {
	return &GetExperienceCommand{Id: id}
}
