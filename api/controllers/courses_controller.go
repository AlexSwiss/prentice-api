package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/AlexSwiss/prentice/api/models"
	"github.com/AlexSwiss/prentice/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

func (server *Server) CreateCourse(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	course := models.Course{}

	err = json.Unmarshal(body, &course)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// uid, err := auth.ExtractTokenID(c.Request)
	// if err != nil {
	// 	errList["Unauthorized"] = "Unauthorized"
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status": http.StatusUnauthorized,
	// 		"error":  errList,
	// 	})
	// 	return
	// }

	// check if the user exist:
	// user := models.User{}
	// err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	// if err != nil {
	// 	errList["Unauthorized"] = "Unauthorized"
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status": http.StatusUnauthorized,
	// 		"error":  errList,
	// 	})
	// 	return
	// }

	//course.AuthorID = uid //the authenticated user is the one creating the course

	course.Prepare()
	errorMessages := course.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	courseCreated, err := course.SaveCourse(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": courseCreated,
	})
}

func (server *Server) GetCourses(c *gin.Context) {

	course := models.Course{}

	courses, err := course.FindAllCourses(server.DB)
	if err != nil {
		errList["No_course"] = "No Course Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": courses,
	})
}

func (server *Server) GetCourse(c *gin.Context) {

	courseID := c.Param("id")
	pid, err := strconv.ParseUint(courseID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	course := models.Course{}

	courseReceived, err := course.FindCourseByID(server.DB, pid)
	if err != nil {
		errList["No_course"] = "No Course Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": courseReceived,
	})
}

func (server *Server) UpdateCourse(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	courseID := c.Param("id")
	// Check if the course id is valid
	pid, err := strconv.ParseUint(courseID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	// //CHeck if the auth token is valid and  get the user id from it
	// uid, err := auth.ExtractTokenID(c.Request)
	// if err != nil {
	// 	errList["Unauthorized"] = "Unauthorized"
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status": http.StatusUnauthorized,
	// 		"error":  errList,
	// 	})
	// 	return
	// }
	//Check if the course exist
	origCourse := models.Course{}
	err = server.DB.Debug().Model(models.Course{}).Where("id = ?", pid).Take(&origCourse).Error
	if err != nil {
		errList["No_course"] = "No Course Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// if uid != origCourse.AuthorID {
	// 	errList["Unauthorized"] = "Unauthorized"
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status": http.StatusUnauthorized,
	// 		"error":  errList,
	// 	})
	// 	return
	// }
	// Read the data courseed
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Start processing the request data
	course := models.Course{}
	err = json.Unmarshal(body, &course)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	course.ID = origCourse.ID //this is important to tell the model the course id to update, the other update field are set above
	//course.AuthorID = origCourse.AuthorID

	course.Prepare()
	errorMessages := course.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	courseUpdated, err := course.UpdateACourse(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": courseUpdated,
	})
}

func (server *Server) DeleteCourse(c *gin.Context) {

	courseID := c.Param("id")
	// Is a valid course id given to us?
	pid, err := strconv.ParseUint(courseID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	fmt.Println("this is delete course sir")

	// // Is this user authenticated?
	// uid, err := auth.ExtractTokenID(c.Request)
	// if err != nil {
	// 	errList["Unauthorized"] = "Unauthorized"
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status": http.StatusUnauthorized,
	// 		"error":  errList,
	// 	})
	// 	return
	// }
	// Check if the course exist
	course := models.Course{}
	err = server.DB.Debug().Model(models.Course{}).Where("id = ?", pid).Take(&course).Error
	if err != nil {
		errList["No_course"] = "No Course Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// // Is the authenticated user, the owner of this course?
	// if uid != course.AuthorID {
	// 	errList["Unauthorized"] = "Unauthorized"
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status": http.StatusUnauthorized,
	// 		"error":  errList,
	// 	})
	// 	return
	// }
	// If all the conditions are met, delete the course
	_, err = course.DeleteACourse(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	// comment := models.Comment{}
	// like := models.Like{}

	// // Also delete the likes and the comments that this course have:
	// _, err = comment.DeleteCourseComments(server.DB, pid)
	// if err != nil {
	// 	errList["Other_error"] = "Please try again later"
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status": http.StatusInternalServerError,
	// 		"error":  errList,
	// 	})
	// 	return
	// }
	// _, err = like.DeleteCourseLikes(server.DB, pid)
	// if err != nil {
	// 	errList["Other_error"] = "Please try again later"
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status": http.StatusInternalServerError,
	// 		"error":  errList,
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Course deleted",
	})
}

// func (server *Server) GetUserCourses(c *gin.Context) {

// 	userID := c.Param("id")
// 	// Is a valid user id given to us?
// 	uid, err := strconv.ParseUint(userID, 10, 64)
// 	if err != nil {
// 		errList["Invalid_request"] = "Invalid Request"
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status": http.StatusBadRequest,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	course := models.Course{}
// 	courses, err := course.FindUserCourses(server.DB, uint32(uid))
// 	if err != nil {
// 		errList["No_course"] = "No Course Found"
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status": http.StatusNotFound,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":   http.StatusOK,
// 		"response": courses,
// 	})
// }
