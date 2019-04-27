# Namespace-manager
Automate namespace management ( like; delete after certain hours, or create snapshot of a pvc, etc )

### Current Resources
1. AutoKill: Deletes a namespace that this resource gets deployed into. You can set up how long to wait for before the controller deletes the namespace. There is also an option to delete all associated helm releases with it.
