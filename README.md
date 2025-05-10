# tCLI
## Unoffical Tines CLI

tCLI allows users to run their [Tines](tines.com) stories via a CLI. tCLI generates a CLI tool from configuration stored within the Tines tenant, and allows for JSON webhooks to be converted into a tab completable CLI.

tCLI is tenant-cententic CLI,  which means that updates & changes to commands, and their configuration is done on the server, downloaded and locally cached to ensure that the CLI does not need to be rebuilt on change, and provide user-specific functionality.

This approach allows for teams to quick run stories without pages, Slack bots, etc. Additionally it allows the response data to be easily accessible to other popular tools such as `jq`, `awk`, and `find`.

This tool is a PoC, and not production suitable, yet. See the Planned Work section for pre-release changes.

For further information run the `tcli help` command.


# Configuration

Configuring tCLI requires two components, tCLI binary config, and a per-story configuration change by the tenant story owner.

Each story config will create a new arguement in the `cmd` section of the tool. i.e.

```sh
tcli cmd foo
```

Simple requests are usually unhelpful, as such you can extend the command with arguments, which you can define as optional or required in the tenant configuration i.e.

```sh
tcli cmd foo --user=bar --ip=192.168.24.156
```
This will send the following JSON payload, to the required webhook.
```json
{
   "user": "bar",
   "ip": "192.168.0.1"
}
```


## CLI Config
Setting up tCLI requires each user to generate an API for their user and configure the auth file under `$HOME/.tcli/auth`
The auth file will look like:
```yaml
tenantURL: https://example-tenant-1234.tines.com/
apiKey: <User API key here>
```

## Tenant config

To create a "tCLI enabled" you must implement the follow 3 changes to the story:
 - Tag the story with `tcli`
 - Created schema resources
 - Implement request/response inline with schema

The schema for story configs is below, along with an example:

Schema:
```json
{
   "foo": "TODO"
}
```

Example
```json
{
  "cmd": "foo",
  "url": "<redacted>",
  "description": "This is in a resource."
}
```

### Creating tCLI Schema
To create a schema you will need to create a resource named `tcli_<story ID>`, i.e. `tcli_2561`. This should map to a matching story with the `tcli` tag.

If not matching story/resource schema are found the command will not be added to the config cache & will not be accessible to the tool.


# Planned Work
The following is work that needs completed prior to the tool being in a alpha worthy state.

- [x] Reworking of the config file URL -> Tenant slug.
- [ ] Authentication flow rework. 
- [ ] Implement log & log levels `-v/-vv`
- [ ] Implement response parsing for custom output
- [ ] Track & migrate from custom Tines API -> SDK as functions become avaliable.
- [x] Add expiration of cache within the cache file
- [ ] Linting/tests/security auditing
- [ ] Migrate to seperate modules
- [ ] Mature/extend request schema.
- [ ] Review if the `cmd` prefix is *actually* needed.
- [ ] Rework URL -> Path in schema due to full URL risks.
- [ ] Add `--no-cache`, forcing redownload of all stories

