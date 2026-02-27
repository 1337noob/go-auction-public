package event_store

type SnapshotPolicy interface {
	ShouldTakeSnapshot(currentVersion int) bool
}

type CountBasedSnapshotPolicy struct {
	Interval int
}

func NewCountBasedSnapshotPolicy(interval int) *CountBasedSnapshotPolicy {
	return &CountBasedSnapshotPolicy{
		Interval: interval,
	}
}

func (c CountBasedSnapshotPolicy) ShouldTakeSnapshot(currentVersion int) bool {
	//log.Println("currentVersion:", currentVersion)
	return currentVersion%c.Interval == 0
}
