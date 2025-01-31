package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	bookingHandler "gitlab.com/final_project1240930/api_gateway/internal/handlers/booking_handler"
	dashboardHandler "gitlab.com/final_project1240930/api_gateway/internal/handlers/dashboard_handler"
	menuHandler "gitlab.com/final_project1240930/api_gateway/internal/handlers/menu_handler"
	tableHandler "gitlab.com/final_project1240930/api_gateway/internal/handlers/table_handler"
	userHandler "gitlab.com/final_project1240930/api_gateway/internal/handlers/user_handler"
	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	internalMiddleware "gitlab.com/final_project1240930/api_gateway/internal/middleware"
	bookingService "gitlab.com/final_project1240930/api_gateway/internal/services/booking"
	dashboardService "gitlab.com/final_project1240930/api_gateway/internal/services/dashboard"
	menuService "gitlab.com/final_project1240930/api_gateway/internal/services/menu"
	tableService "gitlab.com/final_project1240930/api_gateway/internal/services/table"
	userService "gitlab.com/final_project1240930/api_gateway/internal/services/user_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	var err error

	// Load .env file
	if err := godotenv.Load("/app/.env"); err != nil {
		logs.Fatal("Error loading .env file", zap.Error(err))
		return
	}

	// if err := godotenv.Load("../.env"); err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	e := echo.New()

	// ตั้งค่า CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000", // development
			"http://frontend.com",   // production
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true, // cookies , session
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				logs.Error("Error occurred during request processing", zap.Error(err))
				c.Error(err)
			}
			return nil
		}
	})

	servicePort_1 := os.Getenv("SERVICE_1_PORT")
	if servicePort_1 == "" {
		logs.Fatal("SERVICE_1_PORT is not set in .env file")
		return
	}

	servicePort_2 := os.Getenv("SERVICE_2_PORT")
	if servicePort_2 == "" {
		logs.Fatal("SERVICE_2_PORT is not set in .env file")
		return
	}

	servicePort_3 := os.Getenv("SERVICE_3_PORT")
	if servicePort_2 == "" {
		logs.Fatal("SERVICE_3_PORT is not set in .env file")
		return
	}

	grpcUserAddress := "user-management-service:" + servicePort_1
	grpcBookingAddress := "booking-service:" + servicePort_2
	grpcRestaurantAddress := "restaurant-service:" + servicePort_3
	// gRPC User
	userCC, err := grpc.Dial(grpcUserAddress, grpc.WithInsecure())
	if err != nil {
		logs.Fatal("Failed to connect to gRPC server", zap.Error(err))
		return
	}
	defer userCC.Close()

	// gRPC Booking
	bookingCC, err := grpc.Dial(grpcBookingAddress, grpc.WithInsecure())
	if err != nil {
		logs.Fatal("Failed to connect to gRPC server", zap.Error(err))
		return
	}
	defer bookingCC.Close()

	// gRPC Dashboard
	dashboardCC, err := grpc.Dial(grpcBookingAddress, grpc.WithInsecure())
	if err != nil {
		logs.Fatal("Failed to connect to gRPC server", zap.Error(err))
		return
	}
	defer dashboardCC.Close()

	// gRPC Menu
	menuCC, err := grpc.Dial(grpcRestaurantAddress, grpc.WithInsecure())
	if err != nil {
		logs.Fatal("Failed to connect to gRPC server", zap.Error(err))
		return
	}
	defer menuCC.Close()

	// gRPC Table
	tableCC, err := grpc.Dial(grpcRestaurantAddress, grpc.WithInsecure())
	if err != nil {
		logs.Fatal("Failed to connect to gRPC server", zap.Error(err))
		return
	}
	defer tableCC.Close()

	userServiceClient := userService.NewUserServiceClient(userCC)
	userService := userService.NewUserService(userServiceClient)
	userHandler := userHandler.NewUserHandler(userService)

	bookingServiceClient := bookingService.NewBookingServiceClient(bookingCC)
	bookingService := bookingService.NewBookingService(bookingServiceClient)
	bookingHandler := bookingHandler.NewBookingHandler(bookingService)

	menuServiceClient := menuService.NewMenuServiceClient(menuCC)
	menuService := menuService.NewMenuService(menuServiceClient)
	menuHandler := menuHandler.NewMenuHandler(menuService)

	tableServiceClient := tableService.NewTableServiceClient(tableCC)
	tableService := tableService.NewTableService(tableServiceClient)
	tableHandler := tableHandler.NewTableHandler(tableService)

	dashboardServiceClient := dashboardService.NewDashboardServiceClient(dashboardCC)
	dashboardService := dashboardService.NewDashboardService(dashboardServiceClient)
	dashboardHandler := dashboardHandler.NewDashboardHandler(dashboardService)

	// Routes Dashboard Service
	dashboardGroup := e.Group("/dashboard")
	{
		// Public routes
		dashboardGroup.GET("/daily-summary", dashboardHandler.GetDailySummary)                            // สรุปรายวัน
		dashboardGroup.GET("/monthly-sales", dashboardHandler.GetMonthlySales)                            // ยอดขายรายเดือน
		dashboardGroup.GET("/monthly-bookings-customers", dashboardHandler.GetMonthlyBookingAndCustomers) // การจองและลูกค้ารายเดือน
		dashboardGroup.GET("/best-sellers", dashboardHandler.GetBestSellers)                              // สินค้าขายดี
	}

	// Routes User Service
	userGroup := e.Group("/usermanagement")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)

		// Secured routes
		securedUserGroup := userGroup.Group("")
		{
			securedUserGroup.GET("/users", internalMiddleware.AuthMiddleware("user", "manager", "admin")(userHandler.GetAllUsers))
			securedUserGroup.GET("/user/:id", internalMiddleware.AuthMiddleware("user", "manager", "admin")(userHandler.GetUser))
			securedUserGroup.PUT("/user/:id/role", internalMiddleware.AuthMiddleware("admin", "manager")(userHandler.UpdateUserRole))
			securedUserGroup.PUT("/user/:id/password", internalMiddleware.AuthMiddleware("admin", "manager")(userHandler.UpdateUserPassword))
			securedUserGroup.DELETE("/user/:id", internalMiddleware.AuthMiddleware("manager")(userHandler.DeleteUser))
		}
	}

	// Routes Booking Service
	bookingGroup := e.Group("/booking")
	{
		bookingGroup.GET("", bookingHandler.GetBookings)
		bookingGroup.GET("/:booking_id", bookingHandler.GetBookingById)

		securedBookingGroup := bookingGroup.Group("")
		{
			securedBookingGroup.POST("/create", internalMiddleware.AuthMiddleware("user", "manager", "admin")(bookingHandler.CreateBooking))
			securedBookingGroup.PUT("/edit/:booking_id", internalMiddleware.AuthMiddleware("user", "manager", "admin")(bookingHandler.UpdateBooking))
			securedBookingGroup.DELETE("/delete/:booking_id", internalMiddleware.AuthMiddleware("admin")(bookingHandler.DeleteBooking))
		}
	}

	// Routes Menu Service
	menuGroup := e.Group("/menu")
	{
		menuGroup.GET("/items", menuHandler.GetMenuItems)
		menuGroup.GET("/item/:id", menuHandler.GetMenuItemById)
		menuGroup.GET("/sets", menuHandler.GetMenuSets)
		menuGroup.GET("/set/:id", menuHandler.GetMenuSetById)
		menuGroup.GET("/set-items", menuHandler.GetMenuSetItems)
		menuGroup.GET("/set-item/:id", menuHandler.GetMenuSetItemById)

		securedMenuGroup := menuGroup.Group("")
		{
			securedMenuGroup.POST("/item", internalMiddleware.AuthMiddleware("admin", "manager")(menuHandler.CreateMenuItem))
			securedMenuGroup.PUT("/item/:id", internalMiddleware.AuthMiddleware("admin", "manager")(menuHandler.UpdateMenuItem))
			securedMenuGroup.DELETE("/item/:id", internalMiddleware.AuthMiddleware("manager")(menuHandler.DeleteMenuItem))

			securedMenuGroup.POST("/set", internalMiddleware.AuthMiddleware("admin", "manager")(menuHandler.CreateMenuSet))
			securedMenuGroup.PUT("/set/:id", internalMiddleware.AuthMiddleware("admin", "manager")(menuHandler.UpdateMenuSet))
			securedMenuGroup.DELETE("/set/:id", internalMiddleware.AuthMiddleware("manager")(menuHandler.DeleteMenuSet))

			securedMenuGroup.POST("/set-item", internalMiddleware.AuthMiddleware("admin", "manager")(menuHandler.CreateMenuSetItem))
			securedMenuGroup.PUT("/set-item", internalMiddleware.AuthMiddleware("admin", "manager")(menuHandler.UpdateMenuSetItem))
			securedMenuGroup.DELETE("/set-item/:id", internalMiddleware.AuthMiddleware("manager")(menuHandler.DeleteMenuSetItem))

		}
	}

	// Routes Table Service
	tableGroup := e.Group("/tables")
	{
		// Public routes
		tableGroup.GET("/:number", tableHandler.GetTableByNumTable) // Get table by number
		tableGroup.GET("", tableHandler.GetTables)                  // Get all tables
		tableGroup.GET("/types", tableHandler.ListTableTypes)       // List all table types
		tableGroup.GET("/available/:date", tableHandler.GetAvailableTables)

		// Secured routes
		securedTableGroup := tableGroup.Group("")
		{
			// Table management
			securedTableGroup.POST("/create", internalMiddleware.AuthMiddleware("admin", "manager")(tableHandler.CreateTable)) // Create a new table
			securedTableGroup.PUT("/:id", internalMiddleware.AuthMiddleware("admin", "manager")(tableHandler.UpdateTable))     // Update a table
			securedTableGroup.DELETE("/:id", internalMiddleware.AuthMiddleware("admin")(tableHandler.DeleteTable))             // Delete a table

			// Table type management
			securedTableGroup.PUT("/type", internalMiddleware.AuthMiddleware("admin", "manager")(tableHandler.UpdateTableType)) // Update table type
		}
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		logs.Fatal("APP_PORT is not set in .env file")
		return
	}
	logs.Info("Server running on port " + appPort)

	// Start the server
	if err := e.Start(":" + appPort); err != nil {
		logs.Fatal("Failed to serve", zap.Error(err))
	}
}
