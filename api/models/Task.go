package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Summary   string    `gorm:"size:2500;not null" json:"summary"`
	User      User      `json:user`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (t *Task) Prepare() {
	t.ID = 0
	t.Summary = html.EscapeString(strings.TrimSpace(t.Summary))
	t.User = User{}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *Task) Validate() error {
	if t.Summary == "" {
		return errors.New("Required Summary")
	}
	if t.UserID < 1 {
		return errors.New("Required User")
	}
	return nil
}

func (t *Task) SaveTask(db *gorm.DB) (*Task, error) {
	var err error
	err = db.Debug().Model(&Task{}).Create(&t).Error
	if err != nil {
		return &Task{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.UserID).Take(&t.User).Error
		if err != nil {
			return &Task{}, err
		}
	}
	return t, nil
}

func (t *Task) FindAllTasks(db *gorm.DB) (*[]Task, error) {
	var err error
	tasks := []Task{}
	err = db.Debug().Model(&Task{}).Limit(100).Find(&tasks).Error
	if err != nil {
		return &[]Task{}, err
	}
	if len(tasks) > 0 {
		for i := range tasks {
			err := db.Debug().Model(&User{}).Where("id = ?", tasks[i].UserID).Take(&tasks[i].User).Error
			if err != nil {
				return &[]Task{}, err
			}
		}
	}
	return &tasks, nil
}

func (t *Task) FindTaskByID(db *gorm.DB, tid uint64) (*Task, error) {
	var err error
	err = db.Debug().Model(&Task{}).Where("id = ?", tid).Take(&t).Error
	if err != nil {
		return &Task{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.UserID).Take(&t.User).Error
		if err != nil {
			return &Task{}, err
		}
	}
	return t, nil
}

func (t *Task) FindAllTasksByUserID(db *gorm.DB, uid uint32) (*[]Task, error) {
	var err error
	tasks := []Task{}
	err = db.Debug().Model(&Task{UserID: uid}).Limit(100).Find(&tasks).Error
	// err = db.Debug().Model(&Task{}).Limit(100).Find(&tasks, Task{UserID: uid}).Error // try this if above doesn't work
	if err != nil {
		return &[]Task{}, err
	}
	if len(tasks) > 0 {
		for i := range tasks {
			err := db.Debug().Model(&User{}).Where("id = ?", tasks[i].UserID).Take(&tasks[i].User).Error
			if err != nil {
				return &[]Task{}, err
			}
		}
	}
	return &tasks, nil
}

func (t *Task) UpdateATask(db *gorm.DB) (*Task, error) {
	var err error

	err = db.Debug().Model(&Task{}).Where("id = ?", t.ID).Updates(Task{Summary: t.Summary, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Task{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.UserID).Take(&t.User).Error
		if err != nil {
			return &Task{}, err
		}
	}
	return t, nil
}

func (t *Task) DeleteATask(db *gorm.DB, tid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Task{}).Where("id = ? and user_id = ?", tid, uid).Take(&Task{}).Delete(&Task{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Task not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
