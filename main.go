package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"steelonion/unoserver/internal"
)

// 游戏实例表
var games map[int]internal.UnoGame

// 中间件函数，用于鉴权
func AuthMiddleware(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "param error", http.StatusBadRequest)
		return false
	}
	token := cookie.Value

	// 调用鉴权方法
	uid, err := internal.GetUid(token)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return false
	}

	// 将用户 ID 存储到请求的上下文中
	ctx := r.Context()
	ctx = context.WithValue(ctx, "uid", uid)
	r = r.WithContext(ctx)
	return true
}

// 尝试将任意函数转换为 http.Handler 接口
func ConvertToHandler(f interface{}) (http.Handler, error) {
	if handlerFunc, ok := f.(func(http.ResponseWriter, *http.Request)); ok {
		// 如果可以转换为 http.HandlerFunc，则返回它
		return http.HandlerFunc(handlerFunc), nil
	}

	// 如果无法转换，返回错误
	return nil, fmt.Errorf("cannot convert the function to http.Handler")
}

// 请求处理器
func RequestHandler(handler interface{}, method string) http.Handler {
	next, err := ConvertToHandler(handler)
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			// 调用鉴权中间件
			if AuthMiddleware(w, r) {
				// 鉴权通过，调用下一个处理器
				next.ServeHTTP(w, r)
				return
			}
		}
		// 在其他情况下，返回错误响应
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// 处理请求的示例处理程序
func MyHandler(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户 ID
	uid, ok := r.Context().Value("uid").(int)
	if !ok {
		http.Error(w, "无法获取用户 ID", http.StatusInternalServerError)
		return
	}
	// 处理请求，使用用户 ID
	fmt.Fprintf(w, "Hello, User %d!", uid)
}

func CreateGame(w http.ResponseWriter, r *http.Request) {
	for {
		randomNumber := rand.Intn(100)
		_, ok := games[randomNumber]
		if !ok {
			games[randomNumber] = internal.UnoGame{}
			ret := make(map[string]int)
			ret["id"] = randomNumber
			jsonData, err := json.Marshal(ret)
			if err != nil {
				http.Error(w, "marshal error", http.StatusInternalServerError)
				return
			}
			w.Write(jsonData)
			break
		}
	}
}

func ClearGames(w http.ResponseWriter, r *http.Request) {
	games = make(map[int]internal.UnoGame)
}

func InitGame(w http.ResponseWriter, r *http.Request) {
	r.FormValue("gameid")
}

func main() {
	games = make(map[int]internal.UnoGame)
	// 创建一个路由处理器
	mux := http.NewServeMux()

	mux.Handle("/manage/game/create", RequestHandler(CreateGame, http.MethodPost))

	// 启动 HTTP 服务器
	http.ListenAndServe(":8080", mux)
}
