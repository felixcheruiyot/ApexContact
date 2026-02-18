package api

import (
	"github.com/livestreamify/backend/internal/api/handlers"
	"github.com/livestreamify/backend/internal/api/middleware"
	"github.com/livestreamify/backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// NewApp wires together the Fiber app, middleware, and all route handlers.
func NewApp(cfg *config.Config) (*fiber.App, func(), error) {
	db, rdb, cleanup, err := setupInfrastructure(cfg)
	if err != nil {
		return nil, nil, err
	}

	app := fiber.New(fiber.Config{
		AppName:      "Live Streamify API v1",
		ErrorHandler: errorHandler,
	})

	registerMiddleware(app, cfg)
	registerRoutes(app, cfg, db, rdb)

	return app, cleanup, nil
}

func registerMiddleware(app *fiber.App, cfg *config.Config) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
}

func registerRoutes(app *fiber.App, cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) {
	authHandler := handlers.NewAuthHandler(cfg, db, rdb)
	eventHandler := handlers.NewEventHandler(cfg, db)
	streamHandler := handlers.NewStreamHandler(cfg, db, rdb)
	paymentHandler := handlers.NewPaymentHandler(cfg, db, rdb)
	adminHandler := handlers.NewAdminHandler(cfg, db, rdb)
	promoterHandler := handlers.NewPromoterHandler(cfg, db)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	v1 := app.Group("/api/v1")

	// ── Auth (public) ──────────────────────────────────────────────────────────
	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", middleware.RequireAuth(cfg), authHandler.Logout)

	// ── Events (public read, protected write) ──────────────────────────────────
	events := v1.Group("/events")
	events.Get("/", eventHandler.List)
	events.Get("/:id", eventHandler.Get)
	events.Post("/", middleware.RequireAuth(cfg), middleware.RequireRole("promoter", "admin"), eventHandler.Create)
	events.Put("/:id", middleware.RequireAuth(cfg), middleware.RequireRole("promoter", "admin"), eventHandler.Update)

	// ── Streaming ──────────────────────────────────────────────────────────────
	// Ingest callback first — no JWT (nginx-rtmp can't send auth headers).
	stream := v1.Group("/stream")
	stream.Post("/ingest/:key/callback", middleware.RequireMediaServerKey(cfg), streamHandler.IngestCallback)
	stream.Post("/:eventId/subscribe", middleware.RequireAuth(cfg), streamHandler.Subscribe)
	stream.Get("/:eventId/token", middleware.RequireAuth(cfg), middleware.AntiPiracy(rdb), streamHandler.GetToken)

	// ── Payments ───────────────────────────────────────────────────────────────
	payments := v1.Group("/payments")
	payments.Post("/initiate", middleware.RequireAuth(cfg), paymentHandler.Initiate)
	payments.Post("/callback", paymentHandler.Callback) // IntaSend webhook (no auth)
	payments.Get("/status/:id", middleware.RequireAuth(cfg), paymentHandler.Status)

	// ── Promoter dashboard ─────────────────────────────────────────────────────
	promoter := v1.Group("/promoter", middleware.RequireAuth(cfg), middleware.RequireRole("promoter", "admin"))
	promoter.Get("/events", promoterHandler.MyEvents)
	promoter.Get("/stream-key/:eventId", promoterHandler.StreamKey)
	promoter.Get("/analytics/:eventId", promoterHandler.Analytics)
	promoter.Get("/revenue", promoterHandler.Revenue)

	// ── Admin ──────────────────────────────────────────────────────────────────
	admin := v1.Group("/admin", middleware.RequireAuth(cfg), middleware.RequireRole("admin"))
	admin.Get("/users", adminHandler.ListUsers)
	admin.Get("/fraud", adminHandler.ListFraudFlags)
	admin.Post("/users/:id/lock", adminHandler.LockUser)
	admin.Post("/users/:id/unlock", adminHandler.UnlockUser)
	admin.Get("/analytics", adminHandler.PlatformAnalytics)
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{"error": err.Error()})
}
