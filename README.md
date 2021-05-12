# utils

工具包

###  bar 进度条
    一个简单的进度条

###  buffer 
    byte缓存区，支持读写偏移。

###  heap
    堆，使用 container/heap 的封装

###  log 异步日志
    异步日志工具，
    支持时间分割，每天切分日志文件。 
    支持文件大小分割，达到日志存储上限，切分日志。
    支持日志等级划分，异步输出。

###  orderMap 
    有序的map

###  pipeline 
    流水线 ， step

###  queue
    队列
    channelQueue  blockQueue  

###  timer
    定时器
    heapTimer 小根堆定时器，高精度 timer。调用系统timer，提供统一管理
    
    timingWheel 时间轮定时器，低精度 timer, 最低精度为毫秒。系统ticker驱动。
    如果时间轮精度为10ms， 那么他的误差在 （0，10）ms之间。如果一个任务延迟 500ms，那它的执行时间在490～500ms之间。
    按平均来讲，出错的概率均等的情况下，那么这个出错可能会延迟或提前最小刻度的一半，在这里就是10ms/2=5ms.
    故，时间轮的 tick 单位 在总延迟时间上，应该不足以影响延迟执行函数处理的事务。
