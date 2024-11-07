package in

type ListTagsCommand struct {
	Page int
	Size int
}

func NewListTagsCommand(page int, size int) *ListTagsCommand {
	return &ListTagsCommand{Page: page, Size: size}
}

type GetTagByIdCommand struct {
	Id string
}

func NewGetTagByIdCommand(id string) *GetTagByIdCommand {
	return &GetTagByIdCommand{Id: id}
}

type CreateTagCommand struct {
	Name string
}

func NewCreateTagCommand(name string) *CreateTagCommand {
	return &CreateTagCommand{Name: name}
}
