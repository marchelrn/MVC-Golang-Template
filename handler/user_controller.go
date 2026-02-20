package handler

import (
	"fmt"
	"mini_jira/contract"
	"mini_jira/dto"
	"mini_jira/middleware"
	"net/http"

	errs "mini_jira/pkg/error"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service contract.UserService
}

func (c *UserController) InitService(s *contract.Service) {
	fmt.Println("DEBUG: Initializing UserController with UserService")
	if s == nil {
		fmt.Println("ERROR: Provided service is nil")
		return
	}

	if s.User == nil {
		fmt.Println("ERROR: UserService in provided service is nil")
		return
	}
	c.service = s.User
	fmt.Println("DEBUG: UserController initialized successfully with UserService")
}

func (u *UserController) GetAll(c *gin.Context) {
	response, err := u.service.GetAll()
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Users retrieved successfully",
		"Users":   response,
	})
}

func (u *UserController) GetById(c *gin.Context) {
	userIdParam := c.Param("Id")
	var userId uint
	_, err := fmt.Sscanf(userIdParam, "%d", &userId)
	if err != nil {
		HandleError(c, errs.BadRequest("Invalid user ID format"))
		return
	}

	response, err := u.service.GetById(userId)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "User retrieved successfully",
		"User":    response,
	})
}

func (u *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	response, err := u.service.GetUserByUsername(username)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "User retrieved successfully",
		"User":    response,
	})
}

func (u *UserController) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	response, err := u.service.GetUserByEmail(email)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "User retrieved successfully",
		"User":    response,
	})
}

func (u *UserController) Login(c *gin.Context) {
	var payload dto.UserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleError(c, errs.BadRequest("Invalid request body"))
		return
	}

	response, err := u.service.Login(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(response.StatusCode, response)
}

func (u *UserController) Register(c *gin.Context) {
	var payload dto.UserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleError(c, errs.BadRequest("Invalid request body"))
		return
	}

	response, err := u.service.Register(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(response.StatusCode, response)
}

func (u *UserController) RefreshToken(c *gin.Context) {
	var payload dto.TokenRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleError(c, errs.BadRequest("Invalid request body"))
		return
	}

	response, err := u.service.RefreshToken(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(response.StatusCode, response)
}

func (u *UserController) Logout(c *gin.Context) {
	var payload dto.TokenRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleError(c, errs.BadRequest("Invalid request body"))
		return
	}

	response, err := u.service.Logout(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(response.StatusCode, response)
}

func (u *UserController) Update(c *gin.Context) {
	userIdParam := c.Param("Id")
	var userId uint
	_, err := fmt.Sscanf(userIdParam, "%d", &userId)
	if err != nil {
		HandleError(c, errs.BadRequest("Invalid user ID format"))
		return
	}

	tokenUserID, _ := c.Get(middleware.ContextUserID)
	tokenRole, _ := c.Get(middleware.ContextRole)
	if tokenUserID.(uint) != userId && tokenRole.(string) != string("admin") {
		HandleError(c, errs.Forbidden("You can only update your own account"))
		return
	}

	var payload dto.UserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleError(c, errs.BadRequest("Invalid request body"))
		return
	}
	payload.Id = userId
	payload.RequestingRole = c.GetString(middleware.ContextRole)

	response, err := u.service.Update(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(response.StatusCode, response)
}

func (u *UserController) Delete(c *gin.Context) {
	userIdParam := c.Param("Id")
	var userId uint
	_, err := fmt.Sscanf(userIdParam, "%d", &userId)
	if err != nil {
		HandleError(c, errs.BadRequest("Invalid user ID format"))
		return
	}

	// IDOR protection: only allow self-delete or admin
	tokenUserID, _ := c.Get(middleware.ContextUserID)
	tokenRole, _ := c.Get(middleware.ContextRole)
	if tokenUserID.(uint) != userId && tokenRole.(string) != string("admin") {
		HandleError(c, errs.Forbidden("You can only delete your own account"))
		return
	}

	err = u.service.Delete(userId)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "User deleted successfully",
	})
}

func (u *UserController) UpdateStatus(c *gin.Context) {
	userIdParam := c.Param("Id")
	var userId uint
	_, err := fmt.Sscanf(userIdParam, "%d", &userId)
	if err != nil {
		HandleError(c, errs.BadRequest("Invalid user ID format"))
		return
	}

	var payload dto.UserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleError(c, errs.BadRequest("Invalid request body"))
		return
	}

	payload.Id = userId

	response, err := u.service.UpdateStatus(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(response.StatusCode, response)
}

func (u *UserController) ResetPassword(c *gin.Context) {
	userIdParam := c.Param("Id")
	var userId uint
	_, err := fmt.Sscanf(userIdParam, "%d", &userId)
	if err != nil {
		HandleError(c, errs.BadRequest("Invalid user ID format"))
		return
	}

	tokenUserID, _ := c.Get(middleware.ContextUserID)
	tokenRole, _ := c.Get(middleware.ContextRole)
	if tokenUserID.(uint) != userId && tokenRole.(string) != string("admin") {
		HandleError(c, errs.Forbidden("You can only reset your own password account"))
		return
	}

	var payload dto.UserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleError(c, errs.BadRequest("Invalid request body"))
		return
	}

	payload.Id = userId

	response, err := u.service.ResetPassword(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(response.StatusCode, response)
}