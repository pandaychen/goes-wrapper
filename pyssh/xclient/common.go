package xclient

import (
	"golang.org/x/crypto/ssh"
)

var (
	NotRecommendCiphers = []string{
		"arcfour256", "arcfour128", "arcfour",
		"aes128-cbc", "3des-cbc",
	}

	NotRecommendKeyExchanges = []string{
		"diffie-hellman-group1-sha1", "diffie-hellman-group-exchange-sha1",
		"diffie-hellman-group-exchange-sha256",
	}
)

func InitSSHConfig() ssh.Config {
	var (
		cfg ssh.Config
	)
	cfg.SetDefaults()
	cfg.Ciphers = append(cfg.Ciphers, NotRecommendCiphers...)
	cfg.KeyExchanges = append(cfg.KeyExchanges, NotRecommendKeyExchanges...)
	return cfg
}
