## cherryctl server reset-bmc-password

Reset server BMC password.

### Synopsis

Reset BMC password for the specified server.

```
cherryctl server reset-bmc-password ID [flags]
```

### Examples

```
  # Reset BMC password for the specified server:
  cherryctl server reset-bmc-password 12345
```

### Options

```
  -h, --help   help for reset-bmc-password
```

### Options inherited from parent commands

```
      --api-url string   Override default API endpoint (default "https://api.cherryservers.com/v1/")
      --config string    Path to JSON or YAML configuration file
      --context string   Specify a custom context name (default "default")
      --fields strings   Comma separated object field names to output in result. Fields can be used for list and get actions.
  -o, --output string    Output format (*table, json, yaml)
      --token string     API Token (CHERRY_AUTH_TOKEN)
```

### SEE ALSO

* [cherryctl server](cherryctl_server.md)	 - Server operations. For more information on server types check Product Docs: https://docs.cherryservers.com/knowledge/product-docs#compute

