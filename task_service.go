package gosis1

type TaskService struct {
	tasks  map[string]*Task
	parser URLParser
	nextId int
}

func NewTaskService(parser URLParser) *TaskService {
	return &TaskService{
		tasks:  make(map[string]*Task),
		parser: parser,
		nextId: 1,
	}
}
