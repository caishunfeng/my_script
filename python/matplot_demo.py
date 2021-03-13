# 折线图
# from matplotlib import pyplot as plt
# x = [5,2,7]
# y = [2,16,4]
# plt.plot(x,y)
# plt.title('Image Title')
# plt.ylabel('Y axis')
# plt.xlabel('X axis')
# plt.show()

# from matplotlib import pyplot as plt
# from matplotlib import style
# style.use('ggplot')
# x = [5,8,10]
# y = [12,16,6]
# x2 = [6,9,11]
# y2 = [6,15,7]
# plt.plot(x,y,'g',label='line one', linewidth=5)
# plt.plot(x2,y2,'r',label='line two',linewidth=5)
# plt.title('Epic Info')
# plt.ylabel('Y axis')
# plt.xlabel('X axis')
# # 设置图例位置
# plt.legend()
# plt.grid(True,color='k')
# plt.show()

# 条形图
from matplotlib import pyplot as plt
plt.bar([0.25,1.25,2.25,3.25,4.25],[50,40,70,80,20],label="BMW", color='b', width=.5)
plt.bar([.75,1.75,2.75,3.75,4.75],[80,20,20,50,60],label="Audi", color='r',width=.5)
plt.legend()
plt.xlabel('Days')
plt.ylabel('Distance (kms)')
plt.title('Information')
plt.show()

# 直方图
# import matplotlib.pyplot as plt
# population_age = [22,55,62,45,21,22,34,42,42,4,2,102,95,85,55,110,120,70,65,55,111,115,80,75,65,54,44,43,42,48]
# bins = [0,10,20,30,40,50,60,70,80,90,100]
# plt.hist(population_age, bins, histtype='bar', color='b', rwidth=0.8)
# plt.xlabel('age groups')
# plt.ylabel('Number of people')
# plt.title('Histogram')
# plt.show()

# 散点图
# import matplotlib.pyplot as plt
# x = [1,1.5,2,2.5,3,3.5,3.6]
# y = [7.5,8,8.5,9,9.5,10,10.5]
# x1=[8,8.5,9,9.5,10,10.5,11]
# y1=[3,3.5,3.7,4,4.5,5,5.2]
# plt.scatter(x,y, label='high income low saving',color='r')
# plt.scatter(x1,y1,label='low income high savings',color='b')
# plt.xlabel('saving*100')
# plt.ylabel('income*1000')
# plt.title('Scatter Plot')
# plt.legend()
# plt.show()

# 面积图
# import matplotlib.pyplot as plt
# days = [1,2,3,4,5]
# sleeping =[7,8,6,11,7]
# eating = [2,3,4,3,2]
# working =[7,8,7,2,2]
# playing = [8,5,7,8,13]
# plt.plot([],[],color='m', label='Sleeping', linewidth=5)
# plt.plot([],[],color='c', label='Eating', linewidth=5)
# plt.plot([],[],color='r', label='Working', linewidth=5)
# plt.plot([],[],color='k', label='Playing', linewidth=5)
# plt.stackplot(days, sleeping,eating,working,playing, colors=['m','c','r','k'])
# plt.xlabel('x')
# plt.ylabel('y')
# plt.title('Stack Plot')
# plt.legend()
# plt.show()

# 饼图
# import matplotlib.pyplot as plt
# days = [1,2,3,4,5]
# sleeping =[7,8,6,11,7]
# eating = [2,3,4,3,2]
# working =[7,8,7,2,2]
# playing = [8,5,7,8,13]
# slices = [7,2,2,13]
# activities = ['sleeping','eating','working','playing']
# cols = ['c','m','r','b']
# plt.pie(slices,  labels=activities,  colors=cols,  startangle=90,  shadow= True,  explode=(0,0,0.1,0),  autopct='%1.1f%%')
# plt.title('Pie Plot')
# plt.show()

# 多图合并
# import numpy as np
# import matplotlib.pyplot as plt
# def f(t):
#     return np.exp(-t) * np.cos(2*np.pi*t)
# t1 = np.arange(0.0, 5.0, 0.1)
# t2 = np.arange(0.0, 5.0, 0.02)
# plt.subplot(221)
# plt.plot(t1, f(t1), 'bo', t2, f(t2))
# plt.subplot(222)
# plt.plot(t2, np.cos(2*np.pi*t2))
# plt.show()

# 双y轴绘制 关键函数：twinx()
# -*- coding: utf-8 -*-
# import numpy as np
# import matplotlib.pyplot as plt
# from matplotlib import rc

# rc('mathtext', default='regular')

# time = np.arange(10)
# temp = np.random.random(10) * 30
# Swdown = np.random.random(10) * 100 - 10
# Rn = np.random.random(10) * 100 - 10

# fig = plt.figure()
# ax = fig.add_subplot(111)
# ax.plot(time, Swdown, '-', label='Swdown')
# ax.plot(time, Rn, '-', label='Rn')
# ax2 = ax.twinx()
# ax2.plot(time, temp, '-r', label='temp')
# ax.legend(loc=2)
# ax.grid()
# ax.set_xlabel("Time (h)")
# ax.set_ylabel(r"Radiation ($MJ\,m^{-2}\,d^{-1}$)")
# ax2.set_ylabel(r"Temperature ($^\circ$C)")
# ax2.set_ylim(0, 35)
# ax.set_ylim(-20, 100)
# ax2.legend(loc=0)
# plt.savefig('0.png')
# plt.show()



# import numpy as np
# import matplotlib.pyplot as plt

# X = [1, 2, 3, 4, 5, 6]
# Y = [2.6, 3.3, 4.7, 5.2, 6.8, 7.1]
# # Y=[2,3,4,5,6,7]
# # 用一次多项式拟合，相当于线性拟合
# z1 = np.polyfit(X, Y, 1)
# p1 = np.poly1d(z1)
# print(z1)
# print(p1)

# x = np.arange(1, 7)
# y = z1[0] * x + z1[1]
# plt.figure()
# plt.scatter(X, Y)
# plt.plot(x, y)
# plt.show()
