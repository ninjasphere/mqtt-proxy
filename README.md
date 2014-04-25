# mqtt-proxy

This service acts as a front end for mqtt servers peforming preauthentication, load balancing and rate limiting.

# setup

Create a table with tokens in it.

```sql
CREATE TABLE users (
	uid int(11) NOT NULL AUTO_INCREMENT,
	mqtt_id varchar(128),
	PRIMARY KEY (uid),
	UNIQUE KEY mqtt_id_UNIQUE (mqtt_id)
)
```

Configure rabbitmq as a backend for mqtt-proxy by editing or creating `/usr/local/etc/rabbitmq/rabbitmq.config`.

```
[{rabbit,        [{tcp_listeners,    {"0.0.0.0", 5672}}]},
 {rabbitmq_mqtt, [{default_user,     <<"guest">>},
                  {default_pass,     <<"guest">>},
                  {allow_anonymous,  true},
                  {vhost,            <<"/">>},
                  {exchange,         <<"amq.topic">>},
                  {subscription_ttl, 1800000},
                  {prefetch,         10},
                  {ssl_listeners,    []},
                  {tcp_listeners,    [2883]},
                  {tcp_listen_options, [binary,
                                        {packet,    raw},
                                        {reuseaddr, true},
                                        {backlog,   128},
                                        {nodelay,   true}]}]}
].
```

Enable plugins by modifying running the following commands.

```
rabbitmq-plugins enable rabbitmq_management
rabbitmq-plugins enable rabbitmq_mqtt
rabbitmq-plugins enable rabbitmq_tracing
```

Restart rabbitmq.

```
rabbitmqctl stop
rabbitmq-server -detached
```

# Port redirection for 443 -> WS port

iptables -A PREROUTING -t nat -i eth0 -p tcp --dport 443 -j REDIRECT --to-port 9000

# status

* Preauthentication supports MySQL at the moment

