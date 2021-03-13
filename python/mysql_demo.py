# !/usr/bin/env python
# -*- encoding: utf-8 -*-
# Open database connection
import pymysql
import random

db_url = "数据库连接"
db = pymysql.connect(db_url, "账号", "密码", "数据库")
cursor = db.cursor()

# 指定2年级的数据生成
sql = 'SELECT cmi.id,cmi.master_id,cmi.clazz_id ' \
      'from clazz_master_info cmi join clazz c on c.id = cmi.clazz_id ' \
      'where c.year = 2020 and c.type = 2 and c.time_seq = 0 and c.status != -1 and cmi.master_id !=0 ' \
      # 'and c.grade = "2A,2B" '

cursor.execute(sql)
rows = cursor.fetchall()
# print(rows)

for row in rows:
    rowId = row[0]
    masterId = row[1]
    clazzId = row[2]
    print(rowId, masterId, clazzId)

    serviceNum = random.randint(100, 200)
    sameSubjectConversionClazzs = random.randint(0, 150)
    sumExpandClazzs = random.randint(0, 150)
    quitClazzs = random.randint(0, 10)

    updateSql = 'UPDATE clazz_master_info ' \
          'SET service_num = {},same_subject_conversion_clazzs = {},sum_expand_clazzs = {},quit_clazz_count = {} ' \
          'WHERE id = {}'.format(serviceNum, sameSubjectConversionClazzs, sumExpandClazzs, quitClazzs, rowId)

    # insertSql = 'insert into clazz_master_info(clazz_id,master_id,period_id,service_num,same_subject_conversion_clazzs,sum_expand_clazzs,quit_clazz_count,status) ' \
    #       'values ({},{},{},{},{},{},{},{})' \
    #       .format(clazzId, masterId, 2, serviceNum, sameSubjectConversionClazzs, sumExpandClazzs, quitClazzs, 1)

    try:
        cursor.execute(updateSql)
        # cursor.execute(insertSql)
        db.commit()
    except Exception:
        print("发生异常", Exception)
        db.rollback()

# disconnect from server
db.close()
