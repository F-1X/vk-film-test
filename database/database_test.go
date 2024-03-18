package database

// import (
// 	"context"
// 	"log"
// 	"testing"
// 	"time"
// 	"vk/config"
// 	mock_database "vk/database/mock"
// 	"vk/model"

// 	"github.com/golang/mock/gomock"
// )

// func TestGetActor(t *testing.T) {
// 	ctl := gomock.NewController(t)
// 	defer ctl.Finish()

// 	repo := mock_database.NewMockDBRepository(ctl)

// 	ctx := context.Background()

// 	in := &model.Actor{
// 		Name:     "test",
// 		Gender:   "test",
// 		Birthday: time.Time{},
// 	}

// 	config, err := config.Read("config.yml")
// 	if err != nil {
// 		log.Fatalf("config read failed: %v", err)
// 	}

// 	Usecase, err := New(config.Database)
// 	if err != nil {
// 		log.Fatalf("database initialization failed: %v", err)
// 	}

// 	Usecase.GetActor(ctx, in)

// }
