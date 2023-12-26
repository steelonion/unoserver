package main

import (
	"context"
	"fmt"
	"net/http"
	"steelonion/unoserver/internal"
)

// 中间件函数，用于鉴权
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			http.Error(w, "param error", http.StatusBadRequest)
		}
		token := cookie.Value

		// 调用鉴权方法
		uid, err := internal.GetUid(token)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// 将用户 ID 存储到请求的上下文中
		ctx := r.Context()
		ctx = context.WithValue(ctx, "uid", uid)
		r = r.WithContext(ctx)

		// 调用下一个处理程序
		next.ServeHTTP(w, r)
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

func main() {
	// 创建一个路由处理器
	mux := http.NewServeMux()

	// 将中间件应用于所有请求
	mux.Handle("/", AuthMiddleware(http.HandlerFunc(MyHandler)))

	// 启动 HTTP 服务器
	http.ListenAndServe(":8080", mux)
}
