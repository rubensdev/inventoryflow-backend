package user

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rubensdev/inventoryflow-backend/internal/jsonutil"
)

const UserCreatedMsg = "User created successfully"
const UserDeletedMsg = "User deleted successfully"

type UserHandler struct {
	logger  *log.Logger
	userSrv UserService
}

func NewUserHandler(logger *log.Logger, userService UserService) *UserHandler {
	return &UserHandler{
		logger:  logger,
		userSrv: userService,
	}
}

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userSrv.GetAll()

	jsonRes := jsonutil.NewJSONResponse(h.logger)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.H{"users": users}, nil)
	if err != nil {
		jsonRes.ServerErrorResponse(w, r, err)
	}
}

func (h UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	jsonRes := jsonutil.NewJSONResponse(h.logger)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		jsonRes.NotFoundResponse(w, r)
		return
	}

	user, err := h.userSrv.GetByID(id)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
	}
	if user == nil {
		jsonRes.NotFoundResponse(w, r)
		return
	}
	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.H{"user": user}, nil)
	if err != nil {
		jsonRes.ServerErrorResponse(w, r, err)
	}
}

func (h UserHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	jsonRes := jsonutil.NewJSONResponse(h.logger)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		jsonRes.NotFoundResponse(w, r)
		return
	}

	userUpdateReq := &UserUpdateRequest{}

	err = jsonutil.ReadJSON(w, r, &userUpdateReq)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	valid := userUpdateReq.Validate()
	if !valid {
		err = jsonutil.WriteJSON(w, http.StatusBadRequest, jsonutil.H{"errors": userUpdateReq.GetErrors()}, nil)
		if err != nil {
			jsonRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	user, err := h.userSrv.Update(id, userUpdateReq)
	if err != nil {
		switch {
		case errors.Is(err, ErrEditConflict):
			jsonRes.EditConflictResponse(w, r)
		case errors.Is(err, ErrUserNotFound):
			jsonRes.NotFoundResponse(w, r)
		default:
			jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
		}
		return
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.H{"user": user}, nil)
	if err != nil {
		jsonRes.ServerErrorResponse(w, r, err)
	}
}

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userRegisterReq := &UserRegisterRequest{}

	jsonRes := jsonutil.NewJSONResponse(h.logger)

	err := jsonutil.ReadJSON(w, r, &userRegisterReq)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	valid := userRegisterReq.Validate()
	if !valid {
		err = jsonutil.WriteJSON(w, http.StatusBadRequest, jsonutil.H{
			"errors": userRegisterReq.GetErrors(),
		}, nil)
		if err != nil {
			jsonRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	u := userRegisterReq.GetModel()
	err = h.userSrv.Create(u)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	data := jsonutil.H{"message": UserCreatedMsg, "user": u}
	err = jsonutil.WriteJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		jsonRes.ServerErrorResponse(w, r, err)
	}
}

func (h UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	jsonRes := jsonutil.NewJSONResponse(h.logger)

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		jsonRes.NotFoundResponse(w, r)
		return
	}

	err = h.userSrv.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			jsonRes.NotFoundResponse(w, r)
		default:
			jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
		}
		return
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.H{"message": UserDeletedMsg}, nil)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusInternalServerError, err)
	}
}
