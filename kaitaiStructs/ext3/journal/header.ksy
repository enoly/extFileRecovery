meta:
    id: header
    endian: be
seq:
    - id: signature
      type: u4
    - id: block_type
      type: u4
      enum: block_type_enum
    - id: serial_number
      type: u4
enums:
    block_type_enum:
        1: descriptor_block
        2: commit_block
        3: superblock_v1
        4: superblock_v2
        5: revoke_block