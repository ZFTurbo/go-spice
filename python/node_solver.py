from utils import Node
from typing import List, Union
from tqdm import tqdm
from argparse import ArgumentParser
import time


if __name__ == '__main__':
    start_time = time.time()
    parser = ArgumentParser()
    parser.add_argument("-f", "--file", dest="filename", help="Path to sourse file", metavar="FILE")
    parser.add_argument("-o", "--out", dest="outname", help="Path to out file", metavar="FILE")
    args = parser.parse_args()

    print(f'\nCircuit - {args.filename}')

    with open(args.filename) as f:
        lines = f.readlines()
        e = 1e-6
        max_steps = 10000
        voltage: List[str] = []
        current: List[str] = []
        res: List[str] = []
        name_nodes: List[str] = []
        nodes: List[Node] = []
        temp: List[str] = []

        print('Fetching data - resistance, voltage soruses and current sourses...\n')

        with tqdm(total=len(lines), bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
            for line in lines:
                if('v' == line[0]):
                    voltage.append(line)
                    if(line.split(' ')[1] not in temp):
                        temp.append(line.split(' ')[1])
                elif('i' == line[0]):
                    current.append(line)
                elif('r' == line[0]):
                    if(float(line.split(' ')[-1]) != 0):
                        res.append(line)
                    sp_res = line.split(' ')
                    if(sp_res[1][2:] != sp_res[2][2:]):
                        if(sp_res[1] not in name_nodes):
                            name_nodes.append(sp_res[1])
                        if(sp_res[2] not in name_nodes):
                            name_nodes.append(sp_res[2])
                    else:
                        if((sp_res[1] + ' ' + sp_res[2]) not in name_nodes):
                            name_nodes.append(sp_res[1] + ' ' + sp_res[2])
                            if(sp_res[1] not in temp):
                                temp.append(sp_res[1])
                            if(sp_res[2] not in temp):
                                temp.append(sp_res[2])

                pbar.update(1)
            pbar.close()

        print('\nFiltering parasite nodes...')

        name_nodes = [x for x in name_nodes if x not in temp and 'gnd' not in x]

        temp.clear()

        print(f'Total nodes to be solved: {len(name_nodes)}\n')
        print(f'Creating nodes models...\n')

        with tqdm(total=len(name_nodes), bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
            for node in name_nodes:
                i: float = 0
                node_vol: List[Union[float, str]] = []
                node_res: List[float] = []

                for r in res:
                    s_r = r.split(' ')
                    vol_source = False
                    if(s_r[1] in node):
                        node_res.append(float(s_r[-1]))
                        for v in voltage:
                            if (s_r[2] == v.split(' ')[1]):
                                node_vol.append(float(v.split(' ')[-1]))
                                vol_source = True
                                break
                        if not vol_source and s_r[2] not in node_vol:
                            node_vol.append(s_r[2])
                    elif(s_r[2] in node):
                        node_res.append(float(s_r[-1]))
                        for v in voltage:
                            if (s_r[1] == v.split(' ')[1]):
                                node_vol.append(float(v.split(' ')[-1]))
                                vol_source = True
                                break
                        if not vol_source and s_r[1] not in node_vol:
                            node_vol.append(s_r[1])

                # Find all current soursec for curent node
                for c in current:
                    sp_c = c.split(' ')
                    if(sp_c[1] in node):
                        i += float(sp_c[-2])

                nodes.append(Node(node, i, node_vol, node_res))
                pbar.update(1)
            pbar.close()

        print(f'\nConnecting nodes with each other...\n')

        with tqdm(total=len(name_nodes), bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
            for node_p in nodes:
                for v in range(len(node_p.connected_v)):
                    for node_c in nodes:
                        if(node_c != node_p and type(node_p.connected_v[v]) is str):
                            if (node_p.connected_v[v] in node_c.name):
                                node_p.connected_v[v] = node_c
                pbar.update(1)
            pbar.close()

        print(f'\nSolving nodes...')

        with tqdm(total=max_steps) as pbar:
            for i in range(max_steps):
                solved_nodes = 0
                for node in nodes:
                    solved_nodes += node.calculate(e)
                if(solved_nodes == len(nodes)):
                    print(f"\nSteps count: {i}\n")
                    break
                pbar.update(1)
            pbar.close()

        print(f'Writing results in log...\n')

        with open(args.out, 'w') as w:
            for node in nodes:
                if(len(node.name.split(' ')) != 2):
                    w.write(f'{node.name}  {node.v:e}\n'.lower())
                else:
                    w.write(f"{node.name.split(' ')[0]}  {node.v:e}\n".lower())
                    w.write(f"{node.name.split(' ')[1]}  {node.v:e}\n".lower())

    print("--- %s seconds ---" % (time.time() - start_time))
