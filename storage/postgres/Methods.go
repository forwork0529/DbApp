package postgres

import (
	"DbApp/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

// Выбор метода отлова Null значений в БД (Идеи):
// 1. Исполльзовать типы pgx , конверитировать поструктурно ( реализовал второй версией метода)
// 2. Использовать типы pgx , конверитровать попеременно - как то это неудобно из за
											// метода Rows.Scan()...
// 3. Использовать ссылки на ссылки         // Как то мне эта идея не понравилась, реализую
											// конгда-нибудь..
// 4. Использовать COALESCE - как по мне самый удачный из опробованных методов!
											// по поводу COALESCE прошу отписаться проверяющего..
// 5. Не разрешать пустные поля в БД -  решение отметаю , потому как: что хочу то и храню в своей базе
											// null значит null, обработку выношу на уровень Go!

func(s *Storage) CreateTask(t storage.Task)(int, error){
	var id int
	var err error
	err = s.db.QueryRow(context.Background(),`
INSERT INTO tasks (author_id, assigned_id, title, content)
VALUES ($1, $2, $3, $4) RETURNING id;`, t.AuthorId, t.AuthorId, t.Title, t.Content).Scan(&id)
	if err != nil{

		return 0, err
	}
	return id, nil
}

func(s *Storage) ListTasks()([]storage.Task, error){
	rows, err := s.db.Query(context.Background(),
		`SELECT id, COALESCE(opened, 0) AS opened, COALESCE(closed, 0) AS closed ,
			COALESCE(author_id, 0) AS author_id, COALESCE(assigned_id, 0) AS assigned_id ,
			COALESCE(title, '') AS title, COALESCE(content, '') AS content FROM tasks;`)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	var listTask []storage.Task
	for rows.Next(){
		t := storage.Task{}
		err = rows.Scan(&t.Id, &t.Opened, &t.Closed, &t.AuthorId,
			&t.AssignedId, &t.Title, &t.Content )
		if err != nil{
			return nil, err
		}

		listTask = append(listTask, t)
	}

	err = rows.Err()
	if err != nil{
		return nil, err
	}
	return listTask, nil
}

func (s *Storage) ListTasks2()([]storage.Task, error){ // Метод возврата с ипользованием
	// втроенных типов для отлова нулевых значений кортежей..

	rows, err := s.db.Query(context.Background(),`SELECT id,
	opened, closed, author_id, assigned_id, title, content FROM tasks;`)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	var liTask []storage.Task


	for rows.Next(){
		var medLiTask taskPgType
		err = rows.Scan(&medLiTask.id,
					&medLiTask.opened,
					&medLiTask.closed,
					&medLiTask.authorId,
					&medLiTask.assignedId,
					&medLiTask.title,
					&medLiTask.content)
		if err != nil{
			return nil, err
		}
		var newTask storage.Task
		newTask , err = medLiTask.Convert()
		if err != nil{
			return nil, err
		}

		liTask = append(liTask, newTask)

	}
	err = rows.Err()
	if err != nil{
		return nil, err
	}
	return liTask, nil
}

func(s *Storage) ListTasksAuth(id int)([]storage.Task, error){
	var listTask []storage.Task
/*	`SELECT id, COALESCE(opened, 0) AS opened, COALESCE(closed, 0) AS closed ,
			COALESCE(author_id, 0) AS author_id, COALESCE(assigned_id, 0) AS assigned_id ,
			COALESCE(title, '') AS title, COALESCE(content, '') AS content FROM tasks;`)*/

	rows , err := s.db.Query(context.Background(),`SELECT id, COALESCE(opened, 0)
	AS opened, COALESCE (closed, 0) AS closed, COALESCE (author_id, 0) AS author_id,
	COALESCE (assigned_id, 0) AS assigned_id, COALESCE (title, '') AS title,
	COALESCE(content, '') AS content FROM tasks WHERE author_id = $1;`,id)

	if err == pgx.ErrNoRows{
		return listTask, nil
	}
	if err != nil{
		return nil , err
	}
	defer rows.Close()

	for rows.Next(){
		var t storage.Task
		err = rows.Scan(&t.Id, &t.Opened, &t.Closed, &t.AuthorId,
			&t.AssignedId, &t.Title, &t.Content )
		if err != nil{
			return nil, err
		}
		listTask = append(listTask, t)
	}
	if err = rows.Err(); err != nil{
		return nil, err
	}
	fmt.Println(len(listTask))
	return listTask, nil
}

func(s *Storage) ListTasksLabel(label string)([]storage.Task, error){
	var listTask []storage.Task
	rows, err := s.db.Query(context.Background(), `
	SELECT tasks.id, COALESCE(tasks.opened, 0) AS opened,
	COALESCE (tasks.closed, 0) AS closed, COALESCE (tasks.author_id,0) AS author_id,
	COALESCE (tasks.assigned_id, 0) AS assigned_id, COALESCE (tasks.title, '') AS title,
	COALESCE (tasks.content, '') from tasks_labels
	join tasks ON tasks.id = tasks_labels.task_id
	join labels ON tasks_labels.label_id = labels.id and labels.name = $1;`, label)
	if err == pgx.ErrNoRows{
		return listTask, nil
	}
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var t storage.Task
		err = rows.Scan(&t.Id, &t.Opened, &t.Closed, &t.AuthorId,
			&t.AssignedId, &t.Title, &t.Content )
		if err != nil{
			return nil, err
		}
		listTask = append(listTask, t)
	}
	err = rows.Err()
	if err != nil{
		return nil, err
	}
	return listTask, nil
}

func(s *Storage) UpdateTaskId(id int, task storage.Task) (int, error){

	 res, err := s.db.Exec(context.Background(), `UPDATE tasks SET closed = $1 ,
		author_id = $2, assigned_id = $3,
		title = $4, content = $5 WHERE id =  $6;`,
		task.Closed, task.AuthorId, task.AssignedId,
	 	task.Title, task.Content, id)
	 if err != nil{
	 	return 0, err
	 }

	return int(res.RowsAffected()), nil
}

func(s *Storage) RemoveTaskId(id int)(int, error){
	res, err := s.db.Exec(context.Background(), `DELETE FROM tasks WHERE id = $1`, id)

	return int(res.RowsAffected()), err
}

func(s *Storage) TaskDoneId(id int)error{
	return nil
}



func(s *Storage) Ping()error{
	err := s.db.Ping(context.Background())
	if err != nil{
		return err
	}
	return nil
}


func DoNothing(){  // Проверка созданой структуры БД на соответствие контракту интерфейса
	var a storage.InterfaceDB
	var er error
	a, er = NewDb("hello world")
	fmt.Println(a, er)
}
