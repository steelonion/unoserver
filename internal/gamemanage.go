package internal

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
)

type GameManager struct {
	// 游戏实例表
	games map[int]*UnoGame
}

func (gm *GameManager) Init() {
	gm.games = make(map[int]*UnoGame)
}

func (gm *GameManager) CreateGame(w http.ResponseWriter, r *http.Request) {
	for {
		randomNumber := rand.Intn(100)
		_, ok := gm.games[randomNumber]
		if !ok {
			gm.games[randomNumber] = &UnoGame{}
			gm.games[randomNumber].Reset()
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

func (gm *GameManager) ClearGames(w http.ResponseWriter, r *http.Request) {
	gm.games = make(map[int]*UnoGame)
	w.WriteHeader(http.StatusOK)
}

func (gm *GameManager) GetGame(w http.ResponseWriter, r *http.Request) *UnoGame {
	// 通过 r.FormValue 获取表单值
	valueStr := r.URL.Query().Get("gameid")
	// 尝试将字符串转换为整数
	id, err := strconv.Atoi(valueStr)
	if err != nil {
		// 处理转换错误
		http.Error(w, "Invalid gameid", http.StatusBadRequest)
		return nil
	}
	game, ok := gm.games[id]
	if !ok {
		// 处理转换错误
		http.Error(w, "Invalid gameid", http.StatusBadRequest)
		return nil
	}
	return game
}

func (gm *GameManager) InitGame(w http.ResponseWriter, r *http.Request) {

	game := gm.GetGame(w, r)
	if game == nil {
		return
	}
	game.Reset()
	w.WriteHeader(http.StatusOK)
}

func (gm *GameManager) StartGame(w http.ResponseWriter, r *http.Request) {
	game := gm.GetGame(w, r)
	if game == nil {
		return
	}
	game.Start()
	w.WriteHeader(http.StatusOK)
}

func (gm *GameManager) GetCards(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户 ID
	uid, ok := r.Context().Value("uid").(int)
	if !ok {
		http.Error(w, "无法获取用户 ID", http.StatusInternalServerError)
		return
	}
	game := gm.GetGame(w, r)
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

func (gm *GameManager) JoinGame(w http.ResponseWriter, r *http.Request) {
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
	game := gm.GetGame(w, r)
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

func (gm *GameManager) GetGameInfo(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户 ID
	_, ok := r.Context().Value("uid").(int)
	if !ok {
		http.Error(w, "无法获取用户 ID", http.StatusInternalServerError)
		return
	}
	game := gm.GetGame(w, r)
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

func (gm *GameManager) PlayCard(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户 ID
	uid, ok := r.Context().Value("uid").(int)
	if !ok {
		http.Error(w, "无法获取用户 ID", http.StatusInternalServerError)
		return
	}
	game := gm.GetGame(w, r)
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
