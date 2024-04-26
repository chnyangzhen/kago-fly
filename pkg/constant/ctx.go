package constant

const (
	// RequestScopeContextKey is the key of request scope context，value type is component.WebContext
	RequestScopeContextKey = "request-scope-context"
	// XRequestID 指定请求ID
	XRequestID = "X-Request-ID"
	// XRequestGroup 指定请求分组，如dubbo服务的group
	XRequestGroup = "X-Request-Group"
	// XRequestVersion 指定请求版本，如dubbo服务的version
	XRequestVersion = "X-Request-Version"
	// XVersion 指定请求版本，如注册的接口的version
	XVersion = "X-Version"
	// ContentType 指定请求的Content-Type
	ContentType = "Content-Type"
)
