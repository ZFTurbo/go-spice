'''
Compare results from IBM solution and from node-model method
'''
from tqdm import tqdm
import numpy as np
from argparse import ArgumentParser
import sys

sys.path.append('../')

if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-m", "--modeling", dest="modelingname", help="Path to sourse file",
                        default='../formated/new_ibmpg0.spice', metavar="FILE")
    parser.add_argument("-s", "--solution", dest="solutionname", help="Path to out file",
                        default='../out/nd_ibmpg0.solution', metavar="FILE")
    args = parser.parse_args()

    with open(args.modelingname) as f:
        lines_model = f.readlines()

    with open(args.solutionname) as f:
        lines_solution = {}
        for line in f.readlines():
            s_l = line.split(' ')
            lines_solution[s_l[0]] = float(s_l[-1])

    pro = []

    print(f'\nSizes: modeling - {len(lines_model)}, solution - {len(lines_solution)}\n')

    with tqdm(total=len(lines_model), bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
        for model in lines_model:
            s_m = model.split(' ')
            one = lines_solution[s_m[0]]
            two = float(s_m[-1])
            
            try:
                if(two > one):
                    pro.append(abs(abs(one-two)/two)*100)
                else:
                    pro.append(abs(abs(one-two)/one)*100)
            except ZeroDivisionError:
                pro.append(0)

            pbar.update(1)
        pbar.close()

    print(f'\nPercentage difference: {np.mean(np.array(pro))}%\n')
