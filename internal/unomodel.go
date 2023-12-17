package internal

// 卡牌结构体
type UnoCard struct {
	CardIndex int    //卡牌编号
	CardType  string //卡牌类型
	CardName  string //卡牌名称
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
		ok = (card.CardName == value.CardName && card.CardType == value.CardName)
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
	CurrentPlayers UnoPlayer         // 当前准备出牌的玩家
	DiscardPile    UnoCardSet        // 弃牌堆
	RemineCard     UnoCardSet        // 剩余的卡牌堆
	LastCard       UnoCard           // 上一张出牌的卡片
	TotalAddCount  int               // 总共需要累积的牌数
	IsForword      bool              // 是否为正向顺序
}

// 初始化游戏 移除所有玩家 重制牌堆
func (g *UnoGame) Init() {

}

// 添加玩家
func (g *UnoGame) AddPlayer() {

}

// 发牌
func (g *UnoGame) DealCard() {

}

// 出牌
func (g *UnoGame) PlayCard(player int, uc UnoCard) {
	p := g.Players[player]
	//移除卡牌
	if p.CardSet.CheckCard(uc) {
		card, ok := p.CardSet.RemoveCard(uc.CardIndex)
		if ok {
			g.LastCard = card
			//将牌放入弃牌堆中
			g.DiscardPile.AddCard(card)
		}
	}
}
