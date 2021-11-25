package interkassa

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

var (
	errBadSign       = errors.New("bad sign")
	errBadCheckoutID = errors.New("bad checkout id")
)

// Notification уведомление о проведении платежа.
type Notification struct {
	CheckoutID  string              `json:"ik_co_id"`
	PaymentID   OptionalStringValue `json:"ik_pm_no"`
	Description OptionalStringValue `json:"ik_desc"`
	PayWayVia   OptionalStringValue `json:"ik_pw_via"`
	Amount      string              `json:"ik_am"`
	Currency    string              `json:"ik_cur"`
	Action      string              `json:"ik_act"`

	InvoiceID          OptionalStringValue `json:"ik_inv_id"`
	CheckoutPurseID    OptionalStringValue `json:"ik_co_prs_id"`
	TransactionID      OptionalStringValue `json:"ik_trn_id"`
	InvoiceCreatedAt   Time                `json:"ik_inv_crt"`
	InvoiceProcessedAt Time                `json:"ik_inv_prc"`
	InvoiceState       string              `json:"ik_inv_st"`
	PaySystemPrince    OptionalStringValue `json:"ik_ps_price"`
	CheckoutRefund     OptionalStringValue `json:"ik_co_rfn"`
	PayerContact       OptionalStringValue `json:"ik_cli"`

	CardMask  OptionalStringValue `json:"ik_p_card_mask"`
	CardToken OptionalStringValue `json:"ik_p_card_token"`

	AdditionalFields Fields `json:"-"`
	Sign             string `json:"-"`

	Test bool `json:"-"`
}

// ParseNotification разбирает уведомление о платеже по данным формы (url.Values).
func (ik *Interkassa) ParseNotification(urlValues url.Values) (*Notification, error) {
	fields := make(map[string]string)
	additionalFields := make(Fields)
	for k, v := range urlValues {
		if len(v) != 1 {
			continue
		}

		if strings.HasPrefix(k, "ik_x_") && len(k) > 5 {
			additionalFields[k[5:]] = v[0]
		}

		fields[k] = v[0]
	}

	sign, ok := fields["ik_sign"]
	if !ok {
		return nil, errBadSign
	}
	delete(fields, "ik_sign")

	if fields["ik_co_id"] != ik.config.ID {
		return nil, errBadCheckoutID
	}

	var keys []string
	for k, _ := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataSet strings.Builder
	for _, key := range keys {
		dataSet.WriteString(fields[key])
		dataSet.WriteRune(':')
	}

	var test bool
	// Если оплачено через тестовый метод платежа.
	if fields["ik_pw_via"] == "test_interkassa_test_xts" {
		test = true
		dataSet.WriteString(ik.config.SignTestKey)
	} else {
		dataSet.WriteString(ik.config.SignKey)
	}

	requiredSign := ik.sign(dataSet.String())
	if sign != requiredSign {
		return nil, fmt.Errorf("%w: %s != %s", errBadSign, sign, requiredSign)
	}

	var notification Notification
	if err := bind(fields, &notification); err != nil {
		return nil, err
	}

	if err := bind(additionalFields, &notification.AdditionalFields); err != nil {
		return nil, err
	}

	notification.Test = test
	notification.Sign = sign
	return &notification, nil
}
