'''
Compare results from IBM solution and from node-model method
'''
from tqdm import tqdm
import numpy as np

with open('./data/ibmpg1.solution') as f:
    lines_solution = f.readlines()

with open('./nd_ibmpg1.solution') as f:
    lines_model = f.readlines()

pro = []

with tqdm(total=len(lines_model)*len(lines_solution), bar_format='{l_bar}{bar:20}{r_bar}{bar:-20b}') as pbar:
    for model in lines_model:
        for solution in lines_solution:
            if solution.split(' ')[0] == model.split(' ')[0]:
                one = (float(solution.split(' ')[-1]))
                two = (float(model.split(' ')[-1]))
                try:
                    if(two>one):
                        pro.append(abs(abs(one-two)/two)*100)
                    else:
                        pro.append(abs(abs(one-two)/one)*100)
                except ZeroDivisionError:
                    pro.append(0)

            pbar.update(1)
    pbar.close()

print(f'\nPercentage difference: {np.mean(np.array(pro))}%')
