# !/usr/bin/env python
# -*- encoding: utf-8 -*-

import requests
import json

# 将 json 文件读取成字符串
# json_data = open('./data.json').read()
url="请求链接"

r = requests.get(url)
print(r.status_code)

req_data = r.content.decode(encoding='utf-8')
# print(req_data)

req_json = json.loads(req_data)
# print(req_data)
# print(req_json['query_result']['data']['rows'])

print(type(req_json))

rows = req_json['query_result']['data']['rows']
for row in rows:
    print(row)