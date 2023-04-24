- External data

```shell
create /selly
```

```shell
create /selly/mongodb
create /selly/mongodb/warehouse
create /selly/mongodb/warehouse/uri "mongodb://localhost:27017"
create /selly/mongodb/warehouse/db_name "warehouse"
```

------------------------
```shell
create /selly/nats
create /selly/nats/warehouse
create /selly/nats/warehouse/uri "localhost:4222"
create /selly/nats/warehouse/api_key "selly"
create /selly/nats/warehouse/user ""
create /selly/nats/warehouse/password ""
create /selly/mongodb/warehouse/replica_set 
create /selly/mongodb/warehouse/ca_pem 
create /selly/mongodb/warehouse/cert_pem 
create /selly/mongodb/warehouse/cert_key_file_password 
create /selly/mongodb/warehouse/read_pref_mode 
```

------------------------
```shell
create /selly/mongodb
create /selly/mongodb/warehouse_audit
create /selly/mongodb/warehouse_audit/host "mongodb://localhost:27017"
create /selly/mongodb/warehouse_audit/db_name "audit-warehouse"

``` 

------------------------
- Location data

```shell
create /selly_warehouse
create /selly_warehouse/common
create /selly_warehouse/common/max_time_wait_for_global_care_update_status_in_minute 60
```

```shell
create /selly_warehouse
create /selly_warehouse/app
create /selly_warehouse/app/authentication
create /selly_warehouse/app/authentication/auth_secretkey "authsecretkey"
```