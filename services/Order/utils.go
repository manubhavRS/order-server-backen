package Order

import "time"

const ShippingCharges = 50

func FetchExpectDateOfDelivery() time.Time {
	return time.Now().Add(time.Hour * 120)
}
