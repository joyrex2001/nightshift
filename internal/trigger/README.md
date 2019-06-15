# Trigger

Nightshift is able to trigger modules when a scale event occurs. These triggers
can be used to e.g. disable alerting for the scaled down pods, or e.g. to start
a pipeline which will trigger a specific process.

To add a new trigger, implement a factory method that implements the factory
type, and register that method with a new type. This type will then be
available in the config as a new scanner type.

The trigger itself should implement the Trigger interface. The Execute method
is called when the trigger occurs. It will receive a list of scanner.Objects
which were affected during the scaling and caused this trigger.
