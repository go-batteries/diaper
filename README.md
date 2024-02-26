## diaper

Its a wrapper around viper, with some conventions applied:

- this reads config only from dotenv files
- ENV key names are in lowercase
- Ability to add value providers.
- Bring your own struct for validations


Value providers are identified by the prefix in the value. There is no specific
criteria for prefix. Default prefix supports are:

- env (implemented)
- ssm (to implement)

To invoke the config loading, the config loader needs some configuration. There
are two possible ways to do that.

There are two ways to build the the loaders:

**With a providers.yaml**

```go
reader, err := os.Open(providersYAMLPath)
assertErr(err, "failed to read file")

providers := diaper.LoadProviders(reader)

loader := diaper.DiaperConfig{
    DefaultEnvFile: "app.env",
    Providers:      providers,
}
```

**Without providers.yaml**

```
providers := diaper.BuildProviders(CustomProvider{})
```


A `Provider`, should implement a `ValueProvider` interface:

```go
func Deref(value interface{}) (computedValue string)
```

** In all cases, if a Provider cannot dereference a env value, you must return
the original value, so that dereferencing can be chained**

**Ordering of Providers** , In the end there is a NoopProvider, which will
always return the original value.

