package constants

//Status
const (
	StatusCreated = "created"
	StatusDone    = "done"
)

// Context
const (
	UserIdCtx    = "userId"
	UserRoleCtx  = "userRole"
	RequestIdCtx = "requestId"
	TraceIdCtx   = "traceId"
	SpanIdCtx    = "spanId"
	ApiNameCtx   = "apiName"
)

// Errors
const (
	BindBodyError      string = "bind_body"
	BindPathError      string = "bind_path"
	UserIdTypeMismatch string = "type_mismatch"
)

const (
	OrderCreatedEvent  string = "OrderCreated"
	StatusUpdatedEvent string = "StatusUpdated"
)
