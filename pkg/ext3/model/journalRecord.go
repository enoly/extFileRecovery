package model

type JournalRecordType int

const (
	JournalSuperblock JournalRecordType = 0
	JournalDescriptor JournalRecordType = 1
	JournalMetadata   JournalRecordType = 2
	JournalCommit     JournalRecordType = 3
	JournalRevoke     JournalRecordType = 4
	JournalInvalid    JournalRecordType = 5
)

type JournalRecord struct {
	RecordNum   uint32
	RecordBlock uint32
	Type        JournalRecordType
	Description string
}
