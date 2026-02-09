package repositories

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"


	"axis/internal/models"
)

type UserRepo interface {
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID int) error
}

type userRepository struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewUserRepo(db *bun.DB, logger zerolog.Logger) UserRepo {
	return &userRepository{
		db:  db,
		log: logger,
	}
}

func (ur *userRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	user := new(models.User)

	err := ur.db.NewSelect().
		Model(user).
		Where("id = ?", userID).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			ur.log.Info().Int("user_id", userID).Msg("User not found.")
			return nil, err
		}
		ur.log.Error().Err(err).Int("user_id", userID).Msg("Failed to fetch user by ID.")
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)

	err := ur.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			ur.log.Info().Str("email", email).Msg("User not found by email.")
			return nil, err
		}
		ur.log.Error().Err(err).Str("email", email).Msg("Failed to fetch user by email.")
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := new(models.User)

	err := ur.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			ur.log.Info().Str("username", username).Msg("User not found by username.")
			return nil, err
		}
		ur.log.Error().Err(err).Str("username", username).Msg("Failed to fetch user by username.")
		return nil, err
	}

	return user, nil
}


func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := ur.db.NewInsert().
		Model(user).
		Exec(ctx)

	if err != nil {
		ur.log.Error().Err(err).Str("email", user.Email).Msg("Failed to create user.")
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := ur.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)

	if err != nil {
		ur.log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to update user.")
		return err
	}

	return nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, userID int) error {
	_, err := ur.db.NewDelete().
		Model(&models.User{}).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			ur.log.Info().Int("user_id", userID).Msg("User not found for deletion.")
			return err
		}
		ur.log.Error().Err(err).Int("user_id", userID).Msg("Failed to delete user.")
		return err
	}

	return nil
}
