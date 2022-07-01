
from typing import List, Union
# self.v: float = i - sum([x/y for (x, y) in zip(connected_v, connected_r) if type(x) is float])


class Node:
    def __init__(self, name) -> None:
        self.name: str = name
        self.i: float = 0
        self.v = 0
        self.prev_v: float = 0
        self.connected_v: List[Union[Node, float]] = []
        self.connected_r: List[float] = []
        self.via = []
        self.via_r = []
        self.via_v = []

    def add_node(self, node):
        if node not in self.connected_v:
            self.connected_v.append(node)

    def add_via(self, via):
        if via not in self.via:
            self.via.append(via)

    def init(self):
        self.v = self.i - sum([x/y for (x, y) in zip(self.connected_v+self.via_v, self.connected_r+self.via_r) if type(x) is float])
        self.prev_v = self.v

    def calculate(self, e) -> int:
        self.prev_v = self.v
        sum_rr: float = sum([1/x for x in self.connected_r]) + sum([1/x for x in self.via_r])
        sum_v_r: float = 0

        for (x, y) in zip(self.connected_v+self.via_v, self.connected_r+self.via_r):
            if type(x) is Node:
                sum_v_r += x.v/y
            else:
                sum_v_r += x/y

        self.v = sum_v_r/sum_rr-self.i/sum_rr

        if(abs(self.v - self.prev_v) < e):
            return 1
        return 0
