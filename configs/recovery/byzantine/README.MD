# Desgin for byzantine simulation

## Simulation setting

- four nodes: node0, node1, node2(in group1), node2'(in group2)
- node0 and node1 are byzantine, they run two instances respectively in group1 and group2 but use the same identity
- In the views of node2 and node2', they thing there are four nodes participating in the consensus but one losts its connection
