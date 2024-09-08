package directive

import (
	"context"
	generatedAdmin "hshelby-tkcled-product/src/graph/generated/admin"

	"github.com/99designs/gqlgen/graphql"
)

var AdminDirective = generatedAdmin.DirectiveRoot{
	RequiredAuth: func(ctx context.Context, obj interface{}, next graphql.Resolver, action *string) (res interface{}, err error) {
		// if !network.HasToken(ctx) {
		// 	return nil, fmt.Errorf("unauthorized")
		// }

		// tokenStr := network.Token(ctx)
		// result, err := grpc_client.AuthenticatorClient().TokenVerify(ctx, &authenticator.TokenVerifyRequest{
		// 	JwtToken: tokenStr,
		// })
		// if err != nil || result == nil {
		// 	return nil, err
		// }

		// ctx = context.WithValue(ctx, "account_id", result.AccountId)
		// ctx = context.WithValue(ctx, "username", result.Username)
		// ctx = context.WithValue(ctx, "status", result.Status)
		// ctx = context.WithValue(ctx, "staff_id", result.StaffId)
		// ctx = context.WithValue(ctx, "workspace_id", result.WorkspaceId)

		// if action == nil {
		// 	return next(ctx)
		// }

		return next(ctx)
	},
}
