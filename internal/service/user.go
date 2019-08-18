package service

import (
	"sync"
	"time"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Credential is model for user's authentication
type Credential struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// UserService user's service
type UserService struct {
	appConfig *conf.AppConfig
	appRepo   *model.AppRepo
}

// NewUserService creates a struct that implement IUserService
func NewUserService(appRepo *model.AppRepo, appConfig *conf.AppConfig) *UserService {
	return &UserService{
		appConfig,
		appRepo,
	}
}

// UpdateUser updates user
func (us *UserService) UpdateUser(id int64, fullName string, username string, email string, password string,
	verified *bool, active *bool) (*model.User, *comtype.CommonError) {
	userRepo := us.appRepo.UserRepo
	updatingMap := make(map[string]interface{})
	var (
		wg               sync.WaitGroup
		existingUserName *model.User
		existingEmail    *model.User
		hashedPassword   []byte

		userNameErr, emailErr *comtype.CommonError
		passwordErr           error
	)

	if len(username) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			existingUserName, userNameErr = userRepo.GetByUsername(username)
			if userNameErr == nil && existingUserName == nil {
				updatingMap["username"] = username
			}
		}()
	}

	if len(email) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			existingEmail, emailErr = userRepo.GetByEmail(email)
			if emailErr == nil && existingEmail == nil {
				updatingMap["email"] = email
			}
		}()

	}

	if len(password) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hashedPassword, passwordErr = bcrypt.GenerateFromPassword([]byte(password), us.appConfig.BcryptCost)
			if passwordErr == nil {
				updatingMap["hashed"] = string(hashedPassword)
			}
		}()
	}

	wg.Wait()

	if userNameErr != nil {
		return nil, comtype.NewCommonError(userNameErr, "UserService - UpdateUser:", comtype.ErrHandleDataFail,
			map[string]string{"username": "username is invalid"})
	}
	if emailErr != nil {
		return nil, comtype.NewCommonError(emailErr, "UserService - UpdateUser:", comtype.ErrHandleDataFail,
			map[string]string{"email": "email is invalid"})
	}
	if passwordErr != nil {
		return nil, comtype.NewCommonError(passwordErr, "UserService - UpdateUser:", comtype.ErrHandleDataFail,
			map[string]string{"password": "password is invalid"})
	}

	if existingUserName != nil {
		return nil, comtype.NewCommonError(nil, "UserService - CreateUser",
			comtype.ErrDuplicatedData, map[string]string{"username": "username is duplicated"})
	}
	if existingEmail != nil {
		return nil, comtype.NewCommonError(nil, "UserService - CreateUser",
			comtype.ErrDuplicatedData, map[string]string{"email": "email is duplicated"})
	}

	if len(fullName) > 0 {
		updatingMap["full_name"] = fullName
	}

	if active != nil {
		updatingMap["active"] = *active
	}

	if verified != nil {
		updatingMap["verified"] = *verified
	}

	updatingMap["updated_at"] = time.Now()

	err := userRepo.Update(id, updatingMap)
	if err != nil {
		return nil, err
	}

	return userRepo.GetByID(id)
}

// GetUser gets an user by ID
func (us *UserService) GetUser(id int64) (*model.User, *comtype.CommonError) {
	userRepo := us.appRepo.UserRepo
	return userRepo.GetByID(id)
}

// FetchUsers gets a list of users
func (us *UserService) FetchUsers(take int, fullName string, username string, verified *bool, email string,
	active *bool, sortBy string) ([]*model.User, *comtype.CommonError) {
	userRepo := us.appRepo.UserRepo
	filters := make(map[string]interface{})
	sorts := make(map[string]comtype.SortDirection)

	if take == 0 || take > 100 {
		take = 100
	}

	if len(username) > 0 {
		filters["username"] = username
	}

	if len(fullName) > 0 {
		filters["full_name"] = fullName
	}

	if len(email) > 0 {
		filters["email"] = email
	}

	if verified != nil {
		filters["verified"] = *verified
	}

	if active != nil {
		filters["active"] = *active
	}

	if len(sortBy) > 0 {
		switch sortBy[0] {
		case '+':
			sorts[sortBy[1:]] = comtype.Ascending
			break
		case '-':
			sorts[sortBy[1:]] = comtype.Decending
			break
		default:
			sorts["created_at"] = comtype.Decending
			break
		}
	} else {
		sorts["created_at"] = comtype.Decending
	}

	return userRepo.Query(take, filters, sorts)
}

