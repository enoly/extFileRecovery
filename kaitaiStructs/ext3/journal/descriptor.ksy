meta:
  id: descriptor
  endian: be
seq:
  - id: header
    size: 12
  - id: descriptors
    type: descriptor_record
    repeat: until
    repeat-until: _.flags.last_record == true
types:
  descriptor_record:
    seq:
      - id: fs_block_num
        type: u4
      - id: flags
        size: 4
        type: descriptor_flags
      - id: uuid
        size: 16
  descriptor_flags:
    seq:
      - id: reserved
        type: b28
      - id: last_record
        type: b1
      - id: deleted_by_transaction
        type: b1
      - id: same_uuid
        type: b1
      - id: special_handling
        type: b1