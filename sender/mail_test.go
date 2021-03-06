package sender

import (
	"os"
	"strconv"
	"testing"
)

func TestEmail(t *testing.T) {
	addr := os.Getenv("EMAIL_ADDR")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	tls, _ := strconv.ParseBool(os.Getenv("EMAIL_TLS"))
	timeout := 60
	insecure := false
	from := username
	tos := []string{"yangbin@cloudhua.com"}
	subject := "hello"
	message := "email test"

	client := NewMailClient(addr, username, password, from, timeout, tls, insecure)
	err := client.Send(tos, subject, message)
	if err != nil {
		t.Fatal(err)
	}
}
