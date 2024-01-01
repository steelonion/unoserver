package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"steelonion/unoserver/internal"
	"strconv"
)

// 游戏实例表
var games map[int]*internal.UnoGame

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
			games[randomNumber] = &internal.UnoGame{}
			games[randomNumber].Reset()
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
	games = make(map[int]*internal.UnoGame)
	w.WriteHeader(http.StatusOK)
}

func GetGame(w http.ResponseWriter, r *http.Request) *internal.UnoGame {
	// 通过 r.FormValue 获取表单值
	valueStr := r.URL.Query().Get("gameid")
	// 尝试将字符串转换为整数
	id, err := strconv.Atoi(valueStr)
	if err != nil {
		// 处理转换错误
		http.Error(w, "Invalid gameid", http.StatusBadRequest)
		return nil
	}
	game, ok := games[id]
	if !ok {
		// 处理转换错误
		http.Error(w, "Invalid gameid", http.StatusBadRequest)
		return nil
	}
	return game
}

func InitGame(w http.ResponseWriter, r *http.Request) {

	game := GetGame(w, r)
	if game == nil {
		return
	}
	game.Reset()
	w.WriteHeader(http.StatusOK)
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	game := GetGame(w, r)
	if game == nil {
		return
	}
	game.Start()
	w.WriteHeader(http.StatusOK)
}

func GetCards(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户 ID
	uid, ok := r.Context().Value("uid").(int)
	if !ok {
		http.Error(w, "无法获取用户 ID", http.StatusInternalServerError)
		return
	}
	game := GetGame(w, r)
	if game == nil {
		return
	}
	set, err := game.GetCards(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret := make(map[string]interface{})
	ret["cards"] = *set
	ret["uid"] = uid
	jsonData, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func JoinGame(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户 ID
	uid, ok := r.Context().Value("uid").(int)
	if !ok {
		http.Error(w, "无法获取用户 ID", http.StatusInternalServerError)
		return
	}
	playername := r.URL.Query().Get("name")
	if playername == "" {
		http.Error(w, "name empty", http.StatusBadRequest)
		return
	}
	game := GetGame(w, r)
	if game == nil {
		return
	}
	err := game.AddPlayer(playername, uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetGameInfo(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户 ID
	_, ok := r.Context().Value("uid").(int)
	if !ok {
		http.Error(w, "无法获取用户 ID", http.StatusInternalServerError)
		return
	}
	game := GetGame(w, r)
	if game == nil {
		return
	}
	ret := make(map[string]interface{})
	ret["isForword"] = game.IsForword
	ret["isGaming"] = game.Gaming
	if game.LastCard == nil {
		ret["lastCard"] = nil
	} else {
		ret["lastCard"] = *game.LastCard
	}
	ret["totalAddCount"] = game.TotalAddCount

	ret["currentPlayerIndex"] = game.CurrentPlayer
	ret["players"] = make(map[int]interface{})
	for i := 0; i < len(game.Players); i++ {
		p := game.Players[i]
		pi := make(map[string]interface{})
		pi["uid"] = p.ID
		pi["name"] = p.Name
		pi["index"] = p.Index
		pi["cardCount"] = len(p.CardSet.Cards)
		ret["players"].(map[int]interface{})[p.ID] = pi
	}
	jsonData, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func PlayCard(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户 ID
	uid, ok := r.Context().Value("uid").(int)
	if !ok {
		http.Error(w, "无法获取用户 ID", http.StatusInternalServerError)
		return
	}
	game := GetGame(w, r)
	if game == nil {
		return
	}
	uciStr := r.URL.Query().Get("uci")
	uci, err := strconv.Atoi(uciStr)
	if err != nil {
		// 处理转换错误
		http.Error(w, "Invalid uci", http.StatusBadRequest)
		return
	}
	err = game.PlayCard(uid, uci)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	games = make(map[int]*internal.UnoGame)
	// 创建一个路由处理器
	mux := http.NewServeMux()

	//控制方法
	mux.Handle("/manage/game/create", RequestHandler(CreateGame, http.MethodPost))
	mux.Handle("/manage/game/init", RequestHandler(InitGame, http.MethodPost))
	mux.Handle("/manage/game/clear", RequestHandler(ClearGames, http.MethodPost))
	mux.Handle("/manage/game/start", RequestHandler(StartGame, http.MethodPost))

	//玩家方法
	mux.Handle("/player/game/join", RequestHandler(JoinGame, http.MethodPost))
	mux.Handle("/player/game/getcards", RequestHandler(GetCards, http.MethodGet))
	mux.Handle("/player/game/getgameinfo", RequestHandler(GetGameInfo, http.MethodGet))
	mux.Handle("/player/game/playcard", RequestHandler(PlayCard, http.MethodPost))

	// 启动 HTTP 服务器
	http.ListenAndServe(":8080", mux)
}
