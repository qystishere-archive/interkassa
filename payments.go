package interkassa

import (
	"fmt"
	"sort"
	"strings"
)

const sciURL = `https://sci.interkassa.com/`

// Payment сформированный платеж, содержит в себе данные для последующей его обработки.
type Payment struct {
	// Данные формы для процессинга платежа.
	Form       Form
	// Параметры платежа.
	Parameters PaymentParameters
}

// PaymentParameters параметры платежа, требуемые для его инициализации.
type PaymentParameters struct {
	// Номер платежа. Сохраняется в биллинге Интеркассы.
	// Позволяет идентифицировать платеж в системе, а так же связать с заказами в вашем биллинге.
	// Проверяется уникальность, если в настройках кассы установлена данная опция. Обязательный параметр.
	ID                string              `json:"ik_pm_no"`
	// Валюта платежа. Обязательный параметр, если к кассе подключено больше чем одна валюта. См. настройки кассы.
	Currency          OptionalStringValue `json:"ik_cur,omitempty"`
	// Сумма платежа. Обязательный параметр.
	Amount            string              `json:"ik_am"`
	// Описание платежа. Обязательный параметр.
	Description       string              `json:"ik_desc"`
	// Срок истечения платежа. Не позволяет клиенту оплатить платеж позже указанного срока.
	// Если же он совершил оплату, то средства зачисляются ему на лицевой счет в системе Интеркасса.
	// Параметр используется если платеж привязан к заказу, который быстро теряет свою актуальность
	// с истечением времени. Например: онлайн бронирование. Опциональный параметр.
	ExpiredAt         OptionalTimeValue   `json:"ik_exp,omitempty"`
	// Время жизни платежа. Указывает в секундах срок истечения платежа после его создания.
	// Не используется, если установлен срок истечения платежа (ik_exp).
	// Опциональный параметр.
	// По умолчанию используется свойство кассы "Время жизни платежа" (Payment Lifetime).
	Lifetime          OptionalInt32Value  `json:"ik_ltm,omitempty"`
	// Параметр предназначен для передачи контактных данных плательщика, например email или телефон.
	// Данные сохраняются в системе вместе c платежом и могут использоваться в отдельных случаях для передачи
	// на платежную систему. Также при включенных настройках доставки уведомлений на сервер мерчанта - этот параметра
	// будет присутствовать в теле нотификации (см. 3.4. Оповещение о платеже)
	PayerContact      OptionalStringValue `json:"ik_cli,omitempty"`
	// Включенные способы оплаты. Позволяет указывать доступные способы оплаты для клиента. Опциональный параметр.
	PayWayOn          OptionalStringValue `json:"ik_pw_on,omitempty"`
	// Отключенные способы оплаты. Позволяет указывать недоступные способы оплаты для клиента. Опциональный параметр.
	PayWayOff         OptionalStringValue `json:"ik_pw_off,omitempty"`
	// Выбранный способ оплаты. Позволяет указать точный способ оплаты для клиента.
	// Параметр работает только с параметром действия (ik_act) установленного в "process" или "payway".
	// см. действие (ik_act). Опциональный параметр.
	PayWayVia         OptionalStringValue `json:"ik_pw_via,omitempty"`
	// Локаль. Позволяет явно указать язык и регион установленные для клиента.
	// Формируется по шаблону: [language[_territory]. По умолчанию определяется автоматически.
	Locale            OptionalStringValue `json:"ik_loc,omitempty"`
	// URL страницы взаимодействия. Опциональный параметр.
	InteractionURL    OptionalStringValue `json:"ik_ia_u,omitempty"`
	// Метод запроса страницы взаимодействия. Опциональный параметр.
	InteractionMethod OptionalStringValue `json:"ik_ia_m,omitempty"`
	// URL страницы проведенного платежа. Опциональный параметр.
	SuccessURL        OptionalStringValue `json:"ik_suc_u,omitempty"`
	// Метод запроса страницы проведенного платежа. Опциональный параметр.
	SuccessMethod     OptionalStringValue `json:"ik_suc_m,omitempty"`
	// URL страницы ожидания проведения платежа. Опциональный параметр.
	PendingURL        OptionalStringValue `json:"ik_pnd_u,omitempty"`
	// Метод запроса страницы ожидания проведения платежа. Опциональный параметр.
	PendingMethod     OptionalStringValue `json:"ik_pnd_m,omitempty"`
	// URL страницы непроведенного платежа. Опциональный параметр.
	FailURL           OptionalStringValue `json:"ik_fal_u,omitempty"`
	// Метод запроса страницы непроведенного платежа. Опциональный параметр.
	FailMethod        OptionalStringValue `json:"ik_fal_m,omitempty"`
	// Действие. Позволяет переопределить начальное состояние процесса оплаты.
	// Опциональный параметр. process — обработать; payways — способы оплаты; payway — платежное направление.
	Action            OptionalStringValue `json:"ik_act,omitempty"`
	// Интерфейс. Позволяет указать формат интерфейса SCI как "web" или "json". По умолчанию "web".
	Interface         OptionalStringValue `json:"ik_int,omitempty"`

	// Префикс дополнительных полей. Позволяет передавать дополнительные поля на SCI, после чего эти
	// параметры включаются в данные уведомления о совершенном платеже на страницу взаимодействия.
	AdditionalFields Fields `json:"-"`
}

// NewPayment создаёт новый платёж, получая в результате требуемые параметры для формы, которая потребуется клиенту
// для оплаты.
//
// Для опциональных параметров используете соответствующие функции, см. OptionalString, OptionalTime, OptionalInt32.
func (ik *Interkassa) NewPayment(paymentParameters PaymentParameters) (*Payment, error) {
	var fields map[string]string
	if err := bind(paymentParameters, &fields); err != nil {
		return nil, err
	}
	var additionalFields Fields
	if err := bind(paymentParameters.AdditionalFields, &additionalFields); err != nil {
		return nil, err
	}

	// ID кассы.
	fields["ik_co_id"] = ik.config.ID

	// Дополнительные параметры.
	for k, v := range additionalFields {
		fields[fmt.Sprintf("ik_x_%s", k)] = v
	}

	var keys []string
	for k, _ := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataSet strings.Builder
	for _, key := range keys {
		dataSet.WriteString(fmt.Sprint(fields[key]))
		dataSet.WriteRune(':')
	}
	dataSet.WriteString(ik.config.SignKey)

	// Сигнатура.
	fields["ik_sign"] = ik.sign(dataSet.String())

	return &Payment{
		Form:       Form{
			Method: "POST",
			Action: sciURL,
			Fields: fields,
		},
		Parameters: paymentParameters,
	}, nil
}
