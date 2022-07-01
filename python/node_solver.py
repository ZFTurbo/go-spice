from utils import Node
from typing import Dict
from tqdm import tqdm
from argparse import ArgumentParser
import time
import os
import sys

sys.path.append('../')

if __name__ == '__main__':
    start_time = time.time()
    parser = ArgumentParser()
    parser.add_argument("-f", "--file", dest="filename", help="Path to sourse file",
                        default='../formated/new_ibmpg0.spice', metavar="FILE")
    parser.add_argument("-o", "--out", dest="outname", help="Path to out file",
                        default='../out/nd_ibmpg0.solution', metavar="FILE")
    args = parser.parse_args()

    print(f'\nCircuit - {args.filename}')

    with open(args.filename) as f:
        lines = f.readlines()
        # Modeling parameters
        e = 1e-8
        max_steps = 25000
        voltage: Dict = {}
        current: Dict = {}
        node_json: Dict = {}

        print('Fetching data - resistance, voltage soruses and current sourses...\n')

        # Fetchin data from spice net list
        with tqdm(total=len(lines), bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
            for line in lines:
                sp_l = line.split(' ')
                # Collect voltage line, also add to temp for further node filter
                if('v' == line[0]):
                    voltage[sp_l[1]] = float(sp_l[-1])
                # Collect current sources lines
                elif('i' == line[0]):
                    if(sp_l[1] not in current.keys()):
                        current[sp_l[1]] = float(sp_l[-2])
                    else:
                        current[sp_l[1]] += float(sp_l[-2])
                # Collect resistors lines
                elif('r' == line[0] and 'gnd' not in line):
                    # If resistor is not via then just append end and start of resistor if not included
                    if(float(sp_l[-1]) != 0):
                        # res.append(line)
                        if(sp_l[1] not in node_json.keys()):
                            node_json[sp_l[1]] = Node(sp_l[1])
                            node_json[sp_l[1]].add_node(sp_l[2])
                            node_json[sp_l[1]].connected_r.append(float(sp_l[-1]))
                        else:
                            node_json[sp_l[1]].add_node(sp_l[2])
                            node_json[sp_l[1]].connected_r.append(float(sp_l[-1]))

                        if(sp_l[2] not in node_json.keys()):
                            node_json[sp_l[2]] = Node(sp_l[2])
                            node_json[sp_l[2]].add_node(sp_l[1])
                            node_json[sp_l[2]].connected_r.append(float(sp_l[-1]))
                        else:
                            node_json[sp_l[2]].add_node(sp_l[1])
                            node_json[sp_l[2]].connected_r.append(float(sp_l[-1]))
                    else:
                        # Resistor is via, so merge end and start in onde node and add if not included
                        if(sp_l[1] not in node_json.keys()):
                            node_json[sp_l[1]] = Node(sp_l[1])
                            node_json[sp_l[1]].add_via(sp_l[2])
                        else:
                            node_json[sp_l[1]].add_via(sp_l[2])

                        if(sp_l[2] not in node_json.keys()):
                            node_json[sp_l[2]] = Node(sp_l[2])
                            node_json[sp_l[2]].add_via(sp_l[1])
                        else:
                            node_json[sp_l[2]].add_via(sp_l[1])

                pbar.update(1)
            pbar.close()

        print(f'Total nodes to be solved: {len(node_json)}\n')
        print(f'Creating nodes model...\n')

        # Create connection between all nodes
        with tqdm(total=len(node_json), bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
            for node in node_json:
                if(len(node_json[node].via) == 0):
                    if current.get(node) is not None:
                        node_json[node].i += current[node]

                    for idx, v in enumerate(node_json[node].connected_v):
                        if voltage.get(v) is None:
                            node_json[node].connected_v[idx] = node_json.get(v)
                        else:
                            node_json[node].connected_v[idx] = voltage.get(v)
                else:
                    if current.get(node) is not None:
                        node_json[node].i += current[node]

                    for via in node_json[node].via:
                        node_json[node].via_v += [x for x in node_json[via].connected_v]
                        node_json[node].via_r += node_json[via].connected_r

                        if current.get(via) is not None:
                            node_json[node].i += current[via]

                    for idx, v in enumerate(node_json[node].connected_v):
                        if voltage.get(v) is None:
                            if(node_json.get(v) is not None):
                                node_json[node].connected_v[idx] = node_json.get(v)
                        else:
                            node_json[node].connected_v[idx] = voltage.get(v)

                    for idx, v in enumerate(node_json[node].via_v):
                        if voltage.get(v) is None:
                            if(node_json.get(v) is not None):
                                node_json[node].via_v[idx] = node_json.get(v)
                        else:
                            node_json[node].via_v[idx] = voltage.get(v)

                pbar.update(1)
            pbar.close()

        print(f'\nInitialization of values...\n')

        with tqdm(total=len(node_json), bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
            for node in node_json:
                node_json[node].init()
                pbar.update(1)
            pbar.close()

        print(f'\nSolving nodes...')

        # Solving system of equations with Zeideil method
        with tqdm(total=max_steps, bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
            for i in range(max_steps):
                solved_nodes = 0
                for node in node_json:
                    solved_nodes += node_json[node].calculate(e)
                if(solved_nodes == len(node_json)):
                    break
                pbar.update(1)
            pbar.close()

        print(f"\nSteps count: {i}\n")
        print(f'Writing results in log...\n')

        # Write solution in file
        if(not os.path.exists('../out')):
            os.mkdir('../out')

        with open(args.outname, 'w') as w:
            for node in node_json:
                w.write(f'{node_json[node].name}  {node_json[node].v:e}\n')

    print("--- %s seconds ---" % (time.time() - start_time))
