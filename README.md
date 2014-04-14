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

# status

* Preauthentication supports MySQL at the moment