// VerifyLogin check if username and password are correct or not
func (us *UserService) VerifyLogin(username string, password string) (*model.User, []*model.Action, []*model.Role, *comtype.CommonError) {
	userRepo := us.appRepo.UserRepo

	user, err := userRepo.GetByUsername(username)
	if err != nil {
		return nil, nil, nil, err
	}
	if user == nil {
		return nil, nil, nil,
			comtype.NewCommonError(nil, "UserService - VerifyLogin:", comtype.ErrDataNotFound, nil)
	}

	if !us.IsPasswordMatched(user.Hashed, password) {
		return nil, nil, nil,
			comtype.NewCommonError(nil, "UserService - VerifyLogin:", comtype.ErrInvalidCredential, nil)
	}

	var (
		wg      sync.WaitGroup
		actions []*model.Action
		roles   []*model.Role

		errAction *comtype.CommonError
		errRole   *comtype.CommonError
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		actions, errAction = us.GetUserActions(user.ID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		roles, errRole = us.GetUserRoles(user.ID)
	}()

	wg.Wait()

	if errAction != nil {
		return nil, nil, nil, errAction
	}
	if errRole != nil {
		return nil, nil, nil, errRole
	}

	return user, actions, roles, nil
}

// RegisterUser method to add a new user
func (us *UserService) RegisterUser(fullName string, username string, password string, email string) (*model.User, *comtype.CommonError) {
	userRepo := us.appRepo.UserRepo
	type dupResult struct {
		err   *comtype.CommonError
		isDup bool
	}
	dupUsernameChan := make(chan dupResult)
	dupEmailChan := make(chan dupResult)

	go func() {
		foundUsername, err := userRepo.GetByUsername(username)
		dupUsernameChan <- dupResult{
			err:   err,
			isDup: foundUsername != nil,
		}
	}()

	go func() {
		foundEmail, err := userRepo.GetByEmail(email)
		dupEmailChan <- dupResult{
			err:   err,
			isDup: foundEmail != nil,
		}
	}()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), us.appConfig.BcryptCost)
	if err != nil {
		return nil, comtype.NewCommonError(err, "UserService - RegisterUser:", comtype.ErrHandleDataFail, nil)
	}

	if result := <-dupUsernameChan; result.err != nil {
		return nil, result.err
	} else if result.isDup {
		return nil, comtype.NewCommonError(nil, "UserService - RegisterUser:",
			comtype.ErrDuplicatedData, map[string]string{
				"username": "duplicated username",
			})
	}

	if result := <-dupEmailChan; result.err != nil {
		return nil, result.err
	} else if result.isDup {
		return nil, comtype.NewCommonError(nil, "UserService - RegisterUser:",
			comtype.ErrDuplicatedData, map[string]string{
				"email": "duplicated email",
			})
	}

	userID, createdErr := userRepo.Create(fullName, username, string(hashedPassword), email)
	if err != nil {
		return nil, createdErr
	}

	user, getErr := userRepo.GetByID(userID)
	if err != nil {
		return nil, getErr
	}

	return user, nil
}

// GetUserActions gets all user's actions
func (us *UserService) GetUserActions(userID int64) ([]*model.Action, *comtype.CommonError) {
	actions, err := us.appRepo.ActionRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return actions, nil
}

// GetUserRoles gets all user's roles
func (us *UserService) GetUserRoles(userID int64) ([]*model.Role, *comtype.CommonError) {
	roles, err := us.appRepo.RoleRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// IsPasswordMatched check if passwords are matched
func (us *UserService) IsPasswordMatched(hashed string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
