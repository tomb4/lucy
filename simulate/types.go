package simulate

type (
	Config struct {
		PartyId  string
		HttpAddr string
		WsAddr   string
		GrpcAddr string
	}

	User struct {
		Id    string
		Token string
	}

	Resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	MetaJoinReq struct {
		PartyId   string `json:"partyId"`
		FollowUid string `json:"followUid"`
	}

	MetaExitReq struct {
		PartyId string `json:"partyId"`
	}
)
