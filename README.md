# tCLI
## Unoffical Tines CLI

tCLI allows users to run their [Tines](tines.com) stories via a CLI. tCLI generates a CLI tool from configuration stored within the Tines tenant, and allows for JSON webhooks to be interacted with from the users terminal.

tCLI is runtime, tenant-cententic CLI,  which means that updates & changes to commands, and their configuration is managed by Tines tenant, downloaded and locally cached to ensure that the CLI does not need to be rebuilt on change, and provide user-specific functionality.

This approach allows for teams to quickly run stories without pages, Slack bots, etc. Additionally, it allows the response data to be easily accessible to other popular tools such as `jq`, `awk`, and `find`.

For further information run the `tcli help` command.

> This tool is a PoC, and not production suitable, yet. See the Planned Work section for required pre-release changes. Please feel free to raised feature requests and/or bugs.

# Configuration

Configuring tCLI requires two components, tCLI auth config, and a per-story configuration change by the tenant story owner.

Each story config will create a new arguement under the `cmd` section of the tool. i.e.

```sh
tcli cmd <foo>
```

Simple requests are usually unhelpful, as such you can extend the command with arguments, which you can define as optional or required in the tenant configuration i.e.

```sh
tcli cmd foo --user=bar --ip=192.168.24.156
```
This will send the following JSON payload to the required webhook.

```json
{
   "user": "bar",
   "ip": "192.168.0.1"
}
```


## CLI Auth Config
Setting up tCLI requires the user to generate an API key for their Tines account and create a YAML config auth file under `$HOME/.tcli/auth`.

Below is sample auth file:

```yaml
tenant_name: example-tenant-1234
api_key: <User API key>
```

## Story Config

To expose a story to the tCLI tool you need the following setup within your Tines tenant:
 - Tag the story with `tcli`
 - Create a resource object with a valid schema, see [Tines Resources](https://www.tines.com/docs/resources/).

### tCLI Schema
To create a schema you will need to create a resource named `tcli_<story ID>`, i.e. `tcli_1234`. This should map to a matching story with the `tcli` tag.

> If no matching story & resource schema are found the command will not be added to the config cache & will not be accessible to the tool.


Schema:
```json
{
    "cmd": "foo",                      // Required: This is the cli command value 
    "webhook_path": "/path/bar/baz",   // Required: This is the webhook path we will be targetting
    "description": "",                 // Optional
    "request": {                       // Required
        "method": "",                  // Required: Any valid HTTP verb.
        "required": [],                // Optional: Required keypairs to be sent to the story
        "optional": []                 // Optional: Optional keypairs to be send to the story
    },
    "response": {                      // Optional
        "format": ""                   // Optional: Define the output format of the story. If left blank we default to auto. Valid values: auto, text, json.
    } 
}
```
The reponse object is not required, unless you want to specific output format.


Below is an example which takes 2 arguments, `length` and `user`, which would be created within a resource call `tcli_1234` to match the story ID number:

```json
{
   "cmd": "foo",
   "webhook_path": "/webhook/bar/baz",
   "description": "This is a very useful description.",
   "request": {
      "method": "post",
      "required": [],
      "optional": [
         "length",
         "user"
      ]
   },
   "response":{
      "format":"text"
   }
}
```


# Planned Work
The following is work that needs completed prior to the tool being in a alpha worthy state.

- [x] Reworking of the config file URL -> Tenant slug.
- [ ] Authentication flow rework. 
- [ ] Implement log & log levels `-v/-vv`
- [x] Implement response parsing for custom output
- [ ] Track & migrate from custom Tines API -> SDK as functions become avaliable.
- [x] Add expiration of cache within the cache file
- [ ] Linting/tests/security checks
- [ ] Migrate to seperate modules
- [ ] Mature/extend request schema.
   - [ ] Add required/optional fields for client-side validation.
   - [ ] Add bool flag support
- [ ] Review if the `cmd` prefix is *actually* needed.
- [x] Rework URL -> Path in schema due to full URL risks.
- [ ] Add `--no-cache`, forcing redownload of all stories

