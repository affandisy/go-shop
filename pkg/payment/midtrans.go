package payment

import (
	"crypto/sha512"
	"encoding/hex"
	"log"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClient struct {
	snapClient  snap.Client
	serverKey   string
	environment midtrans.EnvironmentType
}

func NewMidtransClient(serverKey, clientKey, environment string) *MidtransClient {
	var env midtrans.EnvironmentType
	if environment == "production" {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}

	snapClient := snap.Client{}
	snapClient.New(serverKey, env)

	log.Printf("Midtrans initialized (Environment: %s)", environment)

	return &MidtransClient{
		snapClient:  snapClient,
		serverKey:   serverKey,
		environment: env,
	}
}

type CreateSnapTokenRequest struct {
	OrderID       string
	GrossAmount   int64
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	Items         []ItemDetail
}

type ItemDetail struct {
	ID       string
	Name     string
	Price    int64
	Quantity int32
}

func (m *MidtransClient) CreateSnapToken(req CreateSnapTokenRequest) (*snap.Response, error) {
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.GrossAmount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
			Phone: req.CustomerPhone,
		},
	}

	if len(req.Items) > 0 {
		items := make([]midtrans.ItemDetails, len(req.Items))
		for i, item := range req.Items {
			items[i] = midtrans.ItemDetails{
				ID:    item.ID,
				Name:  item.Name,
				Price: item.Price,
				Qty:   item.Quantity,
			}
		}
		snapReq.Items = &items
	}

	snapResp, err := m.snapClient.CreateTransaction(snapReq)
	if err != nil {
		log.Printf("Failed to create Snap token: %v", err)
		return nil, err
	}

	log.Printf("Snap token created for order %s: %s", req.OrderID, snapResp.Token)

	return snapResp, nil
}

func (m *MidtransClient) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	input := orderID + statusCode + grossAmount + m.serverKey
	hasher := sha512.New()
	hasher.Write([]byte(input))
	expectedSignature := hex.EncodeToString(hasher.Sum(nil))

	return expectedSignature == signatureKey
}
