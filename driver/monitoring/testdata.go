package monitoring

import (
	"fmt"
	"time"
)

var (
	Node1TestId = Node("A")
	Node2TestId = Node("B")
	Node3TestId = Node("C")

	Node1TestLog = "INFO [05-04|09:34:15.080] New block      index=1 id=2:1:247c79       gas_used=11 txs=10/0 age=7.392s t=\"711.334µs\" \n" +
		"INFO [05-04|09:34:15.537] New block      index=2 id=3:1:3d6fb6       gas_used=22 txs=20/0 age=343.255ms t=1.579ms \n" +
		"INFO [05-04|09:34:16.027] New block      index=3 id=3:4:9bb789       gas_used=33   txs=30/0 age=380.470ms t=1.540ms \n"

	Node2TestLog = "INFO [05-04|09:34:16.512] New block      index=1 id=2:1:247c79       gas_used=11 txs=10/0 age=7.392s t=4.686ms \n" +
		"INFO [05-04|09:34:17.003] New block      index=2 id=3:1:3d6fb6       gas_used=22 txs=20/0 age=343.255ms t=2.579ms \n"

	Node3TestLog = "INFO [05-04|09:38:15.080] New block      index=1 id=2:1:247c79       gas_used=11 txs=10/0 age=7.392s t=5.686ms \n"

	year     = time.Now().Year()
	time1, _ = time.Parse("2006-[01-02|15:04:05.000]", fmt.Sprintf("%d-[05-04|09:34:15.080]", year))
	time2, _ = time.Parse("2006-[01-02|15:04:05.000]", fmt.Sprintf("%d-[05-04|09:34:15.537]", year))
	time3, _ = time.Parse("2006-[01-02|15:04:05.000]", fmt.Sprintf("%d-[05-04|09:34:16.027]", year))
	time4, _ = time.Parse("2006-[01-02|15:04:05.000]", fmt.Sprintf("%d-[05-04|09:34:16.512]", year))
	time5, _ = time.Parse("2006-[01-02|15:04:05.000]", fmt.Sprintf("%d-[05-04|09:34:17.003]", year))
	time6, _ = time.Parse("2006-[01-02|15:04:05.000]", fmt.Sprintf("%d-[05-04|09:38:15.080]", year))

	dur1, _ = time.ParseDuration("711.334µs")
	dur2, _ = time.ParseDuration("1.579ms")
	dur3, _ = time.ParseDuration("1.540ms")
	dur4, _ = time.ParseDuration("4.686ms")
	dur5, _ = time.ParseDuration("2.579ms")
	dur6, _ = time.ParseDuration("5.686ms")

	block1 = Block{Height: 1, Time: time1, Txs: 10, GasUsed: 11, ProcessingTime: dur1}
	block2 = Block{Height: 2, Time: time2, Txs: 20, GasUsed: 22, ProcessingTime: dur2}
	block3 = Block{Height: 3, Time: time3, Txs: 30, GasUsed: 33, ProcessingTime: dur3}

	block4 = Block{Height: 1, Time: time4, Txs: 10, GasUsed: 11, ProcessingTime: dur4}
	block5 = Block{Height: 2, Time: time5, Txs: 20, GasUsed: 22, ProcessingTime: dur5}

	block6 = Block{Height: 1, Time: time6, Txs: 10, GasUsed: 11, ProcessingTime: dur6}

	NodeBlockTestData = map[Node][]Block{
		Node1TestId: {block1, block2, block3},
		Node2TestId: {block4, block5},
		Node3TestId: {block6},
	}

	BlockHeight1TestMap = map[int]Block{
		1: block1,
		2: block2,
		3: block3,
	}

	BlockchainTestData = []Block{block1, block2, block3}
)
