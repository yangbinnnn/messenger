package sender

import (
	_tls "crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type MailClient struct {
	addr     string
	username string
	password string
	from     string
	timeout  int
	tls      bool
	insecure bool
}

func NewMailClient(addr, username, password, from string, timeout int, tls, insecure bool) *MailClient {
	return &MailClient{
		addr:     addr,
		from:     from,
		timeout:  timeout,
		username: username,
		password: password,
		tls:      tls,
		insecure: insecure,
	}
}

func (cli MailClient) connect() (*smtp.Client, error) {
	conn, err := net.DialTimeout("tcp", cli.addr, time.Duration(cli.timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	host, _, err := net.SplitHostPort(cli.addr)
	if err != nil {
		return nil, err
	}
	if cli.tls {
		tlsConn := _tls.Client(conn, &_tls.Config{
			ServerName:         host,
			InsecureSkipVerify: cli.insecure,
		})
		if err = tlsConn.Handshake(); err != nil {
			return nil, err
		}
		conn = tlsConn
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, err
	}
	if ok, _ := client.Extension("AUTH"); ok {
		err = client.Auth(smtp.PlainAuth("", cli.username, cli.password, host))
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func (cli MailClient) Send(tos []string, subject, message string) error {
	conn, err := cli.connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Mail(cli.from); err != nil {
		return err
	}
	for _, t := range tos {
		if err := conn.Rcpt(t); err != nil {
			return err
		}
	}
	w, err := conn.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	template := "From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\n%s\r\n"
	data := fmt.Sprintf(template, cli.from, strings.Join(tos, ","), subject, message)
	_, err = w.Write([]byte(data))
	if err != nil {
		return err
	}
	conn.Quit()
	return nil
}
