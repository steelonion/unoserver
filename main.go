package main

import (
	"context"
	"fmt"
	"net/http"
	"steelonion/unoserver/internal"
)

// 中间件函数，用于鉴权
func AuthMiddleware(w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "param error", http.StatusBadRequest)
		return nil, false
	}
	token := cookie.Value

	// 调用鉴权方法
	uid, err := internal.GetUid(token)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return nil, false
	}

	// 将用户 ID 存储到请求的上下文中
	ctx := r.Context()
	ctx = context.WithValue(ctx, "uid", uid)
	newrequest := r.WithContext(ctx)
	return newrequest, true
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
		if r.Method != method {
			http.Error(w, "Error method", http.StatusBadRequest)
			return
		}
		// 调用鉴权中间件
		ok := false
		r, ok = AuthMiddleware(w, r)
		if ok {
			// 鉴权通过，调用下一个处理器
			next.ServeHTTP(w, r)
			return
		}
		// 在其他情况下，返回错误响应
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func main() {
	gm := &internal.GameManager{}
	gm.Init()
	// 创建一个路由处理器
	mux := http.NewServeMux()

	//控制方法
	mux.Handle("/manage/game/create", RequestHandler(gm.CreateGame, http.MethodPost))
	mux.Handle("/manage/game/init", RequestHandler(gm.InitGame, http.MethodPost))
	mux.Handle("/manage/game/clear", RequestHandler(gm.ClearGames, http.MethodPost))
	mux.Handle("/manage/game/start", RequestHandler(gm.StartGame, http.MethodPost))

	//玩家方法
	mux.Handle("/player/game/join", RequestHandler(gm.JoinGame, http.MethodPost))
	mux.Handle("/player/game/getcards", RequestHandler(gm.GetCards, http.MethodGet))
	mux.Handle("/player/game/getgameinfo", RequestHandler(gm.GetGameInfo, http.MethodGet))
	mux.Handle("/player/game/playcard", RequestHandler(gm.PlayCard, http.MethodPost))

	// 启动 HTTP 服务器
	http.ListenAndServe(":8080", mux)
}
