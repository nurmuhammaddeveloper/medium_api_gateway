package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nurmuhammaddeveloper/medium_api_gateway/api/models"

	pbu "github.com/nurmuhammaddeveloper/medium_api_gateway/genproto/user_service"
)

// @Router /users [post]
// @Summary Create a user
// @Description Create a user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User"
// @Success 201 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		req models.CreateUserRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := h.grpcClient.UserService().Create(context.Background(), &pbu.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Gender:          req.Gender,
		Password:        req.Password,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
		Type:            req.Type,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseUserModel(user))
}

// @Router /users [put]
// @Summary Update a user
// @Description Update a user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.UpdateUserRequest true "User"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		req models.UpdateUserRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, err := h.grpcClient.UserService().Update(context.Background(), &pbu.User{
		Id:              int64(id),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Gender:          req.Gender,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseUserModel(user))
}

// @Router /users/{id} [get]
// @Summary Get user by id
// @Description Get user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.grpcClient.UserService().Get(context.Background(), &pbu.IdRequest{Id: int64(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseUserModel(resp))
}

// @Router /users/email/{email} [get]
// @Summary Get user by email
// @Description Get user by email
// @Tags user
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	resp, err := h.grpcClient.UserService().GetByEmail(context.Background(), &pbu.GetByEmailRequest{Email: email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseUserModel(resp))
}

func parseUserModel(user *pbu.User) models.User {
	return models.User{
		ID:              user.Id,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt,
	}
}

// @Router /users [get]
// @Summary Get all users
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 200 {object} models.GetAllUsersResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	req, err := validateGetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.grpcClient.UserService().GetAll(context.Background(), &pbu.GetAllUsersRequest{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getUsersResponse(result))
}

func getUsersResponse(data *pbu.GetAllUsersResponse) *models.GetAllUsersResponse {
	response := models.GetAllUsersResponse{
		Users: make([]*models.User, 0),
		Count: data.Count,
	}

	for _, user := range data.Users {
		u := parseUserModel(user)
		response.Users = append(response.Users, &u)
	}

	return &response
}

// @Router /users/{id} [delete]
// @Summary Delete user
// @Description Delete user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = h.grpcClient.UserService().Delete(context.Background(), &pbu.IdRequest{Id: int64(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.ResponseOK{
		Message: "success",
	})
}

// func(h *handlerV1)GetProfile(c *gin.Context){
// 	payload, err := h.GetAuthPayload(c)
// 	if err != nil{
// 		c.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// }
