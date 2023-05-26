package main

import (
	"fmt"
	"sort"
	"time"

	"nathejk.dk/internal/data"
	"nathejk.dk/internal/mobilepay"
)

func (app *application) syncMobilePay() error {
	client := mobilepay.NewApiClient("https://api.mobilepay.dk", app.config.mobilepay.reportToken)
	rsp, err := client.Transactions(1000, 1)
	if err != nil {
		return err
	}

	latest := app.models.Payments.Latest()
	payments := []data.Payment{}
	for _, t := range rsp.Transactions {
		if t.Type != mobilepay.TransactionTypePayment || t.Timestamp.Before(latest) {
			continue
		}
		payments = append(payments, data.Payment{
			Timestamp:       t.Timestamp,
			ShopNumber:      t.MyShopNumber,
			Amount:          int(t.Amount),
			Currency:        t.Currency,
			Message:         t.Message,
			UserPhoneNumber: t.UserPhoneNumber,
			UserName:        t.UserName,
		})
	}
	if len(payments) > 0 {
		app.logger.PrintInfo("Found transactions newer than lastest", map[string]string{"latest": latest.Format(time.RFC3339), "count": fmt.Sprintf("%d", len(payments))})
	}
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].Timestamp.Before(payments[j].Timestamp)
	})
	for _, payment := range payments {
		app.models.Payments.Insert(&payment)

	}
	return nil
}
