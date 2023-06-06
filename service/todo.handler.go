package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"todoGin/config"
	"todoGin/model/entity"
	"todoGin/model/request"
	"todoGin/model/respErr"
	"todoGin/repository"
)

type Handler struct {
	TodoRepository repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) *Handler {
	return &Handler{
		TodoRepository: todoRepo,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var user entity.User

	// binding request body ke struct user
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.Error{
			Error: "invalid request Body",
		})
		return
	}

	// cek apakah pengguna sudah ada di database
	existingUser, err := h.TodoRepository.GetUserByUsername(user.Username)
	if existingUser != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.Error{
			Error: "User already exist",
		})
		return
	}

	// hash password pengguna sebelum disimpan ke database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.Error{
			Error: "Failed hash Password",
		})
		return
	}

	// simpan pengguna ke database
	newUser := &entity.User{
		Username: user.Username,
		Password: string(hashedPassword),
	}
	err = h.TodoRepository.CreateUser(newUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.Error{
			Error: "Failed Create User",
		})
		return
	}

	// mengembalikan pesan berhasil sebagai response
	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h *Handler) Login(ctx *gin.Context) {
	var user entity.User

	// binding request body ke struct user
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.Error{
			Error: "invalid request Body",
		})
		return
	}

	// cek apakah pengguna ada di database
	storedUser, err := h.TodoRepository.GetUserByUsername(user.Username)
	if err != nil || storedUser == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, respErr.Error{
			Error: "invalid Username or Password",
		})
		return
	}

	// bandingkan password yang dimasukkan dengan hash password di database
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, respErr.Error{
			Error: "invalid Username or Password",
		})
		return
	}

	// membuat token
	token, err := config.CreateToken(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respErr.Error{
			Error: "Failed to generate Token",
		})
		return
	}

	//// menyimpan informasi pengguna yang sedang login ke dalam konteks Gin
	//ctx.Set("username", storedUser.Username)

	// Menampilkan pesan hello user dengan username yang berhasil login
	// mengembalikan token sebagai response
	ctx.JSON(http.StatusOK, gin.H{"" +
		"message": fmt.Sprintf("Hello %s! You are now logged in.", user.Username),
		"token": token,
	})
}

func (h *Handler) Access(ctx *gin.Context) {
	// ambil username dari konteks
	username, ok := ctx.Get("username")
	if !ok {
		// jika tidak ada username di dalam konteks, berarti pengguna belum terautentikasi
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, respErr.Error{
			Error: "Unauthorized",
		})
		return
	}

	// kirim pesan hello ke pengguna
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Hello %s!", username),
	})
}

func (h *Handler) TodolistHandlerGetAll(ctx *gin.Context) {
	todos, err := h.TodoRepository.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &respErr.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}
	logrus.Info(http.StatusOK, "Success Get All Data")
	//ctx.AbortWithStatusJSON(http.StatusOK, todos)
	ctx.AbortWithStatusJSON(http.StatusOK, request.TodoResponseToGetAll{
		Message: "Success Get All",
		Data:    len(todos),
		Todos:   todos,
	})
}
func (h *Handler) TodolistHandlerCreate(ctx *gin.Context) {
	todolist := new(request.TodolistCreateRequest)
	err := ctx.ShouldBindJSON(todolist)
	if err != nil {
		logrus.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "Invalid input",
			Status:  http.StatusBadRequest,
		})
		return
	}
	newTodo, errCreate := h.TodoRepository.Create(todolist.Title)
	if errCreate != nil {
		logrus.Error(errCreate)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}

	logrus.Info(http.StatusOK, " Success Create Todo", todolist)
	ctx.JSON(http.StatusOK, request.TodoResponse{
		Status:  http.StatusOK,
		Message: "New Todo Created",
		Data:    *newTodo,
	})
}
func (h *Handler) TodolistHandlerGetByID(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "Bad request",
			Status:  http.StatusBadRequest,
		})
		return
	}
	todo, err := h.TodoRepository.GetByID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todo by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	if todo == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, respErr.ErrorResponse{
			Message: "Not Found",
			Status:  http.StatusNotFound,
		})
		return
	}
	logrus.Info(http.StatusOK, " Success Get By ID")
	ctx.JSON(http.StatusOK, request.TodoResponse{
		Status:  http.StatusOK,
		Message: "Success Get Id",
		Data:    *todo,
	})
}

func (h *Handler) TodolistHandlerUpdate(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "parse ID error",
			Status:  http.StatusBadRequest,
		})
		return
	}
	reqBody := new(request.TodolistUpdateRequest)
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "Bad request",
			Status:  http.StatusBadRequest,
		})
		return
	}
	ErrId, err := h.TodoRepository.GetByID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todo by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	if ErrId == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, respErr.ErrorResponse{
			Message: "ID not Found",
			Status:  http.StatusNotFound,
		})
		return
	}
	rowsAffected, err := h.TodoRepository.Update(todoID, reqBody.ReqTodo())
	if err != nil {
		logrus.Errorf("failed when updating todo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	if rowsAffected == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, request.TodoIDResponse{
			Message: "Not Change",
			Data:    reqBody,
		})
		return
	}

	logrus.Info(http.StatusOK, " Success Update Todo")
	ctx.JSON(http.StatusOK, request.TodoUpdateResponse{
		Status:  http.StatusOK,
		Message: "Success Update Todo",
		Todos:   reqBody,
	})

}
func (h *Handler) TodolistHandlerDelete(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "Parse ID Error",
			Status:  http.StatusBadRequest,
		})
		return
	}
	isFound, err := h.TodoRepository.Delete(todoID)
	if err != nil {
		logrus.Errorf("failed when deleting todo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	//fmt.Println(isFound)
	if isFound == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, respErr.ErrorResponse{
			Message: "Not Found",
			Status:  http.StatusNotFound,
		})
		return
	}
	logrus.Info(http.StatusOK, " Success DELETE")
	ctx.JSON(http.StatusOK, request.TodoDeleteResponse{
		Status:  http.StatusOK,
		Message: "Success Delete",
	})
}
