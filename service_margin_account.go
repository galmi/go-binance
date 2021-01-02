package binance

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"strconv"
)

func (as *apiService) NewMarginOrder(or NewMarginOrderRequest) (*ProcessedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = or.Symbol
	params["side"] = string(or.Side)
	params["type"] = string(or.Type)
	params["quantity"] = strconv.FormatFloat(or.Quantity, 'f', -1, 64)
	if or.Price > 0.0 {
		params["price"] = strconv.FormatFloat(or.Price, 'f', -1, 64)
	}
	if or.StopPrice != 0 {
		params["stopPrice"] = strconv.FormatFloat(or.StopPrice, 'f', -1, 64)
	}
	if or.NewClientOrderID != "" {
		params["newClientOrderId"] = or.NewClientOrderID
	}
	if or.IcebergQty != 0 {
		params["icebergQty"] = strconv.FormatFloat(or.IcebergQty, 'f', -1, 64)
	}
	if or.SideEffectType != "" {
		params["sideEffectType"] = string(or.SideEffectType)
	}
	if or.TimeInForce != "" {
		params["timeInForce"] = string(or.TimeInForce)
	}
	if or.NewOrderRespType != "" {
		params["newOrderRespType"] = string(or.NewOrderRespType)
	}
	if or.IsIsolated {
		params["isIsolated"] = "TRUE"
	}
	params["timestamp"] = strconv.FormatInt(unixMillis(or.Timestamp), 10)

	res, err := as.request("POST", "sapi/v1/margin/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrder := struct {
		Symbol             string      `json:"symbol"`
		OrderID            int64       `json:"orderId"`
		ClientOrderID      string      `json:"clientOrderId"`
		TransactTime       float64     `json:"transactTime"`
		Price              json.Number `json:"price"`
		OrigQty            json.Number `json:"origQty"`
		ExecutedQty        json.Number `json:"executedQty"`
		CumulativeQuoteQty json.Number `json:"cummulativeQuoteQty"`
		Status             OrderStatus `json:"status"`
		TimeInForce        TimeInForce `json:"timeInForce"`
		Type               OrderType   `json:"type"`
		Side               OrderSide   `json:"side"`
		IsIsolated         bool        `json:"isIsolated"`
	}{}
	if err := json.Unmarshal(textRes, &rawOrder); err != nil {
		return nil, errors.Wrap(err, "rawOrder unmarshal failed")
	}

	t, err := timeFromUnixTimestampFloat(rawOrder.TransactTime)
	if err != nil {
		return nil, err
	}

	price, _ := rawOrder.Price.Float64()
	origQty, _ := rawOrder.OrigQty.Float64()
	executedQty, _ := rawOrder.ExecutedQty.Float64()
	cumulativeQuoteQty, _ := rawOrder.CumulativeQuoteQty.Float64()
	return &ProcessedOrder{
		Symbol:             rawOrder.Symbol,
		OrderID:            rawOrder.OrderID,
		ClientOrderID:      rawOrder.ClientOrderID,
		TransactTime:       t,
		Price:              price,
		OrigQty:            origQty,
		ExecutedQty:        executedQty,
		CumulativeQuoteQty: cumulativeQuoteQty,
		Status:             rawOrder.Status,
		TimeInForce:        rawOrder.TimeInForce,
		Type:               rawOrder.Type,
		Side:               rawOrder.Side,
		IsIsolated:         rawOrder.IsIsolated,
	}, nil
}

func (as *apiService) NewMarginOrderTest(or NewMarginOrderRequest) error {
	params := make(map[string]string)
	params["symbol"] = or.Symbol
	params["side"] = string(or.Side)
	params["type"] = string(or.Type)
	params["quantity"] = strconv.FormatFloat(or.Quantity, 'f', -1, 64)
	if or.Price > 0.0 {
		params["price"] = strconv.FormatFloat(or.Price, 'f', -1, 64)
	}
	if or.StopPrice != 0 {
		params["stopPrice"] = strconv.FormatFloat(or.StopPrice, 'f', -1, 64)
	}
	if or.NewClientOrderID != "" {
		params["newClientOrderId"] = or.NewClientOrderID
	}
	if or.IcebergQty != 0 {
		params["icebergQty"] = strconv.FormatFloat(or.IcebergQty, 'f', -1, 64)
	}
	if or.SideEffectType != "" {
		params["sideEffectType"] = string(or.SideEffectType)
	}
	if or.TimeInForce != "" {
		params["timeInForce"] = string(or.TimeInForce)
	}
	if or.IsIsolated {
		params["isIsolated"] = "TRUE"
	}
	params["timestamp"] = strconv.FormatInt(unixMillis(or.Timestamp), 10)

	res, err := as.request("POST", "sapi/v1/margin/order/test", params, true, true)
	if err != nil {
		return err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return as.handleError(textRes)
	}
	return nil
}

func (as *apiService) QueryMarginOrder(qor QueryOrderRequest) (*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = qor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(qor.Timestamp), 10)
	if qor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(qor.OrderID, 10)
	}
	if qor.OrigClientOrderID != "" {
		params["origClientOrderId"] = qor.OrigClientOrderID
	}
	if qor.IsIsolated {
		params["isIsolated"] = "TRUE"
	}
	if qor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(qor.RecvWindow), 10)
	}

	res, err := as.request("GET", "sapi/v1/margin/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from order.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrder := &rawExecutedOrder{}
	if err := json.Unmarshal(textRes, rawOrder); err != nil {
		return nil, errors.Wrap(err, "rawOrder unmarshal failed")
	}

	eo, err := executedOrderFromRaw(rawOrder)
	if err != nil {
		return nil, err
	}
	return eo, nil
}

func (as *apiService) CancelMarginOrder(cor CancelOrderRequest) (*CanceledOrder, error) {
	params := make(map[string]string)
	params["symbol"] = cor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(cor.Timestamp), 10)
	if cor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(cor.OrderID, 10)
	}
	if cor.OrigClientOrderID != "" {
		params["origClientOrderId"] = cor.OrigClientOrderID
	}
	if cor.NewClientOrderID != "" {
		params["newClientOrderId"] = cor.NewClientOrderID
	}
	if cor.IsIsolated {
		params["isIsolated"] = "TRUE"
	}
	if cor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(cor.RecvWindow), 10)
	}

	res, err := as.request("DELETE", "sapi/v1/margin/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from order.delete")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawCanceledOrder := struct {
		Symbol            string `json:"symbol"`
		OrigClientOrderID string `json:"origClientOrderId"`
		OrderID           int64  `json:"orderId"`
		ClientOrderID     string `json:"clientOrderId"`
	}{}
	if err := json.Unmarshal(textRes, &rawCanceledOrder); err != nil {
		return nil, errors.Wrap(err, "cancelOrder unmarshal failed")
	}

	return &CanceledOrder{
		Symbol:            rawCanceledOrder.Symbol,
		OrigClientOrderID: rawCanceledOrder.OrigClientOrderID,
		OrderID:           rawCanceledOrder.OrderID,
		ClientOrderID:     rawCanceledOrder.ClientOrderID,
	}, nil
}

