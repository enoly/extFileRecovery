meta:
    id: ext3Superblock
    endian: le
seq:
    - id: inodes_count
      type: u4
    - id: blocks_count
      type: u4
    - id: reserved_blocks_count
      type: u4
    - id: free_blocks_count
      type: u4
    - id: free_inodes_count
      type: u4
    - id: first_data_block
      type: u4
    - id: log_block_size
      type: u4
    - id: log_frag_size
      type: u4
    - id: blocks_per_group
      type: u4
    - id: frags_per_group
      type: u4
    - id: inodes_per_group
      type: u4
    - id: mount_time
      type: u4
    - id: written_time
      type: u4
    - id: mount_count
      type: u2
    - id: max_mount_count
      type: u2
    - id: signature
      contents: [0x53, 0xef]
    - id: fs_state
      type: u2
      enum: state_enum
    - id: errors
      type: u2
      enum: errors_enum
    - id: minor_version
      type: u2
    - id: last_check
      type: u4
    - id: check_interval
      type: u4
    - id: creator_os
      type: u4
      enum: creator_os_enum
    - id: major_version
      type: u4
      enum: major_version_enum
    - id: def_reserved_uid
      type: u2
    - id: def_reserved_gid
      type: u2
    - id: first_inode
      type: u4
    - id: inode_size
      type: u2
    - id: block_group_copy_loc
      type: u2
    - id: feature_compatable
      type: u4
    - id: feature_incompatable
      type: u4
    - id: feature_read_only
      type: u4
    - id: uuid
      size: 16
    - id: volume_name
      size: 16
    - id: last_mounted
      size: 64
    - id: algo_bitmap
      type: u4
    - id: prealloc_blocks
      type: u1
    - id: prealloc_dir_blocks
      type: u1
    - id: padding1
      size: 2
    - id: journal_uuid
      size: 16
    - id: journal_inode_num
      type: u4
    - id: journal_device
      type: u4
    - id: orphan_inodes
      type: u4
enums:
    state_enum:
        1: valid_fs
        2: error_fs
        4: orphans_being_recovered
    errors_enum:
        1: act_continue
        2: act_read_only
        3: act_panic
    creator_os_enum:
        0: linux
        1: hurd
        2: masix
        3: free_bsd
        4: lites
    major_version_enum:
        0: orignial
        1: dynamic