import os
from tqdm import tqdm
import platform
import datetime
import subprocess

linux_bin = "./separa"
win_bin = "./separa.exe"

def read_target(file: str) -> list[str]:
    res: list[str] = []
    with open(file, 'r') as f:
        for line in f.readlines():
            res.append(line.strip())
    return res


def scan_target():
    pass


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
        subprocess.run(command, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)