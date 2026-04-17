package main

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestGetAdBalance(t *testing.T) {
	ctx := context.Background()
	ads := []string{
		"7602565582319517713",
	}

	accessToken := ""
	bcId := ""
	for _, ad := range ads {
		time.Sleep(time.Second)
		res, err := AdvertiserBalanceGet(ctx, accessToken, bcId, ad)
		if err != nil {
			log.Println(err)
			return
		}

		if len(res.Data.AdvertiserAccountList) > 0 {
			log.Println("===:", ad, ",", res.Data.AdvertiserAccountList[0].BudgetRemaining)
		}
	}
}
