package aws_ses

import (
	"crypto/tls"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Client maintains a small set of keep-alive connections.
type Client struct {
	server     *mail.SMTPServer
	connection chan *mail.SMTPClient
}

// newpool prepares—but does *not* open—connections.
func NewSMTP(cfg SmtpConfig) *Client {
	s := mail.NewSMTPClient()

	s.Host = cfg.Host
	s.Port = cfg.Port
	s.Username = cfg.Username
	s.Password = cfg.Password
	s.Encryption = mail.EncryptionSTARTTLS // recommended
	s.Authentication = mail.AuthLogin      // works best with SES
	s.KeepAlive = true                     // reuse TCP + TLS
	s.ConnectTimeout = cfg.ConnectTimeout
	s.SendTimeout = cfg.SendTimeout
	s.TLSConfig = &tls.Config{ServerName: cfg.Host}

	return &Client{
		server:     s,
		connection: make(chan *mail.SMTPClient, cfg.MaxConns),
	}
}

// grab returns a ready client—or dials a new one if the pool is empty.
func (p *Client) grab() (*mail.SMTPClient, error) {
	select {
	case c := <-p.connection:
		return c, nil
	default:
		return p.server.Connect()
	}
}

// release puts the client back or closes it if the pool is full.
func (p *Client) release(c *mail.SMTPClient) {
	select {
	case p.connection <- c:
	default:
		c.Close()
	}
}

// send pushes a fully-built mail.MSG through SES with context support.
func (p *Client) Send(msg *mail.Email) error {
	c, err := p.grab()
	if err != nil {
		return err
	}
	defer p.release(c)
	return msg.Send(c)
}
