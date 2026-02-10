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
	attachmentRepo := repositories.NewAttachmentRepo(bunDB, s.log)
	channelMemberRepo := repositories.NewChannelMemberRepo(bunDB, s.log)
	channelRepo := repositories.NewChannelRepo(bunDB, s.log)
	messageRepo := repositories.NewMessageRepo(bunDB, s.log)
	meetingRepo := repositories.NewMeetingRepo(bunDB, s.log)
	reactionRepo := repositories.NewReactionRepo(bunDB, s.log)
	userRepo := repositories.NewUserRepo(bunDB, s.log)
	workspaceMemberRepo := repositories.NewWorkspaceMemberRepo(bunDB, s.log)
	workspaceRepo := repositories.NewWorkspaceRepo(bunDB, s.log)

	// --- Services ---
	attachmentService := services.NewAttachmentService(attachmentRepo, s.log)
	channelMemberService := services.NewChannelMemberService(channelMemberRepo, s.log)
	channelService := services.NewChannelService(channelRepo, channelMemberRepo, workspaceMemberRepo, s.log)
	messageService := services.NewMessageService(messageRepo, meetingRepo, s.log)
	reactionService := services.NewReactionService(reactionRepo, s.log)
	userService := services.NewUserService(userRepo, s.log)
	meetingService := services.NewMeetingService(meetingRepo, channelRepo, userRepo, channelMemberRepo, s.log)
	workspaceMemberService := services.NewWorkspaceMemberService(workspaceMemberRepo, workspaceRepo, s.log)
	workspaceService := services.NewWorkspaceService(workspaceRepo, workspaceMemberRepo, s.log)
	meetingChatService := services.NewMeetingChatService(meetingRepo, messageRepo, userRepo, attachmentRepo, reactionRepo, s.log) // Initialize MeetingChatService

	// --- Handlers ---
	attachmentHandler := handlers.NewAttachmentHandler(attachmentService, s.log)
	channelMemberHandler := handlers.NewChannelMemberHandler(channelMemberService, s.log)
	channelHandler := handlers.NewChannelHandler(channelService, s.log)
	messageHandler := handlers.NewMessageHandler(messageService, s.log)
	reactionHandler := handlers.NewReactionHandler(reactionService, s.log)
	userHandler := handlers.NewUserHandler(userService, s.log)
	meetingHandler := handlers.NewMeetingHandler(meetingService, s.log)
	workspaceMemberHandler := handlers.NewWorkspaceMemberHandler(workspaceMemberService, s.log)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceService, s.log)
	chatHandler := handlers.NewChatHandler(meetingChatService, s.log) // Initialize ChatHandler

	// --- API Routes ---
	api := r.Group("/api")
	{
		// User Routes
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/users", middlewares.JWTAuth(s.log), userHandler.GetUserByID)
		api.GET("/users/by-email", userHandler.GetUserByEmail)       // Query param: ?email=
		api.GET("/users/by-username", userHandler.GetUserByUsername) // Query param: ?username=
		api.PUT("/users", middlewares.JWTAuth(s.log), userHandler.UpdateUser)
		api.DELETE("/users", middlewares.JWTAuth(s.log), userHandler.DeleteUser)

		// Workspace Routes
		api.POST("/workspaces", middlewares.JWTAuth(s.log), workspaceHandler.CreateWorkspace)
		api.GET("/workspaces/:workspaceID", workspaceHandler.GetWorkspaceByID)
		api.PUT("/workspaces/:workspaceID", middlewares.JWTAuth(s.log), workspaceHandler.UpdateWorkspace)
		api.DELETE("/workspaces/:workspaceID", middlewares.JWTAuth(s.log), workspaceHandler.DeleteWorkspace)
		api.GET("/workspaces", middlewares.JWTAuth(s.log), workspaceHandler.GetWorkspacesForUser)

		// Workspace Member Routes
		api.POST("/workspaces/:workspaceID/members", workspaceMemberHandler.AddMemberToWorkspace)
		api.DELETE("/workspaces/:workspaceID/members/:userID", workspaceMemberHandler.RemoveMemberFromWorkspace)
		api.GET("/workspaces/:workspaceID/members", workspaceMemberHandler.GetWorkspaceMembers)
		api.POST("/workspaces/:workspaceID/join", middlewares.JWTAuth(s.log), workspaceMemberHandler.JoinWorkspace)

		// Channel Routes
		api.POST("/channels", middlewares.JWTAuth(s.log), channelHandler.CreateChannel)
		api.GET("/channels/:channelID", middlewares.JWTAuth(s.log), channelHandler.GetChannelByID)
		api.PUT("/channels/:channelID", middlewares.JWTAuth(s.log), channelHandler.UpdateChannel)
		api.DELETE("/channels/:channelID", middlewares.JWTAuth(s.log), channelHandler.DeleteChannel)
		api.GET("/workspaces/:workspaceID/channels", middlewares.JWTAuth(s.log), channelHandler.GetChannelsForWorkspace)

		// Channel Member Routes
		api.POST("/channels/:channelID/members", channelMemberHandler.AddMemberToChannel)
		api.DELETE("/channels/:channelID/members/:userID", channelMemberHandler.RemoveMemberFromChannel)
		api.GET("/channels/:channelID/members", channelMemberHandler.GetChannelMembers)

		// Message Routes
		api.POST("/messages", middlewares.JWTAuth(s.log), messageHandler.CreateMessage)
		api.GET("/messages/:messageID", messageHandler.GetMessageByID)
		api.PUT("/messages/:messageID", middlewares.JWTAuth(s.log), messageHandler.UpdateMessage)
		api.DELETE("/messages/:messageID", middlewares.JWTAuth(s.log), messageHandler.DeleteMessage)
		api.GET("/meetings/:meetingID/messages", messageHandler.GetMessagesInMeeting)

		// Meeting Routes
		api.POST("/meetings", middlewares.JWTAuth(s.log), meetingHandler.CreateMeeting)
		api.GET("/meetings/:meetingID", meetingHandler.GetMeetingByID)
		api.PUT("/meetings/:meetingID", middlewares.JWTAuth(s.log), meetingHandler.UpdateMeeting)
		api.DELETE("/meetings/:meetingID", middlewares.JWTAuth(s.log), meetingHandler.DeleteMeeting)
		api.GET("/channels/:channelID/meetings", meetingHandler.GetMeetingsByChannelID)
		api.POST("/meetings/:meetingID/participants", middlewares.JWTAuth(s.log), meetingHandler.AddParticipant)
		api.DELETE("/meetings/:meetingID/participants/:participantID", middlewares.JWTAuth(s.log), meetingHandler.RemoveParticipant)

		// Attachment Routes
		api.POST("/attachments", attachmentHandler.CreateAttachment)
		api.GET("/attachments/:attachmentID", attachmentHandler.GetAttachmentByID)
		api.GET("/messages/:messageID/attachments", attachmentHandler.GetAttachmentsForMessage)

		// Reaction Routes
		api.POST("/messages/:messageID/reactions", reactionHandler.AddReaction)
		api.DELETE("/messages/:messageID/reactions/:userID/:emoji", reactionHandler.RemoveReaction)
		api.GET("/messages/:messageID/reactions", reactionHandler.GetReactionsForMessage)
	}

	// --- WebSocket Routes ---
	wsGroup := r.Group("/ws")
	wsGroup.Use(middlewares.JWTAuth(s.log))
	{
		wsGroup.GET("/meeting/:meeting_id/chat", chatHandler.ServeMeetingChatWs)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	s.log.Info().Msg("Handling HelloWorld request")
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	s.log.Info().Msg("HelloWorld response sent")
	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	s.log.Info().Msg("Handling health check request")
	healthStatus := s.db.Health()
	s.log.Info().Interface("health_status", healthStatus).Msg("Database health status retrieved")
	c.JSON(http.StatusOK, healthStatus)
}
