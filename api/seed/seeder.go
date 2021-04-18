package seed

import (
	"log"

	"github.com/AlexSwiss/prentice/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Firstname: "Coded",
		Username:  "Fingers",
		Email:     "codedfingers@example.com",
		Gender:    "male",
		Phone:     "0987654321",
		Country:   "country",
		State:     "state",
		City:      "city",
		Area:      "area",
		Position:  "position",
		Password:  "password",
	},
}

var courses = []models.Course{
	models.Course{
		Name:       "Javascript",
		Start_date: "01 November 2021",
		End_date:   "14 June 2022",
		Cohort_id:  "54CU_2021",
	},
}

// var posts = []models.Post{
// 	models.Post{
// 		Title:   "Title 1",
// 		Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
// 	},
// 	models.Post{
// 		Title:   "Title 2",
// 		Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
// 	},
// }

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.User{}, &models.Course{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.Course{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
	for i, _ := range courses {
		err = db.Debug().Model(&models.Course{}).Create(&courses[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
