package handler

import (
	"net/http"

	"strings"
	"encoding/json"

	"github.com/Yash-Watchguard/Tasknest/internal/logger"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"

	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	
)

type UserHandler struct {
	userService service1.UserServiceInterface
}

func NewUserHandler(userService service1.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}
func (uh *UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        uh.Getuser(w, r)
    case http.MethodDelete:
        uh.DeleteUser(w, r)
	case http.MethodPut:
        uh.PromoteEmployee(w,r)
    case http.MethodPatch:
		uh.UpdateUser(w,r)
	default:
		logger.Error("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
    }


}
func(uh *UserHandler)Getuser(w http.ResponseWriter,r *http.Request){

   

    ctx := r.Context()
    userID, ok := ctx.Value(ContextKey.UserId).(string)
    if !ok {
		logger.Error("user id not found")
        response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
        return
    }

    role, ok := ctx.Value(ContextKey.UserRole).(roles.Role)
    if !ok {
		logger.Error("user role not found")
        response.ErrorResponse(w, http.StatusUnauthorized, "User role not found", 1007)
        return
    }
	 // Get the query parameter "id"
    path := strings.TrimPrefix(r.URL.Path, "/v1/users/")
    id := strings.Trim(path, "/")

    if id != "" {
        // Fetch single user
        if role != roles.Admin && id != userID {
			logger.Error("unauthorized person wants to view profile")
            response.ErrorResponse(w, http.StatusForbidden, "Access denied", 1008)
            return
        }

        user, err := uh.userService.ViewProfile(id)
        if err != nil {
			logger.Error("user not found")
            response.ErrorResponse(w, http.StatusNotFound, "User not found", 1004)
            return
        }
        
		logger.Info("User retrieved successfully")
        response.SuccessResponse(w, user, "User retrieved successfully", http.StatusOK)
        return
    }

	if role!=roles.Admin{
		logger.Error("unauthorized person wants to view profile")
        response.ErrorResponse(w, http.StatusForbidden, "Access denied", 1008)
        return
	}
    users, err :=uh.userService.ViewAllUsers()
        if err != nil {
			logger.Error("no users available")
            response.ErrorResponse(w, http.StatusNotFound, "No user found", 1004)
            return
        }

		logger.Info("users retrived successfully")
		response.SuccessResponse(w,users,"Users Retrived Successfully",http.StatusOK)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        logger.Error("Invalid HTTP method")
        response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
        return
    }

    ctx := r.Context()

    path := strings.TrimPrefix(r.URL.Path, "/v1/users/")
    id := strings.Trim(path, "/")
    userId := ctx.Value(ContextKey.UserId).(string)
    role := ctx.Value(ContextKey.UserRole).(roles.Role)

    
    if userId != id && role != roles.Admin {
        logger.Error("Unauthorized delete attempt")
        response.ErrorResponse(w, http.StatusForbidden, "Access denied", 1008)
        return
    }

    err := uh.userService.DeleteUser(id)
    if err != nil {
        logger.Error("Error deleting user")
        response.ErrorResponse(w, http.StatusInternalServerError, "Error deleting user", 1010)
        return
    }

    logger.Info("User deleted successfully")
    response.SuccessResponse(w, nil, "User deleted successfully", http.StatusOK)
}

func(uh *UserHandler)PromoteEmployee(w http.ResponseWriter,r *http.Request){
    
	ctx := r.Context()
    role := ctx.Value(ContextKey.UserRole).(roles.Role)

	if role != roles.Admin {
		logger.Error("Only admins can promote users")
        response.ErrorResponse(w, http.StatusForbidden, "Only admins can promote users", 1008)
        return
    }

	path :=strings.TrimPrefix(r.URL.Path, "/v1/users/")

	parts:=strings.Split(path,"/")

	if len(parts)!=2 || parts[1]!="promote"{
		logger.Error("Invalid path")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid path", 1001)
        return
	}
	id := parts[0]

	err:=uh.userService.PromoteEmployee(id)

	if err!=nil{
		logger.Error("Failed to promote user")
		response.ErrorResponse(w,http.StatusInternalServerError,"Failed to promote user",1011)
	}

	logger.Info("User promoted to Manager successfully")
	response.SuccessResponse(w,nil,"User promoted to Manager successfully",http.StatusOK)
}

func(uh *UserHandler)UpdateUser(w http.ResponseWriter,r * http.Request){
    
	ctx := r.Context()
    userId := ctx.Value(ContextKey.UserId).(string)
    role := ctx.Value(ContextKey.UserRole).(roles.Role)

	path := strings.TrimPrefix(r.URL.Path, "/v1/users/")
    id := strings.Trim(path, "/")

	if userId != id && role != roles.Admin {
        logger.Error("Unauthorized update attempt")
        response.ErrorResponse(w, http.StatusForbidden, "Access denied",1008)
        return
    }
    var updates map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        logger.Error("Invalid JSON body")
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body",1001)
        return
    }

	err:=uh.userService.UpdateUser(userId,updates)
    
	if err != nil {
        logger.Error("Failed to update user: " + err.Error())
        response.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user",1011)
        return
    }

    logger.Info("User updated successfully")
    response.SuccessResponse(w, nil, "User updated successfully", http.StatusOK)
    
}






