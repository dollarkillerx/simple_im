package api

import (
	"simple_im/internal/models"
	"simple_im/internal/storage"
	"simple_im/internal/ws"
	"simple_im/pkg/common/jwt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestEnv holds test environment components
type TestEnv struct {
	DB         *gorm.DB
	Storage    *storage.Storage
	Hub        *ws.Hub
	JWTManager *jwt.JWTManager
}

// SetupTestEnv creates a test environment with SQLite in-memory database
func SetupTestEnv() (*TestEnv, error) {
	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate tables
	err = db.AutoMigrate(
		&models.User{},
		&models.Friend{},
		&models.Group{},
		&models.GroupMember{},
		&models.Message{},
		&models.File{},
	)
	if err != nil {
		return nil, err
	}

	// Create mock Redis client (nil for tests that don't need Redis)
	var redisClient *redis.Client = nil

	st := storage.NewStorage(redisClient, db)
	hub := ws.NewHub()
	jwtManager := jwt.NewJWTManager("test_secret", 3600)

	go hub.Run()

	return &TestEnv{
		DB:         db,
		Storage:    st,
		Hub:        hub,
		JWTManager: jwtManager,
	}, nil
}

// CreateTestUser creates a test user and returns it
func (env *TestEnv) CreateTestUser(username, password string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Nickname: username,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := env.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// CreateTestFriendship creates a friendship between two users
func (env *TestEnv) CreateTestFriendship(userID, friendID int64, status models.FriendStatus) error {
	friend := &models.Friend{
		UserID:   userID,
		FriendID: friendID,
		Status:   status,
	}
	return env.DB.Create(friend).Error
}

// CreateTestGroup creates a test group
func (env *TestEnv) CreateTestGroup(name string, ownerID int64) (*models.Group, error) {
	group := &models.Group{
		Name:    name,
		OwnerID: ownerID,
	}
	if err := env.DB.Create(group).Error; err != nil {
		return nil, err
	}

	// Add owner as member
	member := &models.GroupMember{
		GroupID: group.ID,
		UserID:  ownerID,
		Role:    models.GroupRoleOwner,
	}
	if err := env.DB.Create(member).Error; err != nil {
		return nil, err
	}

	return group, nil
}
