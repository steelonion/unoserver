package internal

// CardColor 表示卡牌的颜色
type CardColor int

const (
	Red CardColor = iota
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

// 新增卡牌
func (cs *UnoCardSet) AddCard(uc UnoCard) bool {
	_, ok := cs.Cards[uc.CardIndex]
	if ok {
		cs.Cards[uc.CardIndex] = uc
	}
	return ok
}

// 校验卡牌合法性
func (cs *UnoCardSet) CheckCard(card UnoCard) bool {
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
	Players        map[int]UnoPlayer // 参与游戏的玩家
	CurrentPlayers int               // 当前准备出牌的玩家
	DiscardPile    UnoCardSet        // 弃牌堆
	RemineCard     UnoCardSet        // 剩余的卡牌堆
	LastCard       *UnoCard          // 上一张出牌的卡片
	TotalAddCount  int               // 总共需要累积的牌数
	IsForword      bool              // 是否为正向顺序
}

// 初始化游戏 移除所有玩家 重置牌堆
func (g *UnoGame) Init() {

}

// 添加玩家
func (g *UnoGame) AddPlayer() {

}

// 切换到下一个玩家
func (g *UnoGame) NextPlayer() {

}

// 发牌
func (g *UnoGame) DealCard(UnoPlayer) {

}

// 检查牌是否符合出牌规则
func (g *UnoGame) CheckCardLegel(uc UnoCard) bool {
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
func (g *UnoGame) Effect(uc UnoCard) {
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

// 出牌
func (g *UnoGame) PlayCard(player int, uc UnoCard) {
	p := g.Players[player]
	if p.CardSet.CheckCard(uc) {
		// 检查出牌是否符合规则
		if g.CheckCardLegel(uc) {
			card, ok := p.CardSet.RemoveCard(uc.CardIndex)
			if ok {
				g.LastCard = &card
				//将牌放入弃牌堆中
				g.DiscardPile.AddCard(card)
			}
		}
	}
}
