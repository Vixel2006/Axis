package repositories

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog/log"
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
	db *bun.DB
}

func NewUserRepo(db *bun.DB) UserRepo {
	return &userRepository{
		db: db,
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
			log.Info().Msg("User not found.")
			return nil, err
		}
		log.Error().Err(err).Msg("Failed to fetch user.")
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
			log.Info().Msgf("User with email %s not found.", email)
			return nil, err
		}
		log.Error().Err(err).Msgf("Failed to fetch user by email %s.", email)
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
			log.Info().Msgf("User with username %s not found.", username)
			return nil, err
		}
		log.Error().Err(err).Msgf("Failed to fetch user by username %s.", username)
		return nil, err
	}

	return user, nil
}





func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := ur.db.NewInsert().
		Model(user).
		Exec(ctx)

	if err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to create user.")
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
		log.Error().Err(err).Msgf("Failed to update user with ID %d.", user.ID)
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
			log.Info().Msgf("User with ID %d not found for deletion.", userID)
			return err
		}
		log.Error().Err(err).Msgf("Failed to delete user with ID %d.", userID)
		return err
	}

	return nil
}
