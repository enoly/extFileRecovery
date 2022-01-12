package extworker

type SuperblockInfo struct {
	InodeCount     uint64
	BlockCount     uint64
	BlockSize      uint64
	InodeSize      uint64
	BlocksPerGroup uint64
	InodesPerGroup uint64
	JournalInode   uint64
	UUID           []byte
	VolumeName     string
}