func (as *apiService) OpenMarginOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = oor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(oor.Timestamp), 10)
	if oor.IsIsolated {
		params["isIsolated"] = "TRUE"
	}
	if oor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(oor.RecvWindow), 10)
	}

	res, err := as.request("GET", "sapi/v1/margin/openOrders", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from openOrders.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrders := []*rawExecutedOrder{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "openOrders unmarshal failed")
	}

	var eoc []*ExecutedOrder
	for _, rawOrder := range rawOrders {
		eo, err := executedOrderFromRaw(rawOrder)
		if err != nil {
			return nil, err
		}
		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) AllMarginOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = aor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(aor.Timestamp), 10)
	if aor.IsIsolated {
		params["isIsolated"] = "TRUE"
	}
	if aor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(aor.OrderID, 10)
	}
	if aor.Limit != 0 {
		params["limit"] = strconv.Itoa(aor.Limit)
	}
	if aor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(aor.RecvWindow), 10)
	}

	res, err := as.request("GET", "sapi/v1/margin/allOrders", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from allOrders.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrders := []*rawExecutedOrder{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "allOrders unmarshal failed")
	}

	var eoc []*ExecutedOrder
	for _, rawOrder := range rawOrders {
		eo, err := executedOrderFromRaw(rawOrder)
		if err != nil {
			return nil, err
		}
		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) MarginAccount(ar AccountRequest) (*MarginAccount, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(ar.Timestamp.Unix()*1000, 10)
	if ar.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(ar.RecvWindow), 10)
	}
	endpoint := "sapi/v1/margin/account"
	if ar.IsIsolated {
		endpoint = "sapi/v1/margin/isolated/account"
	}
	res, err := as.request("GET", endpoint, params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from account.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawAccount := struct {
		BorrowEnabled       bool        `json:"borrowEnabled"`
		MarginLevel         json.Number `json:"marginLevel"`
		TotalAssetOfBtc     json.Number `json:"totalAssetOfBtc"`
		TotalLiabilityOfBtc json.Number `json:"totalLiabilityOfBtc"`
		TotalNetAssetOfBtc  json.Number `json:"totalNetAssetOfBtc"`
		TradeEnabled        bool        `json:"tradeEnabled"`
		TransferEnabled     bool        `json:"transferEnabled"`
		UserAssets          []struct {
			Asset    string      `json:"asset"`
			Borrowed json.Number `json:"borrowed"`
			Free     json.Number `json:"free"`
			Interest json.Number `json:"interest"`
			Locked   json.Number `json:"locked"`
			NetAsset json.Number `json:"netAsset"`
		} `json:"userAssets"`
		Assets []struct {
			Symbol struct {
				Asset    string      `json:"asset"`
				Borrowed json.Number `json:"borrowed"`
				Free     json.Number `json:"free"`
				Interest json.Number `json:"interest"`
				Locked   json.Number `json:"locked"`
				NetAsset json.Number `json:"netAsset"`
			} `json:"baseAsset"`
			BaseSymbol struct {
				Asset    string      `json:"asset"`
				Borrowed json.Number `json:"borrowed"`
				Free     json.Number `json:"free"`
				Interest json.Number `json:"interest"`
				Locked   json.Number `json:"locked"`
				NetAsset json.Number `json:"netAsset"`
			} `json:"quoteAsset"`
			MarginLevel json.Number `json:"marginLevel"`
		} `json:"assets"`
	}{}
	if err := json.Unmarshal(textRes, &rawAccount); err != nil {
		return nil, errors.Wrap(err, "rawAccount unmarshal failed")
	}

	marginLevel, _ := rawAccount.MarginLevel.Float64()
	totalAssetOfBtc, _ := rawAccount.TotalAssetOfBtc.Float64()
	totalLiabilityOfBtc, _ := rawAccount.TotalLiabilityOfBtc.Float64()
	totalNetAssetOfBtc, _ := rawAccount.TotalNetAssetOfBtc.Float64()
	acc := &MarginAccount{
		BorrowEnabled:       rawAccount.BorrowEnabled,
		MarginLevel:         marginLevel,
		TotalAssetOfBtc:     totalAssetOfBtc,
		TotalLiabilityOfBtc: totalLiabilityOfBtc,
		TotalNetAssetOfBtc:  totalNetAssetOfBtc,
		TradeEnabled:        rawAccount.TradeEnabled,
		TransferEnabled:     rawAccount.TransferEnabled,
	}
	for _, b := range rawAccount.UserAssets {
		borrowed, _ := b.Borrowed.Float64()
		free, _ := b.Free.Float64()
		interest, _ := b.Interest.Float64()
		locked, _ := b.Locked.Float64()
		netAsset, _ := b.NetAsset.Float64()
		acc.Assets = append(acc.Assets, &Asset{
			Asset:    b.Asset,
			Borrowed: borrowed,
			Free:     free,
			Interest: interest,
			Locked:   locked,
			NetAsset: netAsset,
		})
	}
	for _, b := range rawAccount.Assets {
		borrowed, _ := b.Symbol.Borrowed.Float64()
		free, _ := b.Symbol.Free.Float64()
		interest, _ := b.Symbol.Interest.Float64()
		locked, _ := b.Symbol.Locked.Float64()
		netAsset, _ := b.Symbol.NetAsset.Float64()
		acc.Assets = append(acc.Assets, &Asset{
			Asset:    b.Symbol.Asset,
			Borrowed: borrowed,
			Free:     free,
			Interest: interest,
			Locked:   locked,
			NetAsset: netAsset,
		})

		borrowed, _ = b.BaseSymbol.Borrowed.Float64()
		free, _ = b.BaseSymbol.Free.Float64()
		interest, _ = b.BaseSymbol.Interest.Float64()
		locked, _ = b.BaseSymbol.Locked.Float64()
		netAsset, _ = b.BaseSymbol.NetAsset.Float64()
		acc.Assets = append(acc.Assets, &Asset{
			Asset:    b.BaseSymbol.Asset,
			Borrowed: borrowed,
			Free:     free,
			Interest: interest,
			Locked:   locked,
			NetAsset: netAsset,
		})
	}

	return acc, nil
}

