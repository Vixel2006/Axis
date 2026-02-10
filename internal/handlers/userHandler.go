package handlers

import (
	"net/http"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userService services.UserService
	log         zerolog.Logger
}

func NewUserHandler(us services.UserService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		userService: us,
		log:         logger,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	h.log.Info().Msg("Handling Register request")
	var form models.RegisterModel
	if err := c.ShouldBindJSON(&form); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for Register")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Str("email", form.Email).Str("username", form.Username).Msg("Register request body")

	user, err := h.userService.Register(c.Request.Context(), form)
	if err != nil {
		h.log.Error().Err(err).Str("email", form.Email).Msg("Failed to register user via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	h.log.Info().Int("user_id", user.ID).Str("email", user.Email).Msg("User registered successfully")
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	h.log.Info().Msg("Handling Login request")
	var creds models.LoginModel
	if err := c.ShouldBindJSON(&creds); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for Login")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Str("email", creds.Email).Msg("Login request body (email)")

	user, err := h.userService.Login(c.Request.Context(), creds)
	if err != nil {
		h.log.Warn().Err(err).Str("email", creds.Email).Msg("Login failed: invalid credentials or user not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		h.log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to generate token for login")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	h.log.Info().Int("user_id", user.ID).Str("email", user.Email).Msg("User logged in successfully")
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	h.log.Info().Msg("Handling GetUserByID request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in GetUserByID")
		return
	}
	h.log.Debug().Int("user_id_from_context", userID).Msg("Retrieving user by ID")

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("user_id", userID).Msg("User not found by ID")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.log.Error().Err(err).Int("user_id", userID).Msg("Failed to retrieve user via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if user == nil {
		h.log.Info().Int("user_id", userID).Msg("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	h.log.Info().Int("user_id", userID).Msg("User retrieved successfully")
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	h.log.Info().Msg("Handling GetUserByEmail request")
	email := c.Query("email")
	if email == "" {
		h.log.Warn().Msg("Email query parameter is required for GetUserByEmail")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email query parameter is required"})
		return
	}
	h.log.Debug().Str("email_param", email).Msg("Retrieving user by email")

	user, err := h.userService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Str("email", email).Msg("User not found by email")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.log.Error().Err(err).Str("email", email).Msg("Failed to retrieve user by email via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if user == nil {
		h.log.Info().Str("email", email).Msg("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	h.log.Info().Str("email", email).Msg("User retrieved successfully by email")
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	h.log.Info().Msg("Handling GetUserByUsername request")
	username := c.Query("username")
	if username == "" {
		h.log.Warn().Msg("Username query parameter is required for GetUserByUsername")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username query parameter is required"})
		return
	}
	h.log.Debug().Str("username_param", username).Msg("Retrieving user by username")

	user, err := h.userService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Str("username", username).Msg("User not found by username")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.log.Error().Err(err).Str("username", username).Msg("Failed to retrieve user by username via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if user == nil {
		h.log.Info().Str("username", username).Msg("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	h.log.Info().Str("username", username).Msg("User retrieved successfully by username")
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	h.log.Info().Msg("Handling UpdateUser request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in UpdateUser")
		return
	}
	h.log.Debug().Int("user_id_from_context", userID).Msg("Updating user")

	var updateUser models.UpdateUser
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for UpdateUser")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Int("user_id", userID).Interface("update_data", updateUser).Msg("UpdateUser request body")

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, updateUser)
	if err != nil {
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("user_id", userID).Msg("User not found for update")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.log.Error().Err(err).Int("user_id", userID).Msg("Failed to update user via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	if user == nil {
		h.log.Info().Int("user_id", userID).Msg("User not found for update (service returned nil)")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	h.log.Info().Int("user_id", userID).Msg("User updated successfully")
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	h.log.Info().Msg("Handling DeleteUser request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in DeleteUser")
		return
	}
	h.log.Debug().Int("user_id_from_context", userID).Msg("Deleting user")

	err = h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("user_id", userID).Msg("User not found for deletion")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.log.Error().Err(err).Int("user_id", userID).Msg("Failed to delete user via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	h.log.Info().Int("user_id", userID).Msg("User deleted successfully")
	c.JSON(http.StatusNoContent, nil)
}
