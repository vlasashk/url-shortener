//go:build integration

package cmd

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vlasashk/url-shortener/test/suits"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(suits.IntegrationSuite))
}