func (as *apiService) MyMarginTrades(mtr MyTradesRequest) ([]*Trade, error) {
	params := make(map[string]string)
	params["symbol"] = mtr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(mtr.Timestamp), 10)
	if mtr.IsIsolated {
		params["isIsolated"] = "TRUE"
	}
	if mtr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(mtr.RecvWindow), 10)
	}
	if mtr.FromID != 0 {
		params["orderId"] = strconv.FormatInt(mtr.FromID, 10)
	}
	if mtr.Limit != 0 {
		params["limit"] = strconv.Itoa(mtr.Limit)
	}

	res, err := as.request("GET", "sapi/v1/margin/myTrades", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from myTrades.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	var rawTrades []struct {
		ID              int64   `json:"id"`
		Price           string  `json:"price"`
		Qty             string  `json:"qty"`
		Commission      string  `json:"commission"`
		CommissionAsset string  `json:"commissionAsset"`
		Time            float64 `json:"time"`
		IsBuyer         bool    `json:"isBuyer"`
		IsMaker         bool    `json:"isMaker"`
		IsBestMatch     bool    `json:"isBestMatch"`
		IsIsolated      bool    `json:"isIsolated"`
	}
	if err := json.Unmarshal(textRes, &rawTrades); err != nil {
		return nil, errors.Wrap(err, "rawTrades unmarshal failed")
	}

	var tc []*Trade
	for _, rt := range rawTrades {
		price, err := floatFromString(rt.Price)
		if err != nil {
			return nil, err
		}
		qty, err := floatFromString(rt.Qty)
		if err != nil {
			return nil, err
		}
		commission, err := floatFromString(rt.Commission)
		if err != nil {
			return nil, err
		}
		t, err := timeFromUnixTimestampFloat(rt.Time)
		if err != nil {
			return nil, err
		}
		tc = append(tc, &Trade{
			ID:              rt.ID,
			Price:           price,
			Qty:             qty,
			Commission:      commission,
			CommissionAsset: rt.CommissionAsset,
			Time:            t,
			IsBuyer:         rt.IsBuyer,
			IsMaker:         rt.IsMaker,
			IsBestMatch:     rt.IsBestMatch,
			IsIsolated:      rt.IsIsolated,
		})
	}
	return tc, nil
}

func (as *apiService) AllMarginAssets(ar AccountRequest) ([]*MarginAsset, error) {
	assets := []*MarginAsset{}
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(unixMillis(ar.Timestamp), 10)
	if ar.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(ar.RecvWindow), 10)
	}

	res, err := as.request("GET", "sapi/v1/margin/allAssets", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from account.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	var rawAllAssets []struct {
		AssetName     string      `json:"assetName"`
		CanBorrow     bool        `json:"isBorrowable"`
		CanMortage    bool        `json:"isMortgageable"`
		UserMinBorrow json.Number `json:"userMinBorrow"`
		UserMinRepay  json.Number `json:"userMinRepay"`
	}
	if err := json.Unmarshal(textRes, &rawAllAssets); err != nil {
		return nil, errors.Wrap(err, "rawAccount unmarshal failed")
	}

	for _, b := range rawAllAssets {
		minBorrow, _ := b.UserMinBorrow.Float64()
		minRepay, _ := b.UserMinRepay.Float64()
		assets = append(assets, &MarginAsset{
			Asset:         b.AssetName,
			CanBorrow:     b.CanBorrow,
			CanMortage:    b.CanMortage,
			UserMinBorrow: minBorrow,
			UserMinRepay:  minRepay,
		})
	}

	return assets, nil

}

func (as *apiService) MaxBorrow(mbr MaxMarginRequest) (float64, error) {
	params := make(map[string]string)
	params["asset"] = mbr.Symbol
	if mbr.IsIsolated {
		params["isolatedSymbol"] = mbr.PairId
	}
	params["timestamp"] = strconv.FormatInt(unixMillis(mbr.Timestamp), 10)
	if mbr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(mbr.RecvWindow), 10)
	}

	res, err := as.request("GET", "sapi/v1/margin/maxBorrowable", params, true, true)
	if err != nil {
		return 0, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, errors.Wrap(err, "unable to read response from account.maxBorrowable")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, as.handleError(textRes)
	}

	var rawResult struct {
		Amount json.Number `json:"amount"`
	}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return 0, errors.Wrap(err, "rawResult unmarshal failed")
	}
	amount, _ := rawResult.Amount.Float64()

	return amount, nil
}

func (as *apiService) MaxTransfer(mbr MaxMarginRequest) (float64, error) {
	params := make(map[string]string)
	params["asset"] = mbr.Symbol
	if mbr.IsIsolated {
		params["isolatedSymbol"] = mbr.PairId
	}
	params["timestamp"] = strconv.FormatInt(unixMillis(mbr.Timestamp), 10)
	if mbr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(mbr.RecvWindow), 10)
	}

	res, err := as.request("GET", "sapi/v1/margin/maxTransferable", params, true, true)
	if err != nil {
		return 0, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, errors.Wrap(err, "unable to read response from account.maxBorrowable")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, as.handleError(textRes)
	}

	var rawResult struct {
		Amount json.Number `json:"amount"`
	}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return 0, errors.Wrap(err, "rawResult unmarshal failed")
	}
	amount, _ := rawResult.Amount.Float64()

	return amount, nil
}
