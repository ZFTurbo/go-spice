'''
IBM spice -> ngspice formater
'''
import os

files = os.listdir('./data')

for file in files:
    if '.spice' in file:
        new_spice = [f'*{file}\n','rrgnd _X_n_gnd gnd 2.500000e-01\n']
        controls = []

        with open(os.path.join('data',file)) as f:
            old_spice = f.readlines()

            for line in old_spice:
                if '*' not in line:
                    if ('v' in line and 'X' in line) or ('iB' in line):
                        old_line = line.split(' ')
                        old_line.insert(3,'dc')
                        
                        if old_line[1] not in controls:
                            controls.append(old_line[1])
                        
                        if('iB' in line):
                            if('_g' in line and old_line[2]):
                                old_line[1],old_line[2] = old_line[2], '_X_n_gnd'
                                old_line[-2] = '-'+old_line[-2]
                                new_spice.append(' '.join(old_line))
                            if('_v' in line and old_line[1]):
                                old_line[2] = '_X_n_gnd'
                                new_spice.append(' '.join(old_line))
                        else:
                            old_line[2] = '_X_n_gnd'
                            new_spice.append(' '.join(old_line))

                    if ('R' in line) or ('r' in line) or ('V' in line):
                        new_line = 'r' + line[1:]
                        new_spice.append(new_line)

                        if line.split(' ')[1] not in controls:
                            controls.append(line.split(' ')[1])
                        if line.split(' ')[2] not in controls:
                            controls.append(line.split(' ')[2])
            
            new_spice.append('\n.control\n')
            new_spice.append('op\n')

            for control in controls:
                new_spice.append(f"print v({control})\n")

            new_spice.append('.endc\n\n')
            new_spice.append('.end\n')

            with open(f'new_{file}','a') as ff:
                ff.writelines(new_spice)