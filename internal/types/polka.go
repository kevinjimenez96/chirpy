package types

import "github.com/google/uuid"

type PolkaWebHookReq struct {
	Event string `json:"event"`
	Data  struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}
