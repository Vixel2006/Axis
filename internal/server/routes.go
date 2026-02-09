package server

import (
	"net/http"


	"axis/internal/handlers"
	"axis/internal/middlewares"
	"axis/internal/repositories"
	"axis/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	// Access the bun.DB from the database service
	bunDB := s.db.GetDB()

	// --- Repositories ---
	attachmentRepo := repositories.NewAttachmentRepo(bunDB)
	channelMemberRepo := repositories.NewChannelMemberRepo(bunDB)
	channelRepo := repositories.NewChannelRepo(bunDB)
	messageRepo := repositories.NewMessageRepo(bunDB)
	reactionRepo := repositories.NewReactionRepo(bunDB)
	userRepo := repositories.NewUserRepo(bunDB)
	workspaceMemberRepo := repositories.NewWorkspaceMemberRepo(bunDB)
	workspaceRepo := repositories.NewWorkspaceRepo(bunDB)

	// --- Services ---
	attachmentService := services.NewAttachmentService(attachmentRepo)
	channelMemberService := services.NewChannelMemberService(channelMemberRepo)
	channelService := services.NewChannelService(channelRepo, channelMemberRepo, workspaceMemberRepo)
	messageService := services.NewMessageService(messageRepo)
	reactionService := services.NewReactionService(reactionRepo)
	userService := services.NewUserService(userRepo)
	workspaceMemberService := services.NewWorkspaceMemberService(workspaceMemberRepo)
	workspaceService := services.NewWorkspaceService(workspaceRepo, workspaceMemberRepo)

	// --- Handlers ---
	attachmentHandler := handlers.NewAttachmentHandler(attachmentService)
	channelMemberHandler := handlers.NewChannelMemberHandler(channelMemberService)
	channelHandler := handlers.NewChannelHandler(channelService)
	messageHandler := handlers.NewMessageHandler(messageService)
	reactionHandler := handlers.NewReactionHandler(reactionService)
	userHandler := handlers.NewUserHandler(userService)
	workspaceMemberHandler := handlers.NewWorkspaceMemberHandler(workspaceMemberService)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceService)

	// --- API Routes ---
	api := r.Group("/api")
	{
		// User Routes
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/users", middlewares.JWTAuth(), userHandler.GetUserByID)
		api.GET("/users/by-email", userHandler.GetUserByEmail)       // Query param: ?email=
		api.GET("/users/by-username", userHandler.GetUserByUsername) // Query param: ?username=
		api.PUT("/users", middlewares.JWTAuth(), userHandler.UpdateUser)
		api.DELETE("/users", middlewares.JWTAuth(), userHandler.DeleteUser)

		// Workspace Routes
		api.POST("/workspaces", middlewares.JWTAuth(), workspaceHandler.CreateWorkspace)
		api.GET("/workspaces/:workspaceID", workspaceHandler.GetWorkspaceByID)
		api.PUT("/workspaces/:workspaceID", middlewares.JWTAuth(), workspaceHandler.UpdateWorkspace)
		api.DELETE("/workspaces/:workspaceID", middlewares.JWTAuth(), workspaceHandler.DeleteWorkspace)
		api.GET("/workspaces", middlewares.JWTAuth(), workspaceHandler.GetWorkspacesForUser)

		// Workspace Member Routes
		api.POST("/workspaces/:workspaceID/members", workspaceMemberHandler.AddMemberToWorkspace)
		api.DELETE("/workspaces/:workspaceID/members/:userID", workspaceMemberHandler.RemoveMemberFromWorkspace)
		api.GET("/workspaces/:workspaceID/members", workspaceMemberHandler.GetWorkspaceMembers)

		// Channel Routes
		api.POST("/channels", middlewares.JWTAuth(), channelHandler.CreateChannel)
		api.GET("/channels/:channelID", channelHandler.GetChannelByID)
		api.PUT("/channels/:channelID", middlewares.JWTAuth(), channelHandler.UpdateChannel)
		api.DELETE("/channels/:channelID", middlewares.JWTAuth(), channelHandler.DeleteChannel)
		api.GET("/workspaces/:workspaceID/channels", middlewares.JWTAuth(), channelHandler.GetChannelsForWorkspace)

		// Channel Member Routes
		api.POST("/channels/:channelID/members", channelMemberHandler.AddMemberToChannel)
		api.DELETE("/channels/:channelID/members/:userID", channelMemberHandler.RemoveMemberFromChannel)
		api.GET("/channels/:channelID/members", channelMemberHandler.GetChannelMembers)

		// Message Routes
		api.POST("/messages", middlewares.JWTAuth(), messageHandler.CreateMessage)
		api.GET("/messages/:messageID", messageHandler.GetMessageByID)
		api.PUT("/messages/:messageID", middlewares.JWTAuth(), messageHandler.UpdateMessage)
		api.DELETE("/messages/:messageID", middlewares.JWTAuth(), messageHandler.DeleteMessage)
		api.GET("/channels/:channelID/messages", messageHandler.GetMessagesInChannel)

		// Attachment Routes
		api.POST("/attachments", attachmentHandler.CreateAttachment)
		api.GET("/attachments/:attachmentID", attachmentHandler.GetAttachmentByID)
		api.GET("/messages/:messageID/attachments", attachmentHandler.GetAttachmentsForMessage)

		// Reaction Routes
		api.POST("/messages/:messageID/reactions", reactionHandler.AddReaction)
		api.DELETE("/messages/:messageID/reactions/:userID/:emoji", reactionHandler.RemoveReaction)
		api.GET("/messages/:messageID/reactions", reactionHandler.GetReactionsForMessage)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
