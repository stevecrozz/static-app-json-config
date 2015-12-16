# Static App JsObj

This tiny web app serves up JSON. The JSON comes from two sources, 1) a static
JSON configuration file and 2) information from HTTP request headers.

To see how it works, let's take the following example config.json:
```json
{
  "mountPath": "/static-apps/",
  "port": 8080,
  "apps": {
    "app-one": {
      "value1": 2
    }
  },
  "dynamic": {
    "requestHeaders": {
      "X-User-Id": "user"
    }
  }
}
```

Given a request like this one:
```bash
$ curl -H "X-User-Id: myuser" http://localhost:8080/static-apps/app-one
```

The response will be:
```
{
  "value1": 2,
  "user": "myuser"
}
```
