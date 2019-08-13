package worker

//调度 , 我理解为分配者 , 从JobMgr处获得Jobs , 然后按一定规则分配给executor
//类似产品经理的职责 , 划分好时间(下次调度时间) , 记录状态 , 汇报状态(日志to mongodb) , 验收成果
