package helpers

import (
	"context"

	"github.com/TarsCloud/TarsGo/tars/util/current"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
)

func GetAuth(ctx context.Context) (map[string]string, string, error) {
	metadata, ok := current.GetRequestContext(ctx)
	if !ok {
		return nil, "", codes.ErrRequestContextMissing
	}

	authz := metadata[consts.AuthorizationMetadata]
	if len(authz) == 0 {
		return nil, "", codes.ErrRequestAuthrizationMetaMissing
	}

	return metadata, authz, nil
}
