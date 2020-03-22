package binance

// OrderStatus represents order status enum.
type OrderStatus string

// OrderType represents order type enum.
type OrderType string

// OrderSide represents order side enum.
type OrderSide string

type MarginOrderSideEffect string

type NewOrderRespType string

var (
	StatusNew             = OrderStatus("NEW")
	StatusPartiallyFilled = OrderStatus("PARTIALLY_FILLED")
	StatusFilled          = OrderStatus("FILLED")
	StatusCancelled       = OrderStatus("CANCELED")
	StatusPendingCancel   = OrderStatus("PENDING_CANCEL")
	StatusRejected        = OrderStatus("REJECTED")
	StatusExpired         = OrderStatus("EXPIRED")

	TypeLimit  = OrderType("LIMIT")
	TypeMarket = OrderType("MARKET")

	SideBuy  = OrderSide("BUY")
	SideSell = OrderSide("SELL")

	SideEffectNo        = MarginOrderSideEffect("NO_SIDE_EFFECT")
	SideEffectMarginBuy = MarginOrderSideEffect("MARGIN_BUY")
	SideEffectAutoRepay = MarginOrderSideEffect("AUTO_REPAY")

	OrderRespTypeAck    = NewOrderRespType("ACK")
	OrderRespTypeResult = NewOrderRespType("RESULT")
	OrderRespTypeFull   = NewOrderRespType("FULL")
)
