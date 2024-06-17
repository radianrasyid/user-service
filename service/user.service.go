package service

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prithuadhikary/user-service/domain"
	"github.com/prithuadhikary/user-service/helper"
	"github.com/prithuadhikary/user-service/model"
	"github.com/prithuadhikary/user-service/repository"
)

var secretKey = []byte("radianrasyid")

type UserService interface {
	Signup(request *model.SignupRequest) error
	Signin(request *model.SigninRequest) (string, *domain.Session, error)
	IsUserExist(username string) (bool, error)
	Signout(request *model.Signout) error
	Whoami(request *model.Whoami) (*model.WhoamiResponse, error)
	EditUser(request *model.EditUserRequest, currentUser *model.WhoamiResponse) (*model.EditUserRequest, error)
}

type userService struct {
	repository repository.UserRepository
}

// Signin implements UserService.
func (service *userService) Signin(request *model.SigninRequest) (string, *domain.Session, error) {
	exists := service.repository.ExistsByUsername(request.Username)
	if !exists {
		return "", nil, errors.New("username and password might be wrong")
	}

	currentUser, err := service.repository.FindSpecificUsername(request.Username)
	if err != nil {
		return "", nil, err // Pass the error returned by FindSpecificUsername directly
	}

	passwordSame := helper.CheckPasswordHash(request.Password, currentUser.Password)
	fmt.Println(passwordSame)
	if !passwordSame {
		return "", nil, errors.New("username and password might be wrong")
	}

	token, err := service.CreateToken(&currentUser)
	fmt.Println("ini error di service sign in", err)
	if err != nil {
		return "", nil, errors.New("something went wrong when creating your token")
	} // Return nil if the login is successful
	service.repository.EditUser(currentUser.Username, token, "token")

	session, err := service.CreateSession(currentUser.ID)

	if err != nil {
		return "", nil, errors.New("something went wrong when creating your session")
	}

	fmt.Println("ini session", &session)

	return currentUser.ID.String(), session, nil

}

func (service *userService) CreateToken(request *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": request.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"role":     request.Role,
	})

	tokenString, err := token.SignedString(secretKey)
	fmt.Println("ini error di service createToken", err)
	if err != nil {
		return "", errors.New("something went wrong when creating token")
	}

	return tokenString, nil
}

func (service *userService) Signup(request *model.SignupRequest) error {
	if request.Password != request.PasswordConfirmation {
		return errors.New("password and confirm password must match")
	}
	exists := service.repository.ExistsByUsername(request.Username)
	if exists {
		return errors.New("email already exists")
	}
	service.repository.Save(&domain.User{
		Username: request.Username,
		Password: request.Password,
		Role:     "END_USER",
		Email:    request.Email,
	})
	return nil
}

func (service *userService) Signout(request *model.Signout) error {
	err := service.VerifyToken(request.Jwt)

	if err != nil {
		return errors.New("token is not verified")
	}

	token, err := jwt.Parse(request.Jwt, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return errors.New("can not verify jwt")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("ini isi jwt", claims["username"])
		service.repository.EditUser(claims["username"].(string), nil, "token")
		return nil
	} else {
		log.Printf("invalid jwt token")
		return nil
	}
}

func (service *userService) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (service *userService) CreateSession(userID uuid.UUID) (*domain.Session, error) {
	session := &domain.Session{
		ID:        uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := service.repository.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (service *userService) Whoami(request *model.Whoami) (*model.WhoamiResponse, error) {
	SessionID, err := uuid.Parse(request.SessionID)
	if err != nil {
		return nil, err
	}

	result, err := service.repository.FindUserBySessionID(SessionID)

	if err != nil {
		return nil, err
	}

	return &model.WhoamiResponse{
		ID:       result.ID,
		Username: result.Username,
		Role:     result.Role,
		Email:    result.Email,
		Token:    result.Token,
	}, nil
}

func (service *userService) EditUser(request *model.EditUserRequest, currentUser *model.WhoamiResponse) (*model.EditUserRequest, error) {

	currentRequest := &model.EditUserRequest{
		Username: request.Username,
		Email:    request.Email,
	}
	fmt.Print("ini request yang masuk ke edit user")
	fmt.Println(currentRequest)
	fmt.Print("ini data current user")
	fmt.Println(currentUser)

	v := reflect.ValueOf(request).Elem()

	typeOfCurrentRequest := v.Type()
	var loopErrors []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typeOfCurrentRequest.Field(i).Name

		var fieldValue interface{}

		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				fieldValue = field.Elem().Interface()
			}
		} else {
			fieldValue = field.Interface()
		}

		if fieldValue != nil {
			err := service.repository.EditUser(currentUser.Username, fieldValue, strings.ToLower(fieldName))

			if err != nil {
				loopErrors = append(loopErrors, err.Error())
			}
		}
	}

	if len(loopErrors) > 0 {
		return nil, fmt.Errorf("errors occured: %v", loopErrors)
	}

	return request, nil
}

func (service *userService) IsUserExist(username string) (bool, error) {
	err := service.repository.ExistsByUsername(username)
	fmt.Print("this query from exist by username")
	fmt.Println(err)
	if err {
		return false, errors.New(fmt.Sprintf("problem occured when checking validity of username %v", err))
	}

	return true, nil
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}
