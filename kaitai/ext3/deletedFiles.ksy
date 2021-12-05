meta:
    id: deleted_files
    endian: le
seq:
    - id: entries
      type: deleted_entry
      repeat: eos
types:
    deleted_entry:
        seq:
            - id: inode_ptr
              type: u4
            - id: rec_len
              type: u2
            - id: name_len
              type: u1
            - id: file_type
              type: u1
              enum: file_type_enum
            - id: name
              size: name_len
              type: str
              encoding: UTF-8
            - id: padding
              size: '4 - (name_len + 8) % 4 == 4 ? 0 : 4 - (name_len + 8) % 4'
        enums:
            file_type_enum:
                0: unknown
                1: reg_file
                2: dir
                3: chrdev
                4: blkdev
                5: fifo
                6: sock
                7: symlink