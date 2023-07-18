import os
from tqdm import tqdm
import platform
import datetime
import subprocess
import json

linux_bin = "./separa"
win_bin = "./separa.exe"

def read_target(file: str) -> list[str]:
    res: list[str] = []
    with open(file, 'r') as f:
        for line in f.readlines():
            res.append(line.strip())
    return res


def merge_target(dir_path: str):
    res = {}
    for filename in os.listdir(dir_path):
        filepath = os.path.join(dir_path, filename)
        if os.path.isfile(filepath):
            with open(filepath, 'r', encoding='utf-8') as f:
                data = json.load(f)
                res.update(data)
    json_str = json.dumps(res, default=lambda x: x if x is not None else 'null', indent=4, ensure_ascii=False)
    with open(dir_path + ".json", 'w', encoding='utf-8') as f:
        f.write(json_str)


if __name__ == '__main__':
    system = platform.system()
    folder_name = datetime.datetime.now().strftime("%Y%m%d%H%M")
    bin_path = ""

    targets = read_target('target.txt')
    print("load " + str(len(targets)) + " CIDR to scan")
    for target in tqdm(targets):
        output = os.path.join('outputs', folder_name, target[:-3] + ".json")

        if system == 'Windows':
            bin_path = win_bin
        elif system == 'Linux':
            bin_path = linux_bin
        command = [bin_path, "scan", "-t", target, "-o", output, "-d", "5", "-n", "800"]
        subprocess.run(command)

    print("merge result")
    merge_target("outputs/" + folder_name)