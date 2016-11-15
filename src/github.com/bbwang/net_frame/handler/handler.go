package handler

func init() {

}

type Handler interface {
	HandleMessage(content *Content, request *Request) error
}

var (
	handlers       = make(map[string]Handler)
	serverHandlers = make(map[string]Handler)
)

func Register(handlerName string, handler Handler) {
	if handler == nil {
		panic("handler: Register handler is nil")
	}
	if _, ok := adapters[handlerName]; ok {
		panic("handler: Register called twice for adapter " + handlerName)
	}
	handlers[handlerName] = handler
}

func NewHandler(serviceName, handlerName string) (Handler, error) {
	handler, ok := handlers[handlerName]
	if !ok {
		return nil, fmt.Errorf("config: unknown handler name %q (forgotten import?)", handlerName)
	}
	return handler, nil
}

func ConfigHandler(serviceName, handlerName string) {
	if _, ok := serverHandlers[serviceName]; ok {
		panic("handler: ConfigHandler called twice for service name " + serviceName)
	}
	handler, err := NewHandler(serviceName, handlerName)
	if err != nil {
		panic("handler: ConfigHandler new handler failed, err:%s", err.Error())
	}
	serverHandlers[serviceName] = handler
}

func HandleMessage(serviceName string, content *Content, request *Request) error {
	Handler, ok := serverHandlers[serviceName]
	if !ok {
		panic("handler: recv message, but not config Handler for servicename:%s" + serviceName)
		return errors.NewError("recv message, but not config Handler")
	}
	return Handler.HandleMessage(content, request)
}
