//nolint:tagliatelle // ТЗ требует несовместимые с линтером вещи
package e2e

import (
	"net/http"
)

type UserInfo struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedCoin `json:"received"`
	Sent     []SentCoin     `json:"sent"`
}

type ReceivedCoin struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoin struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

// Тест проверяет практически весь флоу
//
// - Получаем третий токен для пользователя employee2
// - От юзера employee2 проверяем баланс монет, он должен быть равен 1000.
// - От юзера director отправляем 50 монет юзеру employee
// - От юзера director отправляем ещё 10 монет employee
// - От юзера director отправляем 40 монет сотруднику employee2
// - От юзера director проверяем, что баланс равен 900 и в отправленных есть записи о том, что 60 монет отправлено юзеру employee и 40 монет юзеру employee2.
// - От юзера employee отправляем 100 монет директору
// - От юзера employee отправляем 50 монет сотруднику employee2
// - От юзера employee покупаем товар под названием t-shirt дважды и socks один раз.
// - От юзера employee проверяем, что в инвентаре содержится 2 футболки и 1 пара носков, а в полученных есть 60 монет от директора.
// - От юзера employee2 проверяем, что в полученных есть 50 монет от employee и 40 монет от director.

func (s *E2ETestSuite) TestUserInteractions() {

	var tokenResponse tokenResponse
	var errorResponse errorResponse
	var infoResponse UserInfo

	tokens := make(map[string]string)
	users := map[string]struct {
		username string
		password string
	}{
		"director":  {username: "director", password: "password"},
		"employee":  {username: "employee", password: "password"},
		"employee2": {username: "employee2", password: "password2"},
	}

	for role, user := range users {
		authReq := s.client.R().
			SetBody(AuthRequest{
				Username: user.username,
				Password: user.password,
			}).
			SetResult(&tokenResponse).
			SetError(&errorResponse)

		resp, err := authReq.Post(s.env.Server.URL + authURL)
		s.Require().NoError(err, "Failed to get token for "+role)
		s.Require().Equal(http.StatusOK, resp.StatusCode(), "Failed to get token for "+role)
		tokens[role] = tokenResponse.Token
	}

	req := s.client.R().
		SetHeader("Authorization", "Bearer "+tokens["employee2"]).
		SetResult(&infoResponse).
		SetError(&errorResponse)

	resp, err := req.Get(s.env.Server.URL + "/api/info")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())
	s.Assert().Equal(1000, infoResponse.Coins, "Employee2 should start with 1000 coins")

	type SendCoinRequest struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	sendCoins := func(fromToken, toUser string, amount int) {
		req := s.client.R().
			SetHeader("Authorization", "Bearer "+fromToken).
			SetBody(SendCoinRequest{
				ToUser: toUser,
				Amount: amount,
			}).
			SetError(&errorResponse)

		resp, err := req.Post(s.env.Server.URL + "/api/sendCoin")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode())
	}

	sendCoins(tokens["director"], "employee", 50)
	sendCoins(tokens["director"], "employee", 10)
	sendCoins(tokens["director"], "employee2", 40)

	req = s.client.R().
		SetHeader("Authorization", "Bearer "+tokens["director"]).
		SetResult(&infoResponse).
		SetError(&errorResponse)

	resp, err = req.Get(s.env.Server.URL + "/api/info")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())
	s.Assert().Equal(900, infoResponse.Coins, "Director should have 900 coins")

	var foundEmployee, foundEmployee2 bool
	for _, sent := range infoResponse.CoinHistory.Sent {
		if sent.ToUser == "employee" && sent.Amount == 60 {
			foundEmployee = true
		}
		if sent.ToUser == "employee2" && sent.Amount == 40 {
			foundEmployee2 = true
		}
	}
	s.Assert().True(foundEmployee, "Should find 60 coins sent to employee")
	s.Assert().True(foundEmployee2, "Should find 40 coins sent to employee2")

	sendCoins(tokens["employee"], "director", 100)
	sendCoins(tokens["employee"], "employee2", 50)

	buyItem := func(token, item string) {
		req := s.client.R().
			SetHeader("Authorization", "Bearer "+token).
			SetError(&errorResponse)

		resp, err := req.Post(s.env.Server.URL + "/api/buy/" + item)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode())
	}

	buyItem(tokens["employee"], "t-shirt")
	buyItem(tokens["employee"], "t-shirt")
	buyItem(tokens["employee"], "socks")

	req = s.client.R().
		SetHeader("Authorization", "Bearer "+tokens["employee"]).
		SetResult(&infoResponse).
		SetError(&errorResponse)

	resp, err = req.Get(s.env.Server.URL + "/api/info")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())

	var tShirtCount, socksCount int
	for _, item := range infoResponse.Inventory {
		if item.Type == "t-shirt" {
			tShirtCount = item.Quantity
		}
		if item.Type == "socks" {
			socksCount = item.Quantity
		}
	}
	s.Assert().Equal(2, tShirtCount, "Should have 2 t-shirts")
	s.Assert().Equal(1, socksCount, "Should have 1 socks")

	var receivedFromDirector int
	for _, received := range infoResponse.CoinHistory.Received {
		if received.FromUser == "director" {
			receivedFromDirector = received.Amount
		}
	}
	s.Assert().Equal(60, receivedFromDirector, "Should have received 60 coins from director")

	req = s.client.R().
		SetHeader("Authorization", "Bearer "+tokens["employee2"]).
		SetResult(&infoResponse).
		SetError(&errorResponse)

	resp, err = req.Get(s.env.Server.URL + "/api/info")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())

	var receivedFromEmployee, receivedFromDirectorE2 int
	for _, received := range infoResponse.CoinHistory.Received {
		if received.FromUser == "employee" {
			receivedFromEmployee = received.Amount
		}
		if received.FromUser == "director" {
			receivedFromDirectorE2 = received.Amount
		}
	}
	s.Assert().Equal(50, receivedFromEmployee, "Should have received 50 coins from employee")
	s.Assert().Equal(40, receivedFromDirectorE2, "Should have received 40 coins from director")
}
