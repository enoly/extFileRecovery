meta:
  id: descriptor
  endian: be
seq:
  - id: header
    size: 12
  - id: descriptors
    type: descriptor_record
    repeat: until
    repeat-until: _.flags & 0x08 == 0x08
types:
  descriptor_record:
    seq:
      - id: fs_block_num
        type: u4
      - id: flags
        type: u4
      - id: uuid
        size: 'flags & 0x02 == 0x02 ? 0 : 16'