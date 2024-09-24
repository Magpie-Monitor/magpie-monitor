# magpie-monitor

Reading logs is for the frogs, let's find insights from them

# To start development run:

`make watch`

# To connect to kafka container from terminal

1. Install kafkacat - https://github.com/edenhill/kcat
2. Set KAFKA_EXTERNAL_HOSTNAME variable in .env to 127.0.0.1
3. Connect with below command

`kcat -b 127.0.0.1:9094 -t pod -p 0 -o  -C -X sasl.username=username -X sasl.password=password -X sasl.mechanism=PLAIN -X security.protocol=SASL_PLAINTEXT -C`

Please note that when you change KAFKA_EXTERNAL_HOSTNAME, kafka container becomes unreachable for other docker containers, so this procedure is meant to be used only for local debugging. 