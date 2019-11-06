package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIssueToken(t *testing.T) {
	user, actions, roles, err := testApp.userService.VerifyLogin("admin", "password")

	require.Nil(t, err)
	require.NotNil(t, user)
	require.NotNil(t, actions)
	require.NotNil(t, roles)

	t.Run("issue_token_success", func(t *testing.T) {

		token, err := testApp.tokenService.IssueToken(user, actions, roles, "remote_addr", "x_forwarded_for", "x_real_ip", "user_agent")

		require.Nil(t, err)
		require.NotNil(t, token)
	})
}
