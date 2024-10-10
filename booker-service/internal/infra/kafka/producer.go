package kafka

type Producer struct {
	ConfigMap *ckafka.ConfigMap
}

func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{
		ConfigMap: configMap,
	}
}
