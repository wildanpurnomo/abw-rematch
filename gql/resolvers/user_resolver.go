package gqlresolvers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	gqldataloaders "github.com/wildanpurnomo/abw-rematch/gql/dataloaders"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	loader = dataloader.NewBatchedLoader(gqldataloaders.UserBatchFn)
)

var GetUserByIdResolver = func(params graphql.ResolveParams) (interface{}, error) {
	source := params.Source.(models.Content)
	userId := fmt.Sprint(source.UserID)

	thunk := loader.Load(params.Context, dataloader.StringKey(userId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]interface{})
	for _, item := range result.([]models.User) {
		userMap[userId] = item
	}

	return userMap[userId], nil
}

var UpdatePasswordResolver = func(params graphql.ResolveParams) (interface{}, error) {
	cookieAccess := libs.GetContextValues(params.Context)
	userId := cookieAccess.UserID
	if userId == "0" {
		return false, errors.New("Invalid token or user not found")
	}

	var input models.UpdatePasswordInput

	// populate inputs
	input.OldPassword = params.Args["old_password"].(string)
	input.NewPassword = params.Args["new_password"].(string)

	// fetch user from DB
	var user models.User
	if err := repositories.Repo.FetchUserById(&user, userId); err != nil {
		return false, errors.New("Invalid token or user not found")
	}

	// verify password
	if !libs.VerifyPassword([]byte(user.Password), []byte(input.OldPassword)) {
		return false, errors.New("Invalid old password or new password")
	}

	// password validation
	if !libs.ValidatePassword(input.NewPassword) {
		return false, errors.New("Invalid old password or new password")
	}

	// hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
	if err != nil {
		return false, errors.New("Whoops")
	}

	// update password
	update := user
	update.Password = string(hash)
	if err := repositories.Repo.UpdateUser(&user, update); err != nil {
		return false, errors.New("Whoops")
	}

	return true, nil
}

var UpdateUsernameResolver = func(params graphql.ResolveParams) (interface{}, error) {
	cookieAccess := libs.GetContextValues(params.Context)
	userId := cookieAccess.UserID
	if userId == "0" {
		return nil, errors.New("Invalid token or user not found")
	}

	var input models.UpdateUsernameInput

	// trim username
	input.Username = strings.TrimSpace(params.Args["username"].(string))

	// username validation
	if !libs.ValidateUsername(input.Username) {
		return nil, errors.New("Invalid new username")
	}

	// fetch user from DB
	var user models.User
	if err := repositories.Repo.FetchUserById(&user, userId); err != nil {
		return nil, errors.New("Invalid token or user not found")
	}

	// update username
	update := user
	update.Username = input.Username
	if err := repositories.Repo.UpdateUser(&user, update); err != nil {
		return nil, errors.New("Whoops!")
	}

	return user, nil
}
