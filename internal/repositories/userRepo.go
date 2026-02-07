package repositories

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	"backend/internal/models"
)

type UserRepo interface {
	GetUser(ctx context.Context, userID int) (*models.User, error)
	Login(ctx context.Context, creds models.LoginModel) (*models.User, error)
	Register(ctx context.Context, form models.RegisterModel) (*models.User, error)
	UpdateUser(ctx context.Context, userID int, newUser models.UpdateUser) (*models.User, error)
	DeleteUser(ctx context.Context, userID int) error
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepo(db *bun.DB) UserRepo {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) GetUser(ctx context.Context, userID int) (*models.User, error) {
	user := new(models.User)

	err := ur.db.NewSelect().
		Model(user).
		Where("id = ?", userID).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msg("User not found.")
			return nil, err
		}
		log.Error().Err(err).Msg("Failed to fetch user.")
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) Login(ctx context.Context, creds models.LoginModel) (*models.User, error) {
	user := new(models.User)

	err := ur.db.NewSelect().
		Model(user).
		Where("email = ?", creds.Email).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msg("User not found.")
			return nil, err
		}
		log.Error().Err(err).Msg("Failed to login user")
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Error().Err(err).Msg("Login with Invalid password.")
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) Register(ctx context.Context, form models.RegisterModel) (*models.User, error) {
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

	_, err = ur.db.NewInsert().
		Model(user).
		Exec(ctx)

	if err != nil {
		log.Info().Str("email", user.Email).Msg("User already exists")
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, userID int, newUser models.UpdateUser) (*models.User, error) {
	user := new(models.User)

	err := ur.db.NewSelect().
		Model(user).
		Where("id = ?", userID).
		Scan(ctx)

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

	_, err = ur.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)

	if err != nil {
		log.Error().Err(err).Msgf("Failed to update user with ID %d.", userID)
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, userID int) error {
	_, err := ur.db.NewDelete().
		Model(&models.User{}).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msgf("User with ID %d not found for deletion.", userID)
			return err
		}
		log.Error().Err(err).Msgf("Failed to delete user with ID %d.", userID)
		return err
	}

	return nil
}
