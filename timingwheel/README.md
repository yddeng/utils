# timingwheel

分级时间轮

一开始，第一层的时间轮所能表示时间范围是0~20Ms之间，假设现在出现一个任务的延迟时间是200Ms，再创建一层时间轮，我们称之为第二层时间轮。

第二层时间轮的创建代码如下

```
//超过当前时间轮的总量，创建更大的时间轮
	if int(delayTime/tw.tickDur) > tw.wheelSize {
		levUpWheel := tw.levUpWheel
		if levUpWheel == nil {
			newTW := newTimingWheel(tw.tickDur*time.Duration(tw.wheelSize), tw.wheelSize, tw)
			tw.levUpWheel = newTW
			levUpWheel = tw.levUpWheel
		}
		levUpWheel.addTimer(delayTime, timer)

	} else {
		tw.addTimer(delayTime, timer)

	}
```

也就是第二层时间轮每一个槽所能表示的时间是第一层时间轮所能表示的时间范围，也就是20Ms。槽的数量还是一样，其他的属性也是继承自第一层时间轮。这时第二层时间轮所能表示的时间范围就是0~400Ms了。

同理，如果第二层时间轮的时间范围还容纳不了新的延迟任务，就会创建第三层、第四层...

值得注意的是，只有当前时间轮无法容纳目标延迟任务所能表示的时间时，才需要创建更高一级的时间轮，或者说把该任务加到更高一级的时间轮中(如果该时间轮已创建)。


创建更大时间轮时间，将大轮的事务移交到小轮上存在误差，还需改进。最好一个只用时间轮。
