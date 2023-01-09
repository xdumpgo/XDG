package shoppy

type SellixWebhook struct {
	Event string `json:"event"`
	Data  struct {
		ID                        int         `json:"id"`
		Uniqid                    string      `json:"uniqid"`
		Total                     int         `json:"total"`
		TotalDisplay              int         `json:"total_display"`
		ExchangeRate              int         `json:"exchange_rate"`
		CryptoExchangeRate        float64     `json:"crypto_exchange_rate"`
		Currency                  string      `json:"currency"`
		ShopID                    int         `json:"shop_id"`
		Name                      string      `json:"name"`
		CustomerEmail             string      `json:"customer_email"`
		PaypalEmailDelivery       int         `json:"paypal_email_delivery"`
		ProductID                 string      `json:"product_id"`
		ProductTitle              string      `json:"product_title"`
		ProductType               string      `json:"product_type"`
		Gateway                   string      `json:"gateway"`
		PaypalEmail               interface{} `json:"paypal_email"`
		PaypalOrderID             interface{} `json:"paypal_order_id"`
		PaypalPayerEmail          interface{} `json:"paypal_payer_email"`
		PaypalFee                 int         `json:"paypal_fee"`
		SkrillEmail               interface{} `json:"skrill_email"`
		SkrillSid                 interface{} `json:"skrill_sid"`
		SkrillLink                interface{} `json:"skrill_link"`
		PerfectmoneyID            interface{} `json:"perfectmoney_id"`
		CryptoAddress             string      `json:"crypto_address"`
		CryptoAmount              float64     `json:"crypto_amount"`
		CryptoReceived            float64     `json:"crypto_received"`
		CryptoURI                 string      `json:"crypto_uri"`
		CryptoConfirmationsNeeded int         `json:"crypto_confirmations_needed"`
		Country                   string      `json:"country"`
		Location                  string      `json:"location"`
		IP                        string      `json:"ip"`
		IsVpnOrProxy              bool        `json:"is_vpn_or_proxy"`
		UserAgent                 string      `json:"user_agent"`
		Quantity                  int         `json:"quantity"`
		CouponID                  interface{} `json:"coupon_id"`
		CustomFields              struct {
			DiscordID int `json:"Discord ID"`
		} `json:"custom_fields"`
		DeveloperInvoice   bool        `json:"developer_invoice"`
		DeveloperTitle     interface{} `json:"developer_title"`
		DeveloperWebhook   interface{} `json:"developer_webhook"`
		DeveloperReturnURL interface{} `json:"developer_return_url"`
		Status             string      `json:"status"`
		Discount           int         `json:"discount"`
		FeePercentage      int         `json:"fee_percentage"`
		DayValue           int         `json:"day_value"`
		Day                string      `json:"day"`
		Month              string      `json:"month"`
		Year               int         `json:"year"`
		CreatedAt          int         `json:"created_at"`
		UpdatedAt          int         `json:"updated_at"`
		UpdatedBy          int         `json:"updated_by"`
		IPInfo             struct {
			ID             int     `json:"id"`
			RequestID      string  `json:"request_id"`
			IP             string  `json:"ip"`
			UserAgent      string  `json:"user_agent"`
			UserLanguage   string  `json:"user_language"`
			FraudScore     int     `json:"fraud_score"`
			CountryCode    string  `json:"country_code"`
			Region         string  `json:"region"`
			City           string  `json:"city"`
			Isp            string  `json:"isp"`
			Asn            int     `json:"asn"`
			Organization   string  `json:"organization"`
			Latitude       float64 `json:"latitude"`
			Longitude      float64 `json:"longitude"`
			IsCrawler      int     `json:"is_crawler"`
			Timezone       string  `json:"timezone"`
			Mobile         int     `json:"mobile"`
			Host           string  `json:"host"`
			Proxy          int     `json:"proxy"`
			Vpn            int     `json:"vpn"`
			Tor            int     `json:"tor"`
			ActiveVpn      int     `json:"active_vpn"`
			ActiveTor      int     `json:"active_tor"`
			RecentAbuse    int     `json:"recent_abuse"`
			BotStatus      int     `json:"bot_status"`
			ConnectionType string  `json:"connection_type"`
			AbuseVelocity  string  `json:"abuse_velocity"`
			CreatedAt      int     `json:"created_at"`
			UpdatedAt      int     `json:"updated_at"`
		} `json:"ip_info"`
		Serials  []string `json:"serials"`
		Webhooks []Webhook `json:"webhooks"`
		CryptoPayout            bool `json:"crypto_payout"`
		CryptoPayoutTransaction CryptoPayoutTx `json:"crypto_payout_transaction"`
		PaypalDispute interface{} `json:"paypal_dispute"`
		StatusHistory []StatusHistory `json:"status_history"`
		CryptoTransactions []CryptoTransaction `json:"crypto_transactions"`
	} `json:"data"`
}

type Webhook struct {
	Uniqid       string `json:"uniqid"`
	URL          string `json:"url"`
	Event        string `json:"event"`
	Retries      int    `json:"retries"`
	ResponseCode int    `json:"response_code"`
	CreatedAt    int    `json:"created_at"`
}

type CryptoPayoutTx struct {
	ToAddress    string  `json:"to_address"`
	FromAddress  string  `json:"from_address"`
	CryptoAmount float64 `json:"crypto_amount"`
	Hash         string  `json:"hash"`
	CreatedAt    int     `json:"created_at"`
}

type StatusHistory struct {
	ID        int    `json:"id"`
	InvoiceID string `json:"invoice_id"`
	Status    string `json:"status"`
	Details   string `json:"details"`
	CreatedAt int    `json:"created_at"`
}

type CryptoTransaction struct {
	CryptoAmount  float64 `json:"crypto_amount"`
	Hash          string  `json:"hash"`
	Confirmations int     `json:"confirmations"`
	CreatedAt     int     `json:"created_at"`
	UpdatedAt     int     `json:"updated_at"`
}