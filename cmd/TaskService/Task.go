package main

import (
	"DbApp/storage"
	"DbApp/storage/postgres"
	"fmt"
	"log"
)

/*API пакета storage должен позволять:

/*
Создавать новые задачи,
Получать список всех задач,
Получать список задач по автору,
Получать список задач по метке,
Обновлять задачу по id,
Удалять задачу по id.
*/

var ConnString string = "postgres://postgres:user@localhost:5432/tasks"

// Для тестов:

var id int
var rows []storage.Task
var res int

var task1 = storage.Task{
	AuthorId:    0,
	AssignedId : 0,
	Title :      "Hello title",
	Content :    "There will be very interesting story",
}




func main(){

	var db storage.InterfaceDB
	var err error

	db, err = postgres.NewDb(ConnString)  // Присваиваем переменной интерфейса конкретную реализацию БД
	if err != nil{
		log.Fatal(err)
	}
	err = db.Ping()							// Проверка связи..
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Connection successful")

	id , err = db.CreateTask(task1)        // Загрузка задачи
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(id)


	rows, err = db.ListTasks()				// Вывод полного списка задач
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(len(rows))
	for _, val := range rows{
		fmt.Println(val)
	}


	fmt.Println( "Author with id = 0:")

	rows, err = db.ListTasksAuth(0) 		// Вывод списка задач по автору id = 0
	if err != nil{
		log.Fatal(err)
	}

	for _, r := range rows{
		fmt.Println(r)
	}

	fmt.Println( "Author with id = 1:")


	rows, err = db.ListTasksAuth(1)   // Вывод списка задач по автору id = 1
	if err != nil{
		log.Fatal(err)
	}

	for _, r := range rows{
		fmt.Println(r)
	}

	fmt.Println("List tasks with label work:")
	rows, err = db.ListTasksLabel("travel")
	if err != nil{
		log.Fatal(err)
	}
	for _, r := range rows{
		fmt.Println(r)
	}

	fmt.Println("List tasks with label work END")

	res, err = db.UpdateTaskId(100, 						// Обновление задачи
		storage.Task{Title: "Best story ", Content : " Here we go!"})
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Updated %v rows\n", res )

	res, err = db.RemoveTaskId(13)						// Удаление задачи
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v rows\n", res )
}