package atlanticnet

type Status string

const (
	StatusAwaitingCreation  Status = "AWAITING_CREATION"
	StatusCreating          Status = "CREATING"
	StatusFailed            Status = "FAILED"
	StatusReprovisioning    Status = "REPROVISIONING"
	StatusResettingPassword Status = "RESETTINGPWD"
	StatusRestarting        Status = "RESTARTING"
	StatusRunning           Status = "RUNNING"
	StatusStopped           Status = "STOPPED"
	StatusQueued            Status = "QUEUED"
	StatusRemoving          Status = "REMOVING"
	StatusRemoved           Status = "REMOVED"
	StatusResizingServer    Status = "RESIZINGSERVER"
	StatusSuspending        Status = "SUSPENDING"
	StatusSuspended         Status = "SUSPENDED"
)
