meta:
    id: extent
    endian: le
seq:
    - id: signature
      contents: [0x0a, 0xf3]
    - id: valid_entries_num
      type: u2
    - id: max_entries_num
      type: u2
    - id: depth
      type: u2
    - id: generation
      type: u4
    - id: internal_nodes
      type: internal_node
      repeat: expr
      repeat-expr: 'depth == 0 ? 0 : valid_entries_num'
    - id: leaf_nodes
      type: leaf_node
      repeat: expr
      repeat-expr: 'depth != 0 ? 0 : valid_entries_num'
types:
  internal_node:
    seq:
      - id: block
        type: u4
      - id: leaf_block_ptr_lower
        type: u4
      - id: leaf_block_ptr_higher
        type: u2
      - id: unused
        size: 2
  leaf_node:
    seq:
      - id: first_file_block
        type: u4
      - id: covered_blocks
        type: u2
      - id: first_block_higher
        type: u2
      - id: first_block_lower
        type: u4