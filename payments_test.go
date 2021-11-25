package interkassa

import (
	"fmt"
	"testing"
	"time"
)

func TestInterkassa_NewPayment(t *testing.T) {
	payment, err := instance.NewPayment(PaymentParameters{
		ID:                "1",
		Currency:          OptionalString("RUB"),
		Amount:            "100.0",
		Description:       "sample description",
		ExpiredAt:         OptionalTime(time.Now().Add(time.Hour * 24)),
		PayerContact:      OptionalString("sample@sample.ru"),
		InteractionURL:    OptionalString("https://local.kikree.com/merchant"),
		InteractionMethod: OptionalString("POST"),
		SuccessURL:        OptionalString("http://localhost:5000/donation/success"),
		SuccessMethod:     OptionalString("GET"),
		PendingURL:        OptionalString("http://localhost:5000/donation/pending"),
		PendingMethod:     OptionalString("GET"),
		FailURL:           OptionalString("http://localhost:5000/donation/failed"),
		FailMethod:        OptionalString("GET"),
		AdditionalFields:  Fields{
			"userID": "222",
		},
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("%#v\n", payment)
}
