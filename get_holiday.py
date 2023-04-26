import requests

url = 'http://apis.data.go.kr/B090041/openapi/service/SpcdeInfoService/getRestDeInfo'
params = {'serviceKey': 'kYU0U/7S/1rSI8K9B+0zI6VRUyDuKlImNPSDJAjPlrQiJxeKQ6X25eNH84wDHBX8aEMIfJi9QJ4zTBwQnVWY5w==', 'solYear': '2023','numOfRows': 100}

response = requests.get(url, params=params)
print(response.text)
