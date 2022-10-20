package server

import (
	"bytes"
	"io/ioutil"
	"strings"

	//	"bytes"
	"github.com/sirupsen/logrus"
	"net/http"
	///	"net/http"
)

func (srv *Server) ChannelHub() {
	for {
		channel, data, err := srv.RedisProvider.Get()
		if err == nil && data != "" {
			var key, body string
			switch channel {
			case "order":
				key, body, err = FetchOrderDetails(data)
			case "details":
				key, body, err = SendOrderDetails(data)
			}
			if err == nil && body != "" && key != "" {
				err = srv.RedisProvider.Publish(key, body)
				if err != nil {
					logrus.Info("Hub: Unable to publish redis message ", err)
				}
			}
		}
	}
}
func FetchOrderDetails(orders string) (string, string, error) {
	OrderDetails := strings.Split(orders, "|")
	logrus.Info("FetchOrderDetails: Order Received: ", OrderDetails[1])
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/order/", bytes.NewBuffer([]byte(OrderDetails[1])))
	if err != nil {
		logrus.Info("FetchOrderDetails: Unable to make order request")
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Info("FetchOrderDetails: Unable to make order request")
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("FetchOrderDetails: Unable to read response")
		return "", "", err
	}
	//err = srv.RedisProvider.Publish(OrderDetails[0], body)
	//if err != nil {
	//	logrus.Info("FetchOrderDetails: Unable to publish redis message ", err)
	//}
	return OrderDetails[0], string(body), nil
}
func SendOrderDetails(orders string) (string, string, error) {
	logrus.Info("SendOrderDetails: OrderID Received: ", orders)
	client := &http.Client{}
	link := "http://localhost:8081/order/" + orders
	req, err := http.NewRequest(http.MethodGet, link, bytes.NewBuffer([]byte(orders)))
	if err != nil {
		logrus.Info("SendOrderDetails: Unable to make order request")
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Info("SendOrderDetails: Unable to make order request")
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("SendOrderDetails: Unable to read response")
		return "", "", err
	}
	//err = srv.RedisProvider.Publish(orders, body)
	//if err != nil {
	//	logrus.Info("SendOrderDetails: Unable to publish redis message ", err)
	//}
	return orders, string(body), nil
}
