## Header Validation Toolkit

This toolkit is designed to ensure your headers are working as intended through a proxy.

### Endpoints

* `/`
  * Echoes the headers the app received through the proxy.
  * Does not set any headers.
* `/diff`
  * Compares a target header file against the headers received through the proxy.

### Usage

1. `go build -v -o hs header-server/header-server.go`
1. `./hs -header-ref target-file.json`

### Expectation Header File

Because it's hard to compare lots of different headers, it's generally best to copy the `header-reference.json` file and then put in values you care about.

### Expected Behaviour with Diffs

By default, the `header-reference.json` file only contains registered HTTP/1.1 header fields, but vendors love to put custom fields in their requests.

This is what an untracked (non-standard, custom) header looks like:

```json
{
	"name": "X-Amzn-Trace-Id",
	"have": {},
	"want": null
}
```

The reason this shows up is because there is no matching header in in our expectation header file, so there's no way to compare it against anything. If you want to start tracking the header, just add the header to the expectations file and restart the app, and it will look more normal:

```json
{
	"name": "X-Amzn-Trace-Id.[0]",
	"have": "Root=1-5cd1e3a3-45b120fc9d944ffc16de3c4c",
	"want": ""
}
```

If you're wondering about why header names generally look like `Header.[0]`, it's because of how Go parses the header frame. Go sets the header frames to be `map[string][]string`, which means it looks a bit like this:

```json
{
  "Header": [
    "value",
    "another-value"
  ]
}
```

According to RFC2616: 

> Multiple message-header fields with the same field-name MAY be present in a message if and only if the entire field-value for that header field is defined as a comma-separated list [i.e., #(values)]. It MUST be possible to combine the multiple header fields into one "field-name: field-value" pair, without changing the semantics of the message, by appending each subsequent field-value to the first, each separated by a comma. The order in which header fields with the same field-name are received is therefore significant to the interpretation of the combined field value, and thus a proxy MUST NOT change the order of these field values when a message is forwarded.

So what this means is that each header can use whatever separator it wants, but it has to use the same separator every single time. The problem is that in practice, no one can agree on a separator, which means valid separators can be `\s`, `,`, `-`, `;`, or just about any other character.

Go's header parsing implementation is RFC-compliant, but in reality it doesn't actually separate any of the strings, which means all your index values will be `0`. If Go started splitting the header values into separate string fields as the RFC originally intended, you would end up seeing something like this:

```json
[
    {
        "name": "X-Amzn-Trace-Id.[0]",
        "have": "Root=1-5cd1e3a3-45b120fc9d944ffc16de3c4c",
        "want": ""
    },
    {
        "name": "X-Amzn-Trace-Id.[1]",
        "have": "Root=1-5cd1e3a3-45b120f1523243fc16de3c4c",
        "want": ""
    }
]
```

So this is expected behaviour.
