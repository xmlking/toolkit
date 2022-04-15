# xds

Provide building blocks to bootstrap a **xDS server** to be used as _control plane_ for managing **gRPC Microservices**. <br/>
It can run in Kubernetes as well as local development environment, making it easy to use both in production and during development.

## Features

- [x] **Universal Control Plane**: Easy to use, distributed, runs anywhere on both Kubernetes and local development.
- [x] **Lightweight**: stateless, scale up/down 
- [x] **Traffic Routing**: With dynamic load-balancing for blue/green, canary, versioning and rollback deployments.
- [x] **Metrics**: export _prometheus_ complaint metrics via _OpenTelemetry_
- [x] **ADS**: Implements Aggregated Service Discovery (ADS)
- [ ] **SDS**: Implements Secret Discovery Service (SDS) [example](https://github.com/smallstep/step-sds)
- [x] **Virtual API Gateway**: Setup envoy ingress gateways via annotations
- [ ] **Health Checks**:
- [ ] **Multi-Tenancy**: 
- [ ] **Multi-Source**: Kubernetes, DNS, file watch, static

## Usage

Create xDS server as show in this [example](https://github.com/xmlking/grpc-starter-kit/tree/develop/service/xds) using this SDK.

You run the service locally or in Kubernetes accompanied by [config.yml](https://github.com/xmlking/grpc-starter-kit/blob/develop/config/config.yml#L78):

`make run-xds`

This use your local _kubeconfig_, so if **kubectl** works then it should work, unless you use some sort of authentication plugin.

If you're running in cluster it should also work without any environment variables. 
It requires some read-only access, which you can find the ClusterRole in [rbac.yml](./rbac.yml). 
It is recommended to deploy this as headless service. As it use DNS-based discovery, we recommend not to use autoscaling on this service but keep it always at max pods.

That's it! Your xDS server is ready.

### Connecting to xDS from clients
You need to set xDS bootstrap config on your client application. Here's the xDS bootstrap file:
```json
{
    "xds_servers": [
        {
            "server_uri": "localhost:5000",
            "channel_creds": [{"type": "insecure"}],
            "server_features": ["xds_v3"]
        }
    ],
    "node": {
        "id": "anything",
        "locality": {
            "zone" : "k8s"
        }
    }
}
```

Make sure to change `server_url` to wherever your application can access this xDS server. You then can supply this to your application in two methods:

- Put the entire JSON in `GRPC_XDS_BOOTSTRAP_CONFIG` environment variable
- Put the entire JSON in a file, then point `GRPC_XDS_BOOTSTRAP` environment variable to its path
Then follow the language specific instructions to enable xDS.

Finally, if you were connecting to `appname.appns:8080` write your connection string as `xds:///appname.appns:8080` instead. 
Note the triple slash and that the namespace is not optional. (As this doesn't use DNS)

### Go

Add this somewhere, maybe in your main file

```go
import _ "google.golang.org/grpc/xds"
```
### JavaScript

Install `@grpc/grpc-js-xds` then run

```js
require('@grpc/grpc-js-xds').register();
```

### Java
You need to add grpc-xds dependency along with the common grpc dependencies. gradle [example](https://github.com/xmlking/micro-apps/blob/develop/apps/account-service/build.gradle.kts#L10)

```kotlin
    implementation(libs.grpc.protobuf)
    implementation(libs.grpc.kotlin.stub)
    implementation(libs.grpc.services) // includes grpc-protobuf
    implementation(libs.grpc.xds) // includes grpc-services, grpc-auth,  grpc-alts
```
Then a new channel can be created with xds protocol.

```kotlin
Grpc.newChannelBuilder("xds:///{service}.{namespace}:{port}", InsecureChannelCredentials.create())
```

Note: the serviceConfigLookUp should not be disabled otherwise the xds protocol does not works correctly.

Since environment variable cannot be changed in java, there are 2 system properties which overrides the common bootstrap variables:

- io.grpc.xds.bootstrap to override GRPC_XDS_BOOTSTRAP
- io.grpc.xds.bootstrapConfig to override GRPC_XDS_BOOTSTRAP_CONFIG

### Virtual API Gateway

One feature of xDS is routing. This xDS server supports virtual API gateway by adding the following annotations to Kubernetes service:

```yaml
apiVersion: v1
kind: Service
metadata:
  # ...
  annotations:
    xds.chinthagunta.io/api-gateway: apigw1,apigw2
    xds.chinthagunta.io/grpc-service: package.name.ExampleService,package.name.Example2Service
```
The service also must have a port named `grpc`, which traffic will be sent to.

Then client applications (with xDS support) can connect to `xds:///apigw1` or `xds:///apigw2` (no port or namespace) and any API calls to gRPC service `package.name.ExampleService` and `package.name.Example2Service` will be sent to this service.
