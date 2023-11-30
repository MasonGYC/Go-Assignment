package main

const (
	// message types
	// from server to manager
	RD_REQ  = "READ_REQUEST"
	WR_REQ  = "WRITE_REQUEST"
	RD_CFM  = "READ_CONFIRM"
	WR_CFM  = "WRITE_CONFIRM"
	INV_CFM = "INVALIDATE_COPY_CONFIRM"

	// from manager to server
	RD_FWD  = "READ_FORWARD"
	WR_FWD  = "WRITE_FORWARD"
	INV_REQ = "INVALIDATE_COPY_REQUEST"

	// from server to server
	SD_RD_PAGE = "SEND_READ_PAGE"
	SD_WR_PAGE = "SEND_WRITE_PAGE"

	// operations
	READ  = "READ"
	WRITE = "WRITE"
	INV   = "INVALIDATE_COPY"

	// manager monitor
	NTC_WR_REQ = "NOTICE_WRITE_REQUEST"
)

type Message struct {
	sender_id    int
	receiver_id  int
	page_num     int
	message_type string
	clock        int // vector clock
	page         Page
	requester    int // server to forward to
}

// request to r/w page x
func RequestMessage(sender_id int, page_num int, clock int, operation string) Message {
	var message_type string
	if operation == READ {
		message_type = RD_REQ
	} else if operation == WRITE {
		message_type = WR_REQ
	}
	return Message{
		sender_id:    sender_id,
		receiver_id:  CM_ID,
		page_num:     page_num,
		message_type: message_type,
		clock:        clock,
		requester:    sender_id,
	}
}

// Forward r/w request for page x to owner server
func ForwardMessage(requester int, receiver_id int, page_num int, clock int, operation string) Message {
	var message_type string
	if operation == READ {
		message_type = RD_FWD
	} else if operation == WRITE {
		message_type = WR_FWD
	}
	return Message{
		sender_id:    CM_ID,
		receiver_id:  receiver_id,
		page_num:     page_num,
		message_type: message_type,
		clock:        clock,
		requester:    requester,
	}
}

// Send page to requested server
func SendPageMessage(sender_id int, receiver_id int, page Page, clock int, operation string) Message {
	var message_type string
	if operation == READ {
		message_type = SD_RD_PAGE
	} else if operation == WRITE {
		message_type = SD_WR_PAGE
	}
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		page_num:     page.id,
		message_type: message_type,
		clock:        clock,
		page:         page,
	}
}

// Comfirm inv/r/w to manager
func ConfirmMessage(requester int, sender_id int, page_num int, clock int, operation string) Message {
	var message_type string
	if operation == READ {
		message_type = RD_CFM
	} else if operation == WRITE {
		message_type = WR_CFM
	} else if operation == INV {
		message_type = INV_CFM
	}
	return Message{
		sender_id:    sender_id,
		receiver_id:  CM_ID,
		page_num:     page_num,
		message_type: message_type,
		clock:        clock,
		requester:    requester,
	}
}

// Send invalidate copy
func InvalidReqMessage(requester int, receiver_id int, page_num int, clock int) Message {
	return Message{
		sender_id:    CM_ID,
		receiver_id:  receiver_id,
		page_num:     page_num,
		message_type: INV_REQ,
		clock:        clock,
		requester:    requester,
	}
}

// notice monitor that a wr request is pushed to queue
func NoticeMessage(requester int, page_num int, clock int) Message {
	return Message{
		sender_id:    CM_ID,
		receiver_id:  CM_ID,
		page_num:     page_num,
		message_type: NTC_WR_REQ,
		clock:        clock,
		requester:    requester,
	}
}
