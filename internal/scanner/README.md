# Scanner

Nightshift scans and scales objects by use of scanners. A scanner is
responsible for:

* Discovery of existing objects
* Scaling objects
* Save and load of a state
* Watch for live changes

Currently there are two watcher modules:

* openshift - which scans, scales and watch OpenShift DeploymentConfig resources
* deployment - which scans, scales and watch Kubernetes Deployment resources
* statefulset - which scans, scales and watch Kubernetes/OpenShift Statefulset resources

To add a new scanner, implement a factory method that implements the factory
type, and register that method with a new type. This type will then be
available in the config as a new scanner type.

The scanner itself should implement the Scanner interface. The watch method is
optional, and when implemented will result in live updates. However, the
GetObjects method is called frequently as well (at the configured resync
interval, default 15minutes).

The current scanners are targeted at OpenShift (or Kubernetes), but there is
no limitation which platform a scanner can target. As long as the
Scale/GetObjects methods can be implemented, a basic scanner can be
implemented as well.
