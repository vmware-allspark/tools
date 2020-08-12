# Istio Security Policy Generator

This directory contains information needed to create large scale security policies.

See the [Istio Security](https://istio.io/latest/docs/reference/config/security/) for more information about policies.

The default values of the policies are specifically made to work with the environment that is created in the setup of [Istio Performance Benchmarking](https://github.com/istio/tools/tree/master/perf/benchmark)

## Config file

To generate specific security policies begin by creating a json file that has the format of the struct below:

```go
"SecurityPolicy":
{
  "AuthZ":
  {
    "action":int,           // optional DENY/ALLOW. Default:DENY
    "numNamespaces":int,    // optional. Default:0
    "numPaths":int,         // optional. Default:0
    "numPolicies":int,      // optional. Default:0
    "numPrincipals":int,    // optional. Default:0
    "numSourceIP":int,      // optional. Default:0
    "numValues":int         // optional. Default:0
  },
  "namespace":string,       // optional, the namespace in which all the policies will be applied to. Default:twopods-istio
  "peerAuthN":
  {
    "mtlsMode":string,      // optional STRICT/DISABLE. Default:STRICT
    "numPolicies":int       // optional. Default:0
  }
}
```

An example config file that will create 2 AuthorizationPolicies that each contain 1 sourceIP rule and 3 paths rules is formed as follows in a file called config.json:

```json
{
    "AuthZ":
    {
        "numPolicies": 2,
        "numSourceIP":1,
        "numPaths":3
    }
}
```

To generate the policies from the config file one must pass in the filename into the configFile flag. For example to generate the policy that is described in the above json run:

```bash
go run generate_policies.go generate.go -configFile="config.json"
```

## AuthorizationPolicy

To create an AuthorizationPolicy policy one must create a json file with the AuthZ field as well as set the numPolicies >= 1.
One should also include at least 1 rule. In this example numSourceIP is set to 1.

```json
{
  "AuthZ":
  {
    "numPolicies":1,
    "numSourceIP":1
  }
}
```

Once the wanted json file is created (called config.json) to generate the policies we just need pass in the config.json file to the configFile flag.

```bash
go run generate_policies.go generate.go -configFile="config.json"
```

This will create an Authorization Policy as follows and print it out to the stdout.

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: test-AuthorizationPolicy-1
  namespace: twopods-istio
spec:
 action: DENY
 rules:
 - from:
   - source:
       ipBlocks:
       - 0.0.0.0
```

##

The values which can be used to create custom AuthorizationPolicies are as follows:

```go
  "AuthZ":
  {
    "action":int,           // optional DENY/ALLOW. Default:DENY
    "numNamespaces":int,    // optional. Default:0
    "numPaths":int,         // optional. Default:0
    "numPolicies":int,      // optional. Default:0
    "numPrincipals":int,    // optional. Default:0
    "numSourceIP":int,      // optional. Default:0
    "numValues":int         // optional. Default:0
  }
```

For more information see [AuthorizationPolicy Reference](https://istio.io/latest/docs/reference/config/security/authorization-policy/).

## PeerAuthentication

To create an PeerAuthentication policy one must create a json file with the PeerAuthN field as well as set the PeerAuthN.numPolicies >= 1.
By default this generates a mtlsMode of "STRICT" but this can be overwritten as follows.

```json
{
  "PeerAuthN":
  {
    "numPolicies":1,
    "mtlsMode":"DISABLE"
  }
}
```

Once the wanted json file is created (called config.json) to generate the policies we just need pass in the config.json file to the configFile flag.

```bash
go run generate_policies.go generate.go -configFile="config.json"
```

This will create a PeerAuthentication policy as follows and print it out to the stdout.

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: test-PeerAuthentication-1
  namespace: twopods-istio
spec:
 mtls:
   mode: DISABLE
```

##

The values which can be used to create custom AuthorizationPolicies are as follows:

```go
  "peerAuthN":
  {
    "mtlsMode":string,      // optional STRICT/DISABLE. Default:STRICT
    "numPolicies":int       // optional. Default:0
  }
```

For more information see [PeerAuthentication Reference](https://istio.io/latest/docs/reference/config/security/peer_authentication/).

## Examples

generate_policies.go also allows a user to create mutliple kinds of policies in one command.
To generate 1 AuthorizationPolicy with a principals rule and 1 PeerAuthorization policy with STRICT mtlsMode, create a json file called twoPolicies.json with the following data and then run the following command:

```json
{
    "mtlsMode":
    {
      "numPolicies":1,
      "numPrincipals":1
    },
    "peerAuthN":
    {
      "numPolicies":1
    }
}
```

```bash
 go run generate_policies.go generate.go -configFile="twoPolicies.json"
```

Which outputs the following yaml:

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: test-AuthorizationPolicy-1
  namespace: twopods-istio
spec:
 action: DENY
 rules:
 - from:
   - source:
       principals:
       - cluster.local/ns/twopods-istio/sa/Invalid-0
---
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: test-PeerAuthentication-1
  namespace: twopods-istio
spec:
 mtls:
   mode: STRICT
```

### Output to a yaml file

To create a large AuthorizationPolicy to an output .yaml file, create a file called largeConfig.json which contains the following data:

```json
{
  "AuthZ":
  {
    "numPolicies":1000,
    "numSourceIP":100,
    "numPaths":100,
    "numValues":100
  }
}
```

run the following command:

```bash
go run generate_policies.go generate.go -configFile="largeConfig.json" > largePolicy.yaml
```

### Apply the yaml file

To apply largePolicy.yaml that was just created to istio use the following command:

```bash
kubectl apply -f largePolicy.yaml
```

## Example 1

- By creating a config file called config.json with the following data, and then running the following command:

```json
{
  "AuthZ":
  {
    "numPolicies":10,
    "numSourceIP":10,
    "numPaths":2
  }
}
```

```bash
go run generate_policies.go generate.go -configFile="config.json" > authZPolicy.yaml
```

- This creates 10 AuthorizationPolicies which each contains 10 sourceIP's sources, 2 paths operations, and places the policies in authZPolicy.yaml.

## Example 2

- By creating a config file called config.json with the following data, and then running the following command:

```json
{
  "AuthZ":
  {
    "numPolicies":1,
    "numSourceIP":100,
    "numPaths":100,
    "numNamespaces":100
  }
}
```

```bash
go run generate_policies.go generate.go -configFile="config.json" > authZPolicy.yaml
```

- This creates 1 AuthorizationPolicy which contains 100 sourceIP's sources, 100 paths operations, 100 namespaces sources, and places the policy in authZPolicy.yaml.

## Example 3

- By creating a config file called config.json with the following data, and then running the following command:

```json
{
  "PeerAuthN":
  {
    "numPolicies":1,
    "mtlsMode":"DISABLE"
  }
}
```

```bash
go run generate_policies.go generate.go -generate_policy="PeerAuthentication:1,mtlsMode:DISABLE"
```

- This creates 1 PeerAuthentication policy which has the mtls mode set to DISABLE

## Cleanup

To remove the policies applied navigate to the generate_policies folder and run the following command (update "largePolicy.yaml" if applied to a different .yaml file):

```bash
kubectl delete -f largePolicy.yaml
```
