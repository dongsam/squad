package types

import (
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OrderBook []OrderGroup

func (ob OrderBook) Add(order Order) OrderBook {
	var prices []string
	for _, og := range ob {
		prices = append(prices, og.Price.String())
	}
	i := sort.Search(len(prices), func(i int) bool {
		return prices[i] <= order.Price.String()
	})
	newOrderBook := ob
	if i < len(prices) {
		if prices[i] == order.Price.String() {
			switch order.Direction {
			case SwapDirectionXToY:
				newOrderBook[i].XToYOrders = append(newOrderBook[i].XToYOrders, order)
			case SwapDirectionYToX:
				newOrderBook[i].YToXOrders = append(newOrderBook[i].YToXOrders, order)
			}
		} else {
			// Insert a new order group at index i.
			newOrderBook = append(newOrderBook[:i], append([]OrderGroup{NewOrderGroup(order)}, newOrderBook[i:]...)...)
		}
	} else {
		// Append a new order group at the end.
		newOrderBook = append(newOrderBook, NewOrderGroup(order))
	}
	return newOrderBook
}

func (ob OrderBook) HighestPriceXToYOrderGroupIndex(start int) (idx int, found bool) {
	for i := start; i < len(ob); i++ {
		if len(ob[i].XToYOrders) > 0 {
			idx = i
			found = true
			return
		}
	}
	return
}

func (ob OrderBook) LowestPriceYToXOrderGroupIndex(start int) (idx int, found bool) {
	for i := start; i >= 0; i-- {
		if len(ob[i].YToXOrders) > 0 {
			idx = i
			found = true
			return
		}
	}
	return
}

func (ob OrderBook) String() string {
	lines := []string{
		"+-----buy------+----------price-----------+-----sell-----+",
	}
	for _, og := range ob {
		lines = append(lines,
			fmt.Sprintf("| %12s | %24s | %-12s |", og.XToYAmount(), og.Price.String(), og.YToXAmount()))
	}
	lines = append(lines, "+--------------+--------------------------+--------------+")
	return strings.Join(lines, "\n")
}

type OrderGroup struct {
	Price      sdk.Dec
	XToYOrders []Order
	YToXOrders []Order
}

func NewOrderGroup(order Order) OrderGroup {
	g := OrderGroup{Price: order.Price}
	switch order.Direction {
	case SwapDirectionXToY:
		g.XToYOrders = append(g.XToYOrders, order)
	case SwapDirectionYToX:
		g.YToXOrders = append(g.YToXOrders, order)
	}
	return g
}

func (og OrderGroup) XToYAmount() sdk.Int {
	amt := sdk.ZeroInt()
	for _, order := range og.XToYOrders {
		amt = amt.Add(order.Amount)
	}
	return amt
}

func (og OrderGroup) YToXAmount() sdk.Int {
	amt := sdk.ZeroInt()
	for _, order := range og.YToXOrders {
		amt = amt.Add(order.Amount)
	}
	return amt
}

// Order represents a swap order, which is made by a user or a pool.
// TODO: use SwapRequest instead - all fields are identical?
type Order struct {
	Orderer   string
	Direction SwapDirection
	Price     sdk.Dec
	Amount    sdk.Int
}

func NewOrder(orderer string, dir SwapDirection, price sdk.Dec, amt sdk.Int) Order {
	return Order{orderer, dir, price, amt}
}
