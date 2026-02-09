package services

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
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
	log      zerolog.Logger
}

func NewUserService(userRepo repositories.UserRepo, logger zerolog.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		log:      logger,
	}
}


func (s *userService) Register(ctx context.Context, form models.RegisterModel) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error().Err(err).Msg("Can't hash the password")
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
		s.log.Info().Str("email", user.Email).Msg("User already exists")
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, creds models.LoginModel) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, creds.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Str("email", creds.Email).Msg("User not found during login.")
			return nil, err
		}
		s.log.Error().Err(err).Str("email", creds.Email).Msg("Failed to fetch user by email for login")
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		s.log.Error().Err(err).Str("email", creds.Email).Msg("Login with Invalid password.")
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("user_id", id).Msg("User not found.")
			return nil, err
		}
		s.log.Error().Err(err).Int("user_id", id).Msg("Failed to fetch user by ID.")
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Str("email", email).Msg("User not found by email.")
			return nil, err
		}
		s.log.Error().Err(err).Str("email", email).Msg("Failed to fetch user by email.")
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Str("username", username).Msg("User not found by username.")
			return nil, err
		}
		s.log.Error().Err(err).Str("username", username).Msg("Failed to fetch user by username.")
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID int, newUser models.UpdateUser) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("user_id", userID).Msg("User not found for update.")
			return nil, err
		}
		s.log.Error().Err(err).Int("user_id", userID).Msg("Failed to fetch user for update.")
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
		s.log.Error().Err(err).Int("user_id", userID).Msg("Failed to update user.")
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("user_id", id).Msg("User not found for deletion.")
			return err
		}
		s.log.Error().Err(err).Int("user_id", id).Msg("Failed to delete user.")
		return err
	}
	return nil
}
