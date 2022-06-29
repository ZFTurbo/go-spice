
from typing import List, Union


class Node:
    def __init__(self, name) -> None:
        self.name: str = name
        self.i: float = 0
        # self.v: float = i - sum([x/y for (x, y) in zip(connected_v, connected_r) if type(x) is float])
        self.v = 0
        self.prev_v: float = 0
        self.connected_v: List[Union[Node, float]] = []
        self.connected_r: List[float] = []

    def set(self, v, i, connected_v, connected_r):
        self.v = v
        self.prev_v: float = v
        self.connected_v: List[Union[Node, float]] = connected_v
        self.connected_r: List[float] = connected_r

    def calculate(self, e) -> int:
        self.prev_v = self.v
        sum_rr: float = sum([1/x for x in self.connected_r])
        sum_v_r: float = 0

        for (x, y) in zip(self.connected_v, self.connected_r):
            if type(x) is Node:
                sum_v_r += x.v/y
            else:
                sum_v_r += x/y

        self.v = sum_v_r/sum_rr-self.i/sum_rr

        if(abs(self.v - self.prev_v) < e):
            return 1
        return 0
