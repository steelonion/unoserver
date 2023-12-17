package internal

// 卡牌结构体
type UnoCard struct {
	CardType string
	CardName string
}

// 卡牌堆
type UnoCardSet struct {
	Cards []UnoCard
}

type UnoPlayer struct {
	CardSet UnoCardSet
	Name    string
}

// UnoGame 表示一个 UNO 游戏。
type UnoGame struct {
	Players       []UnoPlayer // 参与游戏的玩家
	RemineCard    UnoCardSet  // 剩余的卡牌堆
	LastCard      UnoCard     // 上一张出牌的卡片
	TotalAddCount int         // 总共需要累积的牌数
	IsForword     bool        // 是否为正向顺序
}
