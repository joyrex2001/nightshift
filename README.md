# NightShift

NightShift is a service that will enable automatic down and upscaling of
deployments within an OpenShift project to save resource usage (or use
resources for something else).

## Configuration

### Schedule

NightShift scales deployments according to a schedule. A schedule can be
added to a pod with an annotation, or can be defined in a config file. An
example of a schedule configuration is: ```Mon-Fri  9:00 replicas=1```.

### Annotations

The annotations supported by NightShift are:

* ```joyrex2001.com/nightshift.schedule``` in which a schedule can be specified.
* ```joyrex2001.com/nightshift.ignore``` which can be set to ```true``` to ignore this deployment.

### Configuration file

The best way to configure NightShift is by using a configuration file. The
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
