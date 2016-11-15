package main

func main() {
	// read config to server
	cnf, err := config.NewConfig("ini", "config.conf")
	if err != nil {
		fmt.Errorf("load config failed. err:%s", err.error())
		return
	}
	frame.Init()

	frame.Router("/healthCheck")
	frame.Router("/dataPerfman")
	frame.ConfigTranslater("client", "ProtoBufferTranslater")
	frame.ConfigTranslater("msg_center", "JsonTranslater")

	frame.ConfigHandler("client", "ClientHandler")
	frame.ConfigHandler("msg_center", "MsgCenterHandler")
	frame.ConfigHandler("push", "PushHandler")
	frame.ConfigHandler("redis", "RedisHandler")

	// http for maintain
	frame.AddListenAddr("client", "http", ":8080")
	// tcp for server
	frame.AddListenAddr("client", "tcp", ":9001")
	frame.AddListenAddr("msg_center", "tcp", ":9002")
	frame.AddConnectAddr("push_server:1", "tcp", ":9003")
	frame.AddConnectAddr("push_server:2", "tcp", ":9004")
	frame.Serve()

}
