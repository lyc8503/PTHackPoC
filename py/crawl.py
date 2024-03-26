# 通过 RSS 爬所有种子信息
import json
import requests
import lxml.etree


def parse_xml(xml_text):
    root = lxml.etree.fromstring(xml_text.replace(b"&nbsp;", b'\n'), parser = lxml.etree.XMLParser(strip_cdata=False))

    ret = []

    for i in root.findall('.//item'):
        tmp = {}
        for child in i:
            tmp[child.tag] = child.text
        ret.append(tmp)

    return ret


def req(start, size):
    r = requests.get(f"https://byr.pt/torrentrss.php?rows={size}&startindex={start}&icat=1&ismalldescr=1&isize=1&iuplder=1&passkey=21d878ce65a3623ee988cbefc329bbc9", timeout=5)
    return parse_xml(r.content)


for i in range(0, 70000, 50):
    print(i)
    for _ in range(0, 10):
        try:
            r = req(i, 50)
            break
        except Exception as e:
            print(e)
    for item in r:
        with open("byr.txt", "a") as f:
            f.write(json.dumps(item, ensure_ascii=False) + "\n")
