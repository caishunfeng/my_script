import csv

dtsTables = open('C:\\Users\\cysp\\desktop\\表A.csv', 'r')
allExistTables = open('C:\\Users\\cysp\\desktop\\表B.csv', 'r')
dtsTablesReader = csv.DictReader(dtsTables)
dtsTablesHeader = next(dtsTablesReader)

allExistTablesReader = csv.DictReader(allExistTables)
allExistTablesHeader = next(allExistTablesReader)

dtsTableMap = {}
allExistTableMap = {}
for row in dtsTablesReader:
    dtsTableMap[row['表A']] = 1

for row in allExistTablesReader:
    allExistTableMap[row['表B']] = 1

print(dtsTableMap.keys())
print(allExistTableMap.keys())

for table in allExistTableMap:
    if(table not in dtsTableMap.keys()):
        print(table)
