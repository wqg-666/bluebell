package controller

import "bluebell/models"

// SwaggerResponse describes the standard API response envelope.
type SwaggerResponse struct {
	Code ResCode     `json:"code" example:"1000"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type SwaggerLoginData struct {
	Token    string `json:"token" example:"eyJhbGciOi..."`
	UserID   string `json:"user_id" example:"14361026703134720"`
	UserName string `json:"user_name" example:"alice"`
}

type SwaggerLoginResponse struct {
	Code ResCode          `json:"code" example:"1000"`
	Msg  string           `json:"msg" example:"success"`
	Data SwaggerLoginData `json:"data"`
}

type SwaggerCommunityListResponse struct {
	Code ResCode            `json:"code" example:"1000"`
	Msg  string             `json:"msg" example:"success"`
	Data []models.Community `json:"data"`
}

type SwaggerCommunityDetailResponse struct {
	Code ResCode                `json:"code" example:"1000"`
	Msg  string                 `json:"msg" example:"success"`
	Data models.CommunityDetail `json:"data"`
}

type SwaggerPostDetailResponse struct {
	Code ResCode              `json:"code" example:"1000"`
	Msg  string               `json:"msg" example:"success"`
	Data models.ApiPostDetail `json:"data"`
}

type SwaggerPostListResponse struct {
	Code ResCode                 `json:"code" example:"1000"`
	Msg  string                  `json:"msg" example:"success"`
	Data []*models.ApiPostDetail `json:"data"`
}
