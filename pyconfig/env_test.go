package env

import (
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

type MysqlConfig struct {
	MysqlHost        string        `default:"serv-mysql"`
	MysqlPort        string        `default:"3306"`
	MysqlUser        string        `default:"root"`
	MysqlPass        string        `default:"xxxx"`
	MysqlDbName      string        `default:"test"`
	MysqlConnTimeout time.Duration `default:"20s"`
}

func InitMysqlConfig() (*MysqlConfig, error) {
	cfg := new(MysqlConfig)
	IgnorePrefix()
	err := FillConfig(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestMysqlConfig(c *C) {
	m, _ := InitMysqlConfig()
	c.Assert(m.MysqlHost, Equals, "serv-mysql")
	c.Assert(m.MysqlPort, Equals, "3306")

	c.Assert(m.MysqlUser, Equals, "root")
	c.Assert(m.MysqlPass, Equals, "xxxx")

	c.Assert(m.MysqlDbName, Equals, "test")
	c.Assert(m.MysqlConnTimeout, Equals, time.Duration(20*time.Second))
}
