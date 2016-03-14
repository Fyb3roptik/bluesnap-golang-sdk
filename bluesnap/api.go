package bluesnap

import "fmt"
import "os"
// You need to change myApp to your app folder
import "myApp/app/models"
import "strings"
import "errors"
import (
	"net/http"
	"net/http/httputil"
)
import "bytes"
import "io/ioutil"
import "log"
import "io"
import "encoding/xml"

var createShopperStr = `<vaulted-shopper xmlns="http://ws.plimus.com">
   <first-name>%s</first-name>
   <last-name>%s</last-name>
   <payment-sources>
      <credit-card-info>
         <credit-card>
            <encrypted-card-number>%s</encrypted-card-number>
            <encrypted-security-code>%s</encrypted-security-code>
            <card-type>%s</card-type>
            <expiration-month>%s</expiration-month>
            <expiration-year>%s</expiration-year>
         </credit-card>
      </credit-card-info>
   </payment-sources>
</vaulted-shopper>`

var createSkuStr = `<catalog-sku xmlns="http://ws.plimus.com">
  <contract-name>One time charge - USD base</contract-name>
  <product-id>%s</product-id>
  <sku-type>DIGITAL</sku-type>
  <pricing-settings>
    <charge-policy-type>ONE TIME PAYMENT</charge-policy-type>
    <charge-policy>
      <one-time-charge>
        <catalog-prices>
          <catalog-price>
            <base-price>true</base-price>
            <currency>USD</currency>
            <amount>%.2f</amount>
          </catalog-price>
        </catalog-prices>
      </one-time-charge>
    </charge-policy>
  </pricing-settings>
</catalog-sku>`

// Need to change remote-host
var createOrderStr = `<order xmlns="http://ws.plimus.com">
  <ordering-shopper>
    <shopper-id>%s</shopper-id>
    <web-info>
      <ip>%s</ip>
      <remote-host>www.myapp.com</remote-host>
      <user-agent>Mozilla/5.0 (Linux; X11)</user-agent>
    </web-info>
  </ordering-shopper>
  <cart>
    <cart-item>
      <sku>
        <sku-id>%s</sku-id>
      </sku>
      <quantity>1</quantity>
    </cart-item>
  </cart>
  <expected-total-price>
    <amount>%.2f</amount>
    <currency>USD</currency>
  </expected-total-price>
</order>`

var createTransactionStr = `<card-transaction xmlns="http://ws.plimus.com">
   <card-transaction-type>AUTH_CAPTURE</card-transaction-type>
   <recurring-transaction>ECOMMERCE</recurring-transaction>
   <soft-descriptor>DescTest</soft-descriptor>
   <amount>%.2f</amount>
   <currency>USD</currency>
   <vaulted-shopper-id>%s</vaulted-shopper-id>
</card-transaction>`

type BlueSnap struct {
	STOREID     string
	BSAPIUSER   string
	BSAPIPASS   string
	BSENDPOINT  string
	BSPRODUCTID string
	Test        bool
}

func (b *BlueSnap) Init() {
	if b.Test {
		b.BSAPIUSER = os.Getenv("TEST_BS_API_USER")
		b.BSAPIPASS = os.Getenv("TEST_BS_PASS")
		b.STOREID = os.Getenv("TEST_BS_STORE_ID")
		b.BSENDPOINT = os.Getenv("TEST_BS_ENDPOINT")
		b.BSPRODUCTID = os.Getenv("TEST_BS_PRODUCT_ID")
	} else {
		b.BSAPIUSER = os.Getenv("BS_API_USER")
		b.BSAPIPASS = os.Getenv("BS_PASS")
		b.STOREID = os.Getenv("BS_STORE_ID")
		b.BSENDPOINT = os.Getenv("BS_ENDPOINT")
		b.BSPRODUCTID = os.Getenv("BS_PRODUCT_ID")
	}
}

