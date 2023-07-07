

def read_file(file: str) -> list[str]:
    res: list[str] = []
    with open(file, 'r', encoding="utf-8") as f:
        for line in f.readlines():
            res.append(line.strip())
    return res


if __name__ == '__main__':

    fingers_a = read_file('./static/fingerprint.txt')
    fingers_b = read_file('./static/tide_fixed.txt')

    fing: dict[str, list[str]] = {}

    print("load " + str(len(fingers_a)) + " fingers_a")
    print("load " + str(len(fingers_b)) + " fingers_b")

    for finger in fingers_a:
        product, point = finger.split('\t', 1)
        if product in fing:
            fing[product].append(point)
        else:
            fing[product] = [point]

    for finger in fingers_b:
        product, point = finger.split('\t', 1)
        if product in fing:
            if product not in fing[product]:
                fing[product].append(point)
        else:
            fing[product] = [point]
    
    cnt = 0
    out = open("./static/merge.txt", "w+", encoding="utf-8")
    for product, points in fing.items():
        for point in points:
            cnt += 1
            out.write(product + '\t' + point)
            if cnt != 31260:
                out.write('\n')
    out.close()
    print("total cnt: " + str(cnt))
    