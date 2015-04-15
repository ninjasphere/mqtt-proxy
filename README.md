# mqtt-proxy

This service acts as a front end for mqtt servers peforming preauthentication, load balancing and rate limiting.

# setup

Create a table with tokens in it.

```sql
CREATE TABLE legends (
	uid int(11) NOT NULL AUTO_INCREMENT,
	mqtt_id varchar(128) COLLATE utf8_bin NOT NULL,
	PRIMARY KEY (uid),
	UNIQUE KEY mqtt_id_UNIQUE (mqtt_id)
) DEFAULT CHARSET=utf8;

CREATE TABLE tokens (
  token_id int(11) NOT NULL AUTO_INCREMENT,
  uid int(11) NOT NULL,
  token varchar(64) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (token_id),
  UNIQUE KEY token_UNIQUE (token)
) DEFAULT CHARSET=utf8;
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

# Licensing

mqtt-proxy is licensed under the MIT License. See LICENSE for the full license text.
