import json

import requests

year = '2023'
url = 'http://apis.data.go.kr/B090041/openapi/service/SpcdeInfoService/getRestDeInfo'
params = {'serviceKey': 'kYU0U/7S/1rSI8K9B+0zI6VRUyDuKlImNPSDJAjPlrQiJxeKQ6X25eNH84wDHBX8aEMIfJi9QJ4zTBwQnVWY5w==', 'solYear': year, 'numOfRows': 100, '_type': 'json'}

response = requests.get(url, params=params)
print(response.text)

jsonObj = json.loads(response.text)
holiDays = jsonObj['response']['body']['items']['item']
for holiDay in holiDays:
    print(holiDay['dateName'])
    print(holiDay['locdate'])
    print(holiDay['isHoliday'])
    print('================================')
