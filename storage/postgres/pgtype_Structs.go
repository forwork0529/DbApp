package postgres

// Паект для преобразования нулевых значений с использованием встроенных типов pgtype

import (
	"DbApp/storage"
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
)

type taskPgType struct{
	id pgtype.Int4
	opened pgtype.Int8
	closed pgtype.Int8
	authorId pgtype.Int4
	assignedId pgtype.Int4
	title pgtype.Varchar
	content pgtype.Varchar
}

func (t *taskPgType)Convert()(storage.Task, error){
	var err error
	realTask  := storage.Task{}

	defer func(err error) {
		if r := interface{}(recover()); r != nil{   // interface{} - моя версия GoLand ругалась на алиас и я добавил это
			err = errors.New(fmt.Sprintf("%v",r))
		}
	}(err)

	realTask.Id = int(t.id.Int)
	realTask.Opened = int(t.opened.Int)
	if t.closed.Status == pgtype.Present{
		realTask.Closed = int(t.closed.Int)
	}
	realTask.AuthorId = int(t.authorId.Int)
	realTask.AssignedId = int(t.assignedId.Int)
	if t.title.Status == pgtype.Present{
		realTask.Title = t.title.String
	}
	if t.content.Status == pgtype.Present{
		realTask.Content = t.content.String
	}

	return realTask, err
}
