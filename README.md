Small Dynamic DNS update client
===

Install program
```shell
go get github.com/Agurato/dyndnsclient
```

Change the [config.yaml](config.yaml) file and launch program with

```shell
dyndnsclient /path/to/config.yaml
```

To run in cron (`crontab -e`):
```cron
*/5 * * * * /path/to/dyndnsclient /path/to/config.yaml >> /path/to/dyndnsclient.log 2>&1
```