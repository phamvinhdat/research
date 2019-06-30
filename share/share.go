package share

type Configuration struct{
	AMQPConnectionURL string
}

type AddTask struct{
	Number1 int
	Number2 int
}

func NewConfiguration() *Configuration{
	return &Configuration{
		AMQPConnectionURL: "amqp://guest:guest@localhost:5672",
	}
}