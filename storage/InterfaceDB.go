package storage

// Интерфейс и эквивалентные структуры

// Интерфейс реализует методы:

/*
Создавать новые задачи,
Получать список всех задач,
Получать список задач по автору,
Получать список задач по метке,
Обновлять задачу по id,
Удалять задачу по id.
*/

type InterfaceDB interface{
	CreateTask(t Task)(int, error)
	ListTasks()([]Task, error)
	ListTasks2()([]Task, error) // Метод для проверки альтернативного способа отлова Null
	ListTasksAuth(id int)([]Task, error)
	ListTasksLabel(l string)([]Task, error)
	UpdateTaskId(id int, task Task) (int, error)
	RemoveTaskId(id int) (int, error)
	TaskDoneId(id int)error
	Ping()error
}

type User struct{
	id int
	name string
}

type Label struct{
	id int
	name string
}

type TaskLabel struct{
	taskId int
	labelId int
}

type Task struct{
	Id         int
	Opened     int
	Closed     int
	AuthorId   int
	AssignedId int
	Title      string
	Content    string
}

