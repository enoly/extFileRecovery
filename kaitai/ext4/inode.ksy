meta:
    id: inode
    endian: le
seq:
    - id: mode
      type: u2
    - id: uid
      type: u2
    - id: size
      type: u4
    - id: atime
      type: u4
    - id: ctime
      type: u4
    - id: mtime
      type: u4
    - id: dtime
      type: u4
    - id: gid
      type: u2
    - id: links_count
      type: u2
    - id: blocks
      type: u4
    - id: flags
      type: u4
    - id: osd1
      type: u4
    - id: i_block
      size: 60
    - id: generation
      type: u4
    - id: file_acl
      type: u4
    - id: dir_acl
      type: u4
    - id: faddr
      type: u4
    - id: osd2
      size: 12