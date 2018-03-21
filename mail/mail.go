package mail

import (
	_tls "crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type Client struct {
	addr     string
	username string
	password string
	from     string
	timeout  int
	_client  *smtp.Client
}

func NewClient(addr, username, password, from string, timeout int, tls, insecure bool) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if tls {
		tlsConn := _tls.Client(conn, &_tls.Config{
			ServerName:         host,
			InsecureSkipVerify: insecure,
		})
		if err = tlsConn.Handshake(); err != nil {
			return nil, err
		}
		conn = tlsConn
	}
	_client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, err
	}
	if ok, _ := _client.Extension("AUTH"); ok {
		err = _client.Auth(smtp.PlainAuth("", username, password, host))
		if err != nil {
			return nil, err
		}
	}
	return &Client{_client: _client, addr: addr, from: from, timeout: timeout, username: username, password: password}, nil
}

func (cli *Client) Send(tos []string, subject, message string) error {
	defer cli._client.Close()
	if err := cli._client.Mail(cli.from); err != nil {
		return err
	}
	for _, t := range tos {
		if err := cli._client.Rcpt(t); err != nil {
			return err
		}
	}
	w, err := cli._client.Data()
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
	cli._client.Quit()
	return nil
}
