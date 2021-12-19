## Golang https://interkassa.com/ api wrapper (WIP)

```go
package main

import (
	"net/http"
	"time"

	"github.com/qystishere/interkassa"
)

func main() {
	// Создаём кассу.
	kassa := interkassa.New(interkassa.Config{
		ID:            "1",                            // ID кассы
		SignAlgorithm: interkassa.SignAlgorithmSha256, // Алгоритм подписи
		SignKey:       "2",                            // Секретный ключ
		SignTestKey:   "3",                            // Тестовый секретный ключ
	})

	// Создаём платеж
	_, err := kassa.NewPayment(interkassa.PaymentParameters{
		ID:                "1",
		Currency:          interkassa.OptionalString("RUB"),
		Amount:            "100.0",
		Description:       "sample description",
		ExpiredAt:         interkassa.OptionalTime(time.Now().Add(time.Hour * 24)),
		PayerContact:      interkassa.OptionalString("sample@sample.ru"),
		InteractionURL:    interkassa.OptionalString("https://local.kikree.com/merchant"),
		InteractionMethod: interkassa.OptionalString("POST"),
		SuccessURL:        interkassa.OptionalString("http://localhost:5000/donation/success"),
		SuccessMethod:     interkassa.OptionalString("GET"),
		PendingURL:        interkassa.OptionalString("http://localhost:5000/donation/pending"),
		PendingMethod:     interkassa.OptionalString("GET"),
		FailURL:           interkassa.OptionalString("http://localhost:5000/donation/failed"),
		FailMethod:        interkassa.OptionalString("GET"),
		AdditionalFields: interkassa.Fields{
			"userID": "222",
		},
	})
	if err != nil {
		panic(err)
	}

	// Обработчик для приёма уведомлений о статусе платежа.
	http.HandleFunc("/merchant", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()

		// Парсим уведомление.
		notification, err := kassa.ParseNotification(r.Form)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		if notification.InvoiceState != "success" {
			w.WriteHeader(200)
			return
		}

		// Используем.
	})
}
```
