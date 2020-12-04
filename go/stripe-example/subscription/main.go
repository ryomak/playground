package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go/v72"
	stripeClient "github.com/stripe/stripe-go/v72/client"
)

// ChargeJSON incoming data for Stripe API
type SubscriptionJSON struct {
	Amount       int64  `json:"amount"`
	ReceiptEmail string `json:"receiptEmail"`
}

var sc *stripeClient.API

func init() {
	sc = stripeClient.New(os.Getenv("SK_TEST_KEY"), stripe.NewBackends(http.DefaultClient))
}

func main() {
	// set up server
	r := gin.Default()

	// basic hello world GET route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	// our basic charge API route
	r.GET("/api/subscription", func(c *gin.Context) {
		itr := sc.Subscriptions.List(&stripe.SubscriptionListParams{
			Customer: "cus_IU1XQjIpARWMzL",
		})
		res := []interface{}{}
		for itr.Next() {
			res = append(res, itr.Subscription())
		}
		c.JSON(http.StatusOK, res)
	})

	// our basic charge API route
	r.POST("/api/subscription", func(c *gin.Context) {
		if _, err := sc.Subscriptions.New(&stripe.SubscriptionParams{
			Customer: stripe.String("cus_IU1XQjIpARWMzL"),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price:    stripe.String("price_1HtsewJIhAVOVKY0dxbyAI8D"),
					Quantity: stripe.Int64(1),
				},
			},
			BillingCycleAnchor: stripe.Int64(time.Now().AddDate(0, 0, 5).Unix()),
			ProrationBehavior:  stripe.String("none"),
		}); err != nil {
			c.String(http.StatusBadRequest, "Request failed")
			return
		}

		c.String(http.StatusCreated, "Successfully charged")
	})

	r.POST("/api/subscription/update", func(c *gin.Context) {
		var json struct {
			SubscriptionID string `json:"subscriptionId"`
		}
		c.BindJSON(&json)
		if _, err := sc.Subscriptions.Update(json.SubscriptionID, &stripe.SubscriptionParams{
			BillingCycleAnchor: stripe.Int64(time.Now().AddDate(0, 0, 5).Unix()),
			ProrationBehavior:  stripe.String("none"),
		}); err != nil {
			c.String(http.StatusBadRequest, "Request failed")
			return
		}

		c.String(http.StatusCreated, "Successfully charged")
	})

	r.Run(":8080")
}
