package connectionManager

import (

	"github.com/stretchr/testify/suite"
	
)

type RedisConnectionTestSuite struct {
	suite.Suite
	connection   DBoperationsFactory
}