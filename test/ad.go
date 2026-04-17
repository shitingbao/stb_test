package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (

	// 获取广告账户余额和预算
	// https://business-api.tiktok.com/portal/docs?id=1739939106470913
	AdvertiserBalanceGetURL = "https://business-api.tiktok.com/open_api/v1.3/advertiser/balance/get/"
)

type TiktokAdvertiserBalanceGetRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AdvertiserAccountList []struct {
			BudgetMode                 string  `json:"budget_mode"`      // 广告账户预算模式。
			BudgetRemaining            float64 `json:"budget_remaining"` // 剩余预算。
			BudgetCost                 float64 `json:"budget_cost"`      // 广告账户已花费的预算
			AccountBalance             float64 `json:"account_balance"`  // 广告账户总余额，保留两位小数。
			BudgetFrequencyRestriction struct {
				UsedCount      int `json:"used_count"`
				TotalCount     int `json:"total_count"`
				RemainingCount int `json:"remaining_count"`
			} `json:"budget_frequency_restriction"`
		} `json:"advertiser_account_list"`
	} `json:"data"`
}

func sget(url string, params map[string]string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func AdvertiserBalanceGet(ctx context.Context, accessToken, bcId, advertiserId string) (TiktokAdvertiserBalanceGetRes, error) {
	bal := TiktokAdvertiserBalanceGetRes{}
	// 预算前置校验（查询余额、频控）
	// 构造查询参数，与 PHP 版保持一致
	fieldsArr := []string{
		"budget_remaining",
		"budget_frequency_restriction",
		"budget_amount_restriction",
		"min_transferable_amount",
	}
	fieldsJSON, err := json.Marshal(fieldsArr)
	if err != nil {
		return bal, err
	}

	filteringJSON, err := json.Marshal(map[string]string{
		"keyword": advertiserId,
	})
	if err != nil {
		return bal, err
	}

	query := map[string]string{
		"bc_id":     bcId,
		"fields":    string(fieldsJSON),
		"filtering": string(filteringJSON),
	}

	b, err := sget(AdvertiserBalanceGetURL, query, map[string]string{
		"Access-Token": accessToken,
		"Accept":       "application/json",
	})
	if err != nil {
		return bal, err
	}

	if err := json.Unmarshal(b, &bal); err != nil {
		return bal, errors.New("接口异常,请重试")
	}

	if bal.Message != "OK" {
		return bal, errors.New("get ad err msg " + bal.Message)
	}
	return bal, nil
}
