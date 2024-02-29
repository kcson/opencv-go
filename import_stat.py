import csv
import json

import pymysql
import requests

year = '2023'
url = 'http://apis.data.go.kr/B090041/openapi/service/SpcdeInfoService/getRestDeInfo'
params = {'serviceKey': 'kYU0U/7S/1rSI8K9B+0zI6VRUyDuKlImNPSDJAjPlrQiJxeKQ6X25eNH84wDHBX8aEMIfJi9QJ4zTBwQnVWY5w==', 'solYear': year, 'numOfRows': 100, '_type': 'json'}

# host = 'parkingcone.ckxv9kglgkje.ap-northeast-2.rds.amazonaws.com' #dev
host = 'itcha-db-az.ckxv9kglgkje.ap-northeast-2.rds.amazonaws.com'  # real
db = 'parkingcone'
user = 'parkingcone'
password = 'last30min'
conn = None
query = '''INSERT INTO B2B_SPOT_USER_COUNT_STAT
           SET
                spot_idx = %s,
                user_id = %s,
                use_count = %s,
                zone_idx = %s
           ON DUPLICATE KEY UPDATE
                use_count = use_count+%s'''

delete_query = 'DELETE FROM HOLIDAY_T WHERE holiday_year = %s'
try:
    # DB connection
    conn = pymysql.connect(host=host, user=user, password=password, db=db, charset='utf8')
    cursor = conn.cursor()

    data_count = 0
    with open('엘지20231011파킹페이퍼고객명단.csv', newline='\r\n') as csvfile:
        csvreader = csv.reader(csvfile, delimiter=',', quotechar='|')
        for row in csvreader:
            cursor.execute(query, (222, row[0], row[1], 124, row[1]))
            data_count += 1
            print(data_count)
        conn.commit()
except Exception as e:
    if conn is not None:
        conn.rollback()
    print(e)
finally:
    if conn is not None:
        conn.close()
