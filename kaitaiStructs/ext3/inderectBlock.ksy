meta:
  id: indirect_block
  endian: le
seq:
  - id: blocks_ptrs
    type: block_ptr
    repeat: eos
types:  
    block_ptr:
        seq:
            - id: ptr
              type: u4