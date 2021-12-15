meta:
    id: gdt_record
    endian: le
seq:
    - id: block_bitmap_block
      type: u4
    - id: inode_bitmap_block
      type: u4
    - id: inode_table_block
      type: u4
    - id: free_blocks_count
      type: u2
    - id: free_inodes_count
      type: u2
    - id: used_dirs_count
      type: u2
    - id: bg_flags
      type: u2