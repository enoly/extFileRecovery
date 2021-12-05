meta:
    id: superblock
    endian: be
seq:
    - id: header
      size: 12
    - id: block_size
      type: u4
    - id: blocks_count
      type: u4
    - id: first_data_block
      type: u4
    - id: first_transaction_number
      type: u4
    - id: first_transaction_block
      type: u4
    - id: error_number
      type: u4
    - id: feature_compatable
      type: u4
    - id: feature_incompatable
      type: u4
    - id: feature_read_only
      type: u4
    - id: uuid
      size: 16
    - id: file_system_count
      type: u4
    - id: superblock_copy
      type: u4
    - id: journal_blocks_per_transaction
      type: u4
    - id: system_blocks_per_transaction
      type: u4
    - id: unused
      size: 176
    - id: fs_uuids
      size: 768