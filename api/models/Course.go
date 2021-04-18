package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Course struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name       string    `gorm:"size:255;not null;unique" json:"name"`
	Start_date string    `gorm:"size:255;not null;unique" json:"start_date"`
	End_date   string    `gorm:"size:255;not null;unique" json:"end_date"`
	Cohort_id  string    `gorm:"size:100;not null;" json:"cohort_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Course) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Start_date = html.EscapeString(strings.TrimSpace(p.Start_date))
	p.End_date = html.EscapeString(strings.TrimSpace(p.End_date))
	p.Cohort_id = html.EscapeString(strings.TrimSpace(p.Cohort_id))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Course) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Name == "" {
		err = errors.New("Required Name")
		errorMessages["Required_name"] = err.Error()

	}
	if p.Start_date == "" {
		err = errors.New("Required Start_date")
		errorMessages["Required_startdate"] = err.Error()

	}
	if p.End_date == "" {
		err = errors.New("Required End_date")
		errorMessages["Required_enddate"] = err.Error()

	}
	if p.Cohort_id == "" {
		err = errors.New("Required cohort ID")
		errorMessages["Required_cohortid"] = err.Error()

	}
	return errorMessages
}

func (p *Course) SaveCourse(db *gorm.DB) (*Course, error) {
	var err error
	err = db.Debug().Model(&Course{}).Create(&p).Error
	if err != nil {
		return &Course{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &Course{}, err
	// 	}
	// }
	return p, nil
}

func (p *Course) FindAllCourses(db *gorm.DB) (*[]Course, error) {
	var err error
	courses := []Course{}
	err = db.Debug().Model(&Course{}).Limit(100).Order("created_at desc").Find(&courses).Error
	if err != nil {
		return &[]Course{}, err
	}
	// if len(courses) > 0 {
	// 	for i, _ := range courses {
	// 		err := db.Debug().Model(&User{}).Where("id = ?", courses[i].AuthorID).Take(&courses[i].Author).Error
	// 		if err != nil {
	// 			return &[]Course{}, err
	// 		}
	// 	}
	// }
	return &courses, nil
}

func (p *Course) FindCourseByID(db *gorm.DB, pid uint64) (*Course, error) {
	var err error
	err = db.Debug().Model(&Course{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Course{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &Course{}, err
	// 	}
	// }
	return p, nil
}

func (p *Course) UpdateACourse(db *gorm.DB) (*Course, error) {

	var err error

	err = db.Debug().Model(&Course{}).Where("id = ?", p.ID).Updates(Course{Name: p.Name, Start_date: p.Start_date, End_date: p.End_date, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Course{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &Course{}, err
	// 	}
	// }
	return p, nil
}

func (p *Course) DeleteACourse(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&Course{}).Where("id = ?", p.ID).Take(&Course{}).Delete(&Course{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// func (p *Course) FindUserCourses(db *gorm.DB, uid uint32) (*[]Course, error) {

// 	var err error
// 	courses := []Course{}
// 	err = db.Debug().Model(&Course{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&courses).Error
// 	if err != nil {
// 		return &[]Course{}, err
// 	}
// 	if len(courses) > 0 {
// 		for i, _ := range courses {
// 			err := db.Debug().Model(&User{}).Where("id = ?", courses[i].AuthorID).Take(&courses[i].Author).Error
// 			if err != nil {
// 				return &[]Course{}, err
// 			}
// 		}
// 	}
// 	return &courses, nil
// }

//When a user is deleted, we also delete the course that the user had
// func (c *Course) DeleteUserCourses(db *gorm.DB, uid uint32) (int64, error) {
// 	courses := []Course{}
// 	db = db.Debug().Model(&Course{}).Where("author_id = ?", uid).Find(&courses).Delete(&courses)
// 	if db.Error != nil {
// 		return 0, db.Error
// 	}
// 	return db.RowsAffected, nil
// }
