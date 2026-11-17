package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

var translatorOnce sync.Once

// sync.Once，保证某个操作在程序运行期间最多执行一次。
var translatorErr error

// 确保测试使用的翻译器已经初始化
func ensureTestTranslator(t *testing.T) {
	t.Helper() //当前函数是一个测试辅助函数。

	translatorOnce.Do(func() {
		translatorErr = InitTrans("zh")
	})
	if translatorErr != nil {
		//Fatalf 只会结束当前测试函数，不会终止整个测试程序。
		t.Fatalf("initialize validator translations: %v", translatorErr)
	}
}

// 模拟发送一次 HTTP 请求，并返回记录的响应
func performRequest(handler gin.HandlerFunc, method, target, body string) *httptest.ResponseRecorder {
	//method：HTTP 方法；
	//target：请求路径；
	//body：请求体内容。
	gin.SetMode(gin.TestMode)
	router := gin.New()
	//例如，如果：
	//method = "POST"
	//target = "/api/v1/post"
	//handler = CreatePostHandler
	//这行代码就相当于：
	//router.POST("/api/v1/post", CreatePostHandler)
	router.Handle(method, target, handler)
	//创建一个用于测试的 HTTP 请求
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	//这个请求头告诉服务器： 当前请求体中的数据是 JSON 格式。
	request.Header.Set("Content-Type", "application/json")
	//创建一个 HTTP 响应记录器。 它不会真的把响应发到网络，而是把处理函数返回的内容记录下来
	response := httptest.NewRecorder()
	//模拟一次 HTTP 请求，并让 Gin 路由器处理这个请求，最后把响应结果写入 response。
	router.ServeHTTP(response, request)
	return response
}

// 把响应 JSON 转换成 ResponseData 结构体
func decodeResponse(t *testing.T, response *httptest.ResponseRecorder) ResponseData {
	t.Helper()
	var data ResponseData
	//response.Body.Bytes()
	//读取响应体中的所有字节数据。
	//json.Unmarshal(response.Body.Bytes(), &data)
	//把响应体中的 JSON 数据解析到 data 变量中。
	if err := json.Unmarshal(response.Body.Bytes(), &data); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, response.Body.String())
	}
	return data
}

// 验证未登录时不能创建帖子
func TestCreatePostHandlerRequiresLogin(t *testing.T) {
	//调用前面定义的 performRequest，模拟一次创建帖子的请求。
	response := performRequest(CreatePostHandler,
		http.MethodPost,
		"/api/v1/post",
		`{"community_id":1,"title":"test","content":"body"}`,
	)
	if response.Code != http.StatusOK {
		t.Fatalf("expected HTTP %d, got %d", http.StatusOK, response.Code)
	}
	//解析响应体，然后检查业务状态码。把响应体 JSON 解析成 ResponseData。
	if got := decodeResponse(t, response).Code; got != CodeNeedLogin {
		t.Fatalf("expected code %d, got %d", CodeNeedLogin, got)
	}
}

// 验证创建帖子接口会拒绝非法 JSON
func TestCreatePostHandlerRejectsInvalidJSON(t *testing.T) {
	response := performRequest(CreatePostHandler,
		http.MethodPost,
		"/api/v1/post",
		`{}`)
	if got := decodeResponse(t, response).Code; got != CodeInvalidParam {
		t.Fatalf("expected code %d, got %d", CodeInvalidParam, got)
	}
}

// 验证注册时两次密码不一致会被拒绝
func TestSignUpHandlerRejectsMismatchedPasswords(t *testing.T) {
	ensureTestTranslator(t)
	response := performRequest(SignUpHandler,
		http.MethodPost,
		"/api/v1/signup",
		`{"username":"www11","password":"one","re_password":"two"}`)
	if got := decodeResponse(t, response).Code; got != CodeInvalidParam {
		t.Fatalf("expected code %d, got %d", CodeInvalidParam, got)
	}
}

// 验证登录时缺少账号密码会被拒绝
func TestLoginHandlerRejectsMissingCredentials(t *testing.T) {
	ensureTestTranslator(t)
	//它模拟的是：
	//用户点击登录，但是没有填写用户名和密码，就直接提交了请求。
	//注意：这是故意构造一个不符合业务要求的 HTTP 请求
	response := performRequest(LoginHandler,
		http.MethodPost,
		"/api/v1/login",
		"{}")
	if got := decodeResponse(t, response).Code; got != CodeInvalidParam {
		t.Fatalf("expected code %d, got %d", CodeInvalidParam, got)
	}
}

// 验证未登录时不能给帖子投票验证未登录时不能给帖子投票
func TestPostVoteControllerRequiresLogin(t *testing.T) {
	response := performRequest(PostVoteController,
		http.MethodPost,
		"/api/v1/vote",
		`{"post_id":"123","direction":"1"}`)
	if got := decodeResponse(t, response).Code; got != CodeNeedLogin {
		t.Fatalf("expected code %d, got %d", CodeNeedLogin, got)
	}
}

// 验证帖子详情接口会拒绝非法帖子 ID
func TestGetPostDetailHandlerRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/post/:id", GetPostDetailHandler)
	//创建一个 GET 测试请求
	request := httptest.NewRequest(http.MethodGet,
		"/api/v1/post/not-an-id",
		nil)
	//创建响应记录器，用于记录接口返回的 HTTP 状态码、响应头和响应
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if got := decodeResponse(t, response).Code; got != CodeInvalidParam {
		t.Fatalf("expected code %d, got %d", CodeInvalidParam, got)
	}
}
