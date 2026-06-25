package event

const (
	WebSocketMessage_NewMessage = iota
	WebSocketMessage_MessageUpdated
	WebSocketMessage_MessageDeleted
	WebSocketMessage_ReactionAdded
	WebSocketMessage_ReactionRemoved
	//WebSocketMessage_Typing
	WebSocketMessage_ReadReceipt
)
