package services

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"axis/internal/repositories"
	"axis/internal/models"
)

type UserService interface {
	Register(ctx context.Context, form models.RegisterModel) (*models.User, error)
	Login(ctx context.Context, creds models.LoginModel) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, userID int, newUser models.UpdateUser) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	userRepo repositories.UserRepo
}

func NewUserService(userRepo repositories.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}


func (s *userService) Register(ctx context.Context, form models.RegisterModel) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Can't hash the password")
		return nil, err
	}

	user := &models.User{
		Name:        form.Name,
		Username:    form.Username,
		Email:       form.Email,
		Password:    string(hashedPassword),
		Status:      models.Active,
		Timezone:    form.Timezone,
		Locale:      form.Locale,
		IsVerified:  false,
		LastLoginAt: nil,
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Info().Str("email", user.Email).Msg("User already exists")
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, creds models.LoginModel) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, creds.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msg("User not found during login.")
			return nil, err
		}
		log.Error().Err(err).Msg("Failed to fetch user by email for login")
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Error().Err(err).Msg("Login with Invalid password.")
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msgf("User with ID %d not found.", id)
			return nil, err
		}
		log.Error().Err(err).Msgf("Failed to fetch user by ID %d.", id)
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msgf("User with email %s not found.", email)
			return nil, err
		}
		log.Error().Err(err).Msgf("Failed to fetch user by email %s.", email)
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msgf("User with username %s not found.", username)
			return nil, err
		}
		log.Error().Err(err).Msgf("Failed to fetch user by username %s.", username)
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID int, newUser models.UpdateUser) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msgf("User with ID %d not found for update.", userID)
			return nil, err
		}
		log.Error().Err(err).Msgf("Failed to fetch user with ID %d for update.", userID)
		return nil, err
	}

	if newUser.Name != nil {
		user.Name = *newUser.Name
	}
	if newUser.Username != nil {
		user.Username = *newUser.Username
	}
	if newUser.Email != nil {
		user.Email = *newUser.Email
	}
	if newUser.Timezone != nil {
		user.Timezone = *newUser.Timezone
	}
	if newUser.Locale != nil {
		user.Locale = *newUser.Locale
	}

	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to update user with ID %d.", userID)
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msgf("User with ID %d not found for deletion.", id)
			return err
		}
		log.Error().Err(err).Msgf("Failed to delete user with ID %d.", id)
		return err
	}
	return nil
}