func (b *BlueSnap) CreateVaultedShopper(bsCardToken, bsSecurityToken, expMonth, expYear string, cardType string, user *models.User) (string, error) {
	names := strings.Split(user.Name, " ")
	firstName := names[0]
	lastName := firstName
	if len(names) > 1 {
		lastName = names[1]
	}
	createShopperReq := fmt.Sprintf(createShopperStr, firstName, lastName,
		bsCardToken, bsSecurityToken, cardType,
		expMonth, expYear)
	return b.performHttpPostForCreate("vaulted-shoppers", createShopperReq)
}

func (b *BlueSnap) CreateTransaction(shopperId string, amount float64) (string, error) {
	return b.performHttpPost("transactions", fmt.Sprintf(createTransactionStr, amount, shopperId))
}

func (b *BlueSnap) performHttpPostForCreate(resource, body string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", b.BSENDPOINT, resource), bytes.NewReader([]byte(body)))
	req.Header.Add("Content-Type", "application/xml")
	log.Println("user: ", b.BSAPIUSER, " pass:", b.BSAPIPASS)
	req.SetBasicAuth(b.BSAPIUSER, b.BSAPIPASS)
	log.Println("headers ", req.Header)
	log.Println("Body: ", body)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes))
	defer resp.Body.Close()
	u, err := resp.Location()
	if err != nil {
		return "", errors.New("Invalid card data")
	}
	paths := strings.Split(u.Path, "/")
	if len(paths) == 0 {
		return "", errors.New("error in returned url format")
	}
	return paths[len(paths)-1], nil
}

func (b *BlueSnap) performHttpPost(resource, body string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", b.BSENDPOINT, resource), bytes.NewReader([]byte(body)))
	req.Header.Add("Content-Type", "application/xml")
	log.Println("user: ", b.BSAPIUSER, " pass:", b.BSAPIPASS)
	req.SetBasicAuth(b.BSAPIUSER, b.BSAPIPASS)
	log.Println("headers ", req.Header)
	log.Println("Body: ", body)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes))
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return "", nil
	} else {
		return "", errors.New("Error in transaction. Please try again.")
	}
}

func (b *BlueSnap) performHttpGet(resource string) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", b.BSENDPOINT, resource), nil)
	req.Header.Add("Content-Type", "application/xml")
	log.Println("user: ", b.BSAPIUSER, " pass:", b.BSAPIPASS)
	req.SetBasicAuth(b.BSAPIUSER, b.BSAPIPASS)
	log.Println("headers ", req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return resp.Body, err
	}
	dump, _ := httputil.DumpResponse(resp, true)
	log.Println("get api dump ", string(dump))
	if resp.StatusCode == http.StatusOK {
		return resp.Body, nil
	} else {
		return resp.Body, errors.New("Error in transaction. Please try again.")
	}
}

func (b *BlueSnap) GetCards(shopperId string) (c *CardList, err error) {
	rc, err := b.performHttpGet(fmt.Sprintf("vaulted-shoppers/%s", shopperId))
	if err != nil {
		return
	}
	c = &CardList{}
	defer rc.Close()
	d := xml.NewDecoder(rc)
	for {
		t, _ := d.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			log.Println(se.Name.Local)
			if se.Name.Local == "card-last-four-digits" {
				card := Card{}
				no := ""
				d.DecodeElement(&no, &se)
				card.LastFour = no
				c.Values = append(c.Values, &card)
				break
			}
		}
	}
	return
}

type Card struct {
	ID          string `json:"id"`
	Month       uint8  `json:"exp_month"`
	Year        uint16 `json:"exp_year"`
	Fingerprint string `json:"fingerprint"`
	LastFour    string `json:"last4"`
	Brand       string `json:"brand"`
	City        string `json:"address_city"`
	Country     string `json:"address_country"`
	Address1    string `json:"address_line1"`
	Address2    string `json:"address_line2"`
	State       string `json:"address_state"`
	Zip         string `json:"address_zip"`
	CardCountry string `json:"country"`
	Name        string `json:"name"`
	DynLastFour string `json:"dynamic_last4"`
	Deleted     bool   `json:"deleted"`
}

// CardList is a list object for cards.
type CardList struct {
	Values []*Card `json:"data"`
}
