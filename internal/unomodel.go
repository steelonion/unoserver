package internal

import (
	"errors"
	"math/rand"
)

// CardColor 表示卡牌的颜色
type CardColor int

const (
	None CardColor = iota
	Red
	Yellow
	Green
	Blue
)

// CardType 表示卡牌的类型
type CardType int

const (
	Number CardType = iota
	Reverse
	Skip
	DrawTwo
	Wild
	WildDrawFour
)

// 卡牌结构体
type UnoCard struct {
	CardIndex int       //卡牌编号
	Color     CardColor //卡牌颜色
	Type      CardType  //卡牌类型
	Value     int       //卡牌值
}

// 卡牌堆
type UnoCardSet struct {
	Cards map[int]UnoCard
}

func (cs *UnoCardSet) New() *UnoCardSet {
	cs.Cards = map[int]UnoCard{}
	return cs
}

// 新增卡牌
func (cs *UnoCardSet) AddCard(uc *UnoCard) bool {
	_, ok := cs.Cards[uc.CardIndex]
	if !ok {
		cs.Cards[uc.CardIndex] = *uc
	}
	return !ok
}

// 校验卡牌合法性
func (cs *UnoCardSet) CheckCard(card *UnoCard) bool {
	value, ok := cs.Cards[card.CardIndex]
	if ok {
		ok = (card.Value == value.Value && card.Color == value.Color && card.Type == value.Type)
	}
	return ok
}

// 移除卡牌
func (cs *UnoCardSet) RemoveCard(index int) (UnoCard, bool) {
	value, ok := cs.Cards[index]
	if ok {
		delete(cs.Cards, index)
	}
	return value, ok
}

type UnoPlayer struct {
	CardSet UnoCardSet //玩家手牌
	Name    string     //玩家名称
	Index   int        //玩家顺序
	ID      int        //玩家ID
}

// UnoGame 表示一个 UNO 游戏。
type UnoGame struct {
	Gaming        bool              // 游戏是否正在进行中
	Players       []UnoPlayer       // 参与游戏的玩家
	PlayersMap    map[int]UnoPlayer // 玩家uid字典
	CurrentPlayer int               // 当前准备出牌的玩家
	DiscardPile   UnoCardSet        // 弃牌堆
	RemineCard    UnoCardSet        // 剩余的卡牌堆
	LastCard      *UnoCard          // 上一张出牌的卡片
	TotalAddCount int               // 总共需要累积的牌数
	IsForword     bool              // 是否为正向顺序
}

func (s *UnoCardSet) initColor(color CardColor, carduid *int) {
	// 加入数字牌
	for i := 0; i < 10; i++ {
		s.AddCard(&UnoCard{Color: color, Type: CardType(Number), Value: i, CardIndex: *carduid})
		*carduid++
	}
	for i := 1; i < 10; i++ {
		s.AddCard(&UnoCard{Color: color, Type: CardType(Number), Value: i, CardIndex: *carduid})
		*carduid++
	}
	// 加入特殊牌
	for i := 0; i < 2; i++ {
		s.AddCard(&UnoCard{Color: color, Type: CardType(Reverse), Value: i, CardIndex: *carduid})
		*carduid++
		s.AddCard(&UnoCard{Color: color, Type: CardType(Skip), Value: i, CardIndex: *carduid})
		*carduid++
		s.AddCard(&UnoCard{Color: color, Type: CardType(DrawTwo), Value: i, CardIndex: *carduid})
		*carduid++
	}
}

func (s *UnoCardSet) initWild(carduid *int) {
	for i := 0; i < 4; i++ {
		s.AddCard(&UnoCard{Color: CardColor(None), Type: CardType(Wild), Value: i, CardIndex: *carduid})
		*carduid++
		s.AddCard(&UnoCard{Color: CardColor(None), Type: CardType(WildDrawFour), Value: i, CardIndex: *carduid})
		*carduid++
	}
}

// 重置游戏 移除所有玩家 重置牌堆
func (g *UnoGame) Reset() {
	// 清空所有数据
	g.DiscardPile = *(&UnoCardSet{}).New()
	g.RemineCard = *(&UnoCardSet{}).New()
	g.LastCard = nil
	g.IsForword = true
	g.TotalAddCount = 0
	g.CurrentPlayer = 0
	g.Players = []UnoPlayer{}
	g.PlayersMap = map[int]UnoPlayer{}
	g.Gaming = false
	//重置牌堆
	carduid := 0
	g.RemineCard.initColor(CardColor(Red), &carduid)
	g.RemineCard.initColor(CardColor(Yellow), &carduid)
	g.RemineCard.initColor(CardColor(Green), &carduid)
	g.RemineCard.initColor(CardColor(Blue), &carduid)
	g.RemineCard.initWild(&carduid)

}

// 启动游戏
func (g *UnoGame) Start() {
	g.Gaming = true
	// 重新发牌
	for i := 0; i < 7; i++ {
		for _, p := range g.Players {
			g.dealCard(&p)
		}
	}
}

// 添加玩家
func (g *UnoGame) AddPlayer(name string, id int) error {
	if g.Gaming {
		return errors.New("can't join game")
	}
	_, ok := g.PlayersMap[id]
	if ok {
		return errors.New("have joined game")
	}
	p := UnoPlayer{Name: name, ID: id, Index: len(g.Players), CardSet: *(&UnoCardSet{}).New()}
	g.Players = append(g.Players, p)
	g.PlayersMap[id] = p
	return nil
}

// 切换到下一个玩家
func (g *UnoGame) NextPlayer() {
	if g.IsForword {
		g.CurrentPlayer++
	} else {
		g.CurrentPlayer--
	}
	if g.CurrentPlayer > len(g.Players) {
		g.CurrentPlayer = 0
	}
	if g.CurrentPlayer < 0 {
		g.CurrentPlayer = len(g.Players) - 1
	}
}

// 发牌
func (g *UnoGame) dealCard(p *UnoPlayer) {
	var keys []int
	for k := range g.RemineCard.Cards {
		keys = append(keys, k)
	}

	randomIndex := rand.Intn(len(keys))
	selectedKey := keys[randomIndex]

	// 获取相应的卡牌
	selectedCard, ok := g.RemineCard.RemoveCard(selectedKey)
	if ok {
		p.CardSet.AddCard(&selectedCard)
	}
}

// 检查牌是否符合出牌规则
func (g *UnoGame) CheckCardLegal(uc *UnoCard) bool {
	if g.LastCard == nil {
		return true
	}

	// 判断颜色是否匹配
	if uc.Color != g.LastCard.Color && uc.Type != Wild && uc.Type != WildDrawFour {
		return false
	}

	// 判断数字或类型是否匹配
	if uc.Type != Wild && uc.Type != WildDrawFour {
		if uc.Type == g.LastCard.Type || uc.Value == g.LastCard.Value {
			return true
		}
	}

	return false
}

// 结算出牌效果
func (g *UnoGame) Effect(uc *UnoCard) {
	switch uc.Type {
	case CardType(WildDrawFour):
		g.TotalAddCount += 4
	case CardType(Wild):
	case CardType(Number):
		g.TotalAddCount += uc.Value
	case CardType(DrawTwo):
		g.TotalAddCount += 2
	case CardType(Reverse):
		g.IsForword = !g.IsForword
	case CardType(Skip):
		g.NextPlayer()
	}
	g.NextPlayer()
}

// 跳过本轮
func (g *UnoGame) Skip(player int) {
	p := g.Players[player]
	//获取罚牌
	g.dealCard(&p)
	//获取累计罚牌
	for i := 0; i < g.TotalAddCount; i++ {
		g.dealCard(&p)
	}
	//清空累积
	g.TotalAddCount = 0
	g.NextPlayer()
}

// 获取用户手牌
func (g *UnoGame) GetCards(player int) (*UnoCardSet, error) {
	if g.CurrentPlayer != player {
		return nil, errors.New("not your turn")
	}
	p := g.Players[player]
	return &p.CardSet, nil
}

// 出牌
func (g *UnoGame) PlayCard(player int, uc *UnoCard) error {
	if g.CurrentPlayer != player {
		return errors.New("not your turn")
	}
	p := g.Players[player]
	// 检查卡牌是否存在
	if !p.CardSet.CheckCard(uc) {
		return errors.New("card not found")
	}
	// 检查出牌是否符合规则
	if !g.CheckCardLegal(uc) {
		return errors.New("illegal card")
	}
	card, ok := p.CardSet.RemoveCard(uc.CardIndex)
	if ok {
		g.LastCard = &card
		//将牌放入弃牌堆中
		g.DiscardPile.AddCard(&card)
	}
	return nil
}
