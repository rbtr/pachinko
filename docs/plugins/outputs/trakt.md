### Trakt and the Trakt Collector output
Pachinko can interact with Trakt. Currently, Pachinko can:
- add sorted items to your Trakt collection (`trakt_collector` output plugin)

#### Trakt Authorization
To communicate with Trakt, it needs an access token. A helper command is included for authorizating:
```bash
$ pachinko trakt
Authenticating in Trakt!
Please open in your browser:    https://trakt.tv/activate
         and enter the code:            1234A1AA
```

Enter the provided code at the [link](https://trakt.tv/activate) and Pachinko will receive an access token. It will write the access credentials out to a file, by default `/etc/pachinko/trakt`. Specify a different file by using the `--authfile /path/to/file` flag on the `trakt` command.

The [authfile](../../examples/trakt) is JSON and contains authorization credentials.

#### Trakt Collector
To add items to your Trakt collection when Pachinko is done processing them, enable the Trakt Collector output plugin in your Pachinko config file. 

The only configurable option is the authfile location - point it at the authfile created by the authorization step as described [above](#trakt-authorization).

```yaml
#...
outputs:
- name: trakt-collector
  authfile: "/etc/pachinko/trakt"

#...
```

Now when Pachinko identifies and processes TV or Movies they will be automatically collected in Trakt!
