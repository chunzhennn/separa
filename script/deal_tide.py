import os
from tqdm import tqdm
import platform
import datetime
import subprocess
import re

def read_file(file: str) -> list[str]:
    res: list[str] = []
    with open(file, 'r', encoding="utf-8") as f:
        for line in f.readlines():
            res.append(line.strip())
    return res

def read_all_file(file: str) -> str:
    res = ""
    with open(file, 'r', encoding="utf-8") as f:
        res = f.read()
    return res

if __name__ == '__main__':

    fingers = read_file('./static/tide.txt')
    print("load " + str(len(fingers)) + " fingers")
    
    # result = {}
    # for finger in fingers:
    #     product, point = finger.split('\t', 1)
    #     points = point.split('||')
    #     result[product] = []
    #     for p in points:
    #         result[product].append(p.strip())
    
    inpu = read_all_file('./static/tide.txt')
    # print(inpu)
    out = open("./static/tide_fixed.txt", "w", encoding="utf-8")
    rs = ""
    p = re.compile(r'(?<=[Body|Title|Header]\=\")(?:(?!\|\||&&).)*(?=\")')
    for finger in fingers:

        res = p.findall(finger)

        for single in res:
            # print(single)
            if '"' in single:
                single_new = single.replace('"', '\\"')
            #     print(single_new)
                finger = finger.replace(single, single_new)
        rs += finger + '\n'
    out.write(rs.rstrip('\n'))
    out.close()