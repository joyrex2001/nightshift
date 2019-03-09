# Nightshift

Nightshift is a service that will enable automatic down and upscaling of
deployments within an OpenShift project to save resource usage (or use
resources for something else). Typically this service will run in a container
in a seperate namespace, where it will monitor and scale the namespaces
according to a presert configuration.

## Install in OpenShift

It is advised to run the Nightshift agent in a separate namespace. In order to
allow OpenShift scaling the services, it needs to run with a service account
that has both read and edit permissions on the namespace it should control.
Adding the service account can be done with the below commands, where a
service account named "nightshift" is created in the project specified with
"source" which has access to the project "target". If multiple projects
should be scaled by nightshift, the policies should be added for each project
individually.

```bash
oc create sa nightshift -n <source>
oc policy add-role-to-user view system:serviceaccount:<target>:nightshift -n <source>
oc policy add-role-to-user edit system:serviceaccount:<target>:nightshift -n <source>
```

This repository also includes an example OpenShift template, which contains
the basis configuration for installing the service. The configuration is stored
in a configmap, which for convenience has been added to the template as well.
Note that for a production setup, this configmap should be adjusted to reflect
the required setup (see also the Configration section). The template pulls the
nightshift image from docker hub.

```bash
oc create configmap nightshift-config --from-file=examples/config.yaml
oc process -f examples/openshift.yaml | oc apply -f -
oc rollout latest dc/nightshift
```

## Configuration

Nightshift can be configured with both annotations, or via a config file. They
can also be combined, where annotations always override any other configuration
set.

### Schedule

Nightshift scales deployments according to a schedule. A schedule can be
added to a pod with an annotation, or can be defined in a config file. An
example of a schedule configuration is: ```Mon-Fri  9:00 replicas=1```.

### Annotations

The annotations supported by nightshift are:

* ```joyrex2001.com/nightshift.schedule``` in which a schedule can be specified. Multiple schedules are allowed, and should be seperated with a ;.
* ```joyrex2001.com/nightshift.ignore``` which can be set to ```true``` to ignore this deployment.

### Configuration file

The best way to configure nightshift is by using a configuration file. The
configuration file will specify which

```
logging:
    threshold: "info"
    verbose: 3

openshift:
    namespace: "staging"

scanner:
    - namespace:
        - "development"
      default:
        schedule:
          - "Mon-Fri  9:00 replicas=1"
          - "Mon-Fri 18:00 replicas=0"
      deployment:
        - selector:
            - "app=shell"
          schedule:
            - ""
```

In this example, the namespace ```staging``` will be scanned for annotations
only. The namespace ```development``` will be scanned for annotations. If a
deployment doesn't have any schedule annotations, it will apply the default
schedule. For a deployment that matches the label ```app=shell```, no schedule
will be applied at all.
