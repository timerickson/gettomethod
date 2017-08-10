gettomethod is a simple proxy from an HTTP GET to another method or verb.

The original intent and use is to allow webhook integration between something
(such as BitBucket) that only does GET webhooks, to a system (like GoCD) that
requires POST or other methods for their api. Additionally, it allows ignoring
SSL errors, providing Basic Auth, and specifying headers.

# Usage
Simply put, it listens on port 8080 for GET requests to a path corresponding
to an HTTP verb:

GET http://<gettomethod-host>:8080/post -> proxies a POST request
GET http://<gettomethod-host>:8080/put -> proxies a PUT request

The target of the proxied request is defined entirely by querystring
parameters:

## Example
```curl "http://<<gettomethod-host>>:8080/post?protocol=https&host=<<gocd-server-host>>&path=/go/api/pipelines/<<pipeline_name>>/schedule&debug&_Confirm=true&ignoreSslErrors&_Authorization_Basic=<gocd-user>:<gocd-password>"```

Will trigger <<pipeline_name>> on a GoCD server via it's API, assuming it is listening on 443 or 443 is mapped to its default 8154 (as mine is).

## Parameters
| parameter | info |
| --- | --- |
| debug | Enable debug mode to see debug information in response (this is a flag) |
| ignoreSslErrors | Ignore SSL certificate errors on the target (e.g. due to self-signed certs) |
| protocol | http or https |
| host | domain or hostname |
| port | target port (if not protocol default) |
| path | path or target request |
| body | body for PUT or POST request (can be skipped for no body) |
| Authorization_Basic | Generates an Authorization header with a Basic Auth formatted/encoded value of the user:pass provided in the parameter value |

## Headers
Headers can be specified via a querystring parameter with the header name prefixed by an underscore '_'. E.g. ```_Foo=bar``` produces a request header Foo with the value bar. Headers with underscore-prefixed names will simply require a double-underscore-prefixed parameter.

## Query Parameters
Query Parameters can be specified as normal, providing they do not begin with an underscore or have the same name as one of the gettomethod parameters. These will all be preserved.
