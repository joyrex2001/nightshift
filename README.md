# Nightshift

Nightshift is a service that will enable automatic down and upscaling of
deployments within an OpenShift project to save resource usage (or use
resources for something else). Typically this service will run in a container
in a seperate namespace, where it will monitor and scale the namespaces
according to a preset configuration.

## Install in OpenShift

Nightshift can be run both from a seperate namespace, or in the namespace which
it will schedule. In order to allow OpenShift scaling the services, it needs to
run with a service account that has both read and edit permissions on the
namespace it should control. Adding the service account can be done with the
below commands, where a service account named "nightshift" is created in the
project specified with "source" which has access to the project "target". If
multiple projects should be scaled by nightshift, the policies should be added
for each project individually.

```bash
oc create sa nightshift -n <source>
oc policy add-role-to-user view system:serviceaccount:<target>:nightshift -n <source>
oc policy add-role-to-user edit system:serviceaccount:<target>:nightshift -n <source>
```

This repository also includes an example OpenShift template, which contains
the basis configuration for installing the service. The configuration is stored
in a configmap, the example folder contains an example for this as well. The
template includes the OpenShift oauth proxy for authentication with an
OpenShift account. Access is restricted to users that have update permissions
in the project where nightshift is deployed. The example template pulls the
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
added to a pod with an annotation, or can be defined in a config file. A
schedule definition contains out of 3 elements; a section defining the day(s)
in the week this schedule applies to, a time and the action that should be
taken.

The first section defines the day(s) of the week. This can be a combination
of single day(s), a range of days, or both. Each day is specified by the first
three letters of the English name. If multiple days, or ranges are specified,
they should be seperated by a comma. A range is specified by two days seperated
by a hyphen.

The second part defines the time. The time is specified in the timezone that
has been configured in the configuration file (default is Local, which usually
equals to UTC in most deployments).

The last part defines the action that needs to be taken in this time event. At
this point only the number of replicas can bet specified.

An example of a schedule configuration is: ```Mon-Wed,Fri 9:00 replicas=1```.

### Annotations

Nightshift can be configured by both a configuration file, as well as
annotations that should be set on a DeploymentConfig. These annotations will
always override the configuration specified in the configuration file. The
annotations supported by nightshift are:

* ```joyrex2001.com/nightshift.schedule``` in which a schedule can be specified.
Multiple schedules are allowed, and should be seperated with a semicolon.
* ```joyrex2001.com/nightshift.ignore``` which can be set to ```true``` to
ignore this deployment.

### Configuration file

The most flexible way to configure nightshift is by using a configuration file.
The configuration file allow complex schedules, with both default schedules as
well as allowing to fix schedules for certain label criteria, specified with
a selector.

In the below example , namespace ```development``` will be scanned for
annotations. If a deployment in this namespace doesn't have any schedule
annotations, it will apply the default schedule. For a deployment that matches
the label ```app=shell```, no schedule will be applied at all.

```
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

See the examples folder for another example, which also includes basic
nightshift configuration.

## See also

* Build status: [![CircleCI](https://circleci.com/gh/joyrex2001/nightshift.svg?style=svg)](https://circleci.com/gh/joyrex2001/nightshift)
* https://hub.docker.com/r/joyrex2001/nightshift
