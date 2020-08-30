package sqlstorage

import (
	"fmt"
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"
	"net/smtp"
	"os"
	"time"
	"unicode"
)

var (
	storeShare = []byte(os.Getenv("STORE_SHARE"))
)

type PaymentRepository struct {
	storage *Storage
}

func (r *PaymentRepository) Make(p *model.Payment) error {
	if !valid(p.Card) {
		return storage.ErrCardIsInvalid
	}

	err := p.PostCreate()
	if err != nil {
		return err
	}

	r.storage.db.QueryRow(
		"INSERT INTO \"payment\" (id, time, game_name, user_email, seller_email, total, code) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		p.ID,
		time.Now(),
		p.GameName,
		p.UserEmail,
		p.SellerEmail,
		p.Total,
		p.Code,
	)

	// Sender data.
	from := fmt.Sprintf("%s", os.Getenv("EMAIL"))
	password := fmt.Sprintf("<%s>", os.Getenv("EMAIL_PASSWORD"))

	// Receiver email address.
	to := []string{
		"mishakukarkin@mail.ru",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("This is a test email message.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func valid(input string) bool {

	sum := 0
	counter := 0

	for _, r := range reverse(input) {

		if unicode.IsDigit(r) {

			val := int(r - '0')

			if counter%2 == 1 {
				val = val * 2

				if val > 9 {
					val = val - 9
				}
			}
			sum += val

			counter++
			continue
		}

		if unicode.IsSpace(r) {
			continue
		}

		return false
	}

	if counter < 2 {
		return false
	}

	return (sum % 10) == 0
}
