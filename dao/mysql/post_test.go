//go:build integration

package mysql

import (
	"bluebell/models"
	"bluebell/setting"
	"os"
	"testing"
)

func TestCreatePost(t *testing.T) {
	configPath := os.Getenv("BLUEBELL_TEST_CONFIG")
	if configPath == "" {
		t.Skip("set BLUEBELL_TEST_CONFIG to an isolated test database configuration")
	}
	if err := settings.Init(configPath); err != nil {
		t.Fatal(err)
	}
	if err := Init(settings.Conf.MySQLConfig); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(Close)
	post := models.Post{
		ID:          10,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
