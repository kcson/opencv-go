import json

import pymysql
import requests

year = '2025'
url = 'http://apis.data.go.kr/B090041/openapi/service/SpcdeInfoService/getRestDeInfo'
params = {'serviceKey': 'kYU0U/7S/1rSI8K9B+0zI6VRUyDuKlImNPSDJAjPlrQiJxeKQ6X25eNH84wDHBX8aEMIfJi9QJ4zTBwQnVWY5w==', 'solYear': year, 'numOfRows': 100, '_type': 'json'}

host = 'parkingcone.ckxv9kglgkje.ap-northeast-2.rds.amazonaws.com'
db = 'parkingcone-test'
user = 'parkingcone'
password = 'last30min'
conn = None
query = '''INSERT INTO HOLIDAY_T
           SET
                holiday_year = %s,
                date_name = %s,
                is_holiday = %s,
                holiday_date = %s'''
delete_query = 'DELETE FROM HOLIDAY_T WHERE holiday_year = %s'
try:
    # DB connection
    conn = pymysql.connect(host=host, user=user, password=password, db=db, charset='utf8')
    cursor = conn.cursor()

    response = requests.get(url, params=params)
    print(response.text)

    jsonObj = json.loads(response.text)
    holiDays = jsonObj['response']['body']['items']['item']
    cursor.execute(delete_query, year)
    for holiDay in holiDays:
        print(holiDay['dateName'])
        print(holiDay['locdate'])
        print(holiDay['isHoliday'])
        print('================================')
        cursor.execute(query, (year, holiDay['dateName'], holiDay['isHoliday'], holiDay['locdate']))

    conn.commit()
except Exception as e:
    if conn is not None:
        conn.rollback()
    print(e)
finally:
    if conn is not None:
        conn.close()
