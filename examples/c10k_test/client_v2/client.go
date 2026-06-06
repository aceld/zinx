package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	CmdCalculate = 1001
	CmdSmall     = 2001
	CmdMedium    = 3001
	CmdLarge     = 4001
)

type Stats struct {
	success     int64
	failed      int64
	respTimes   []float64
	respMu      sync.Mutex
	histBuckets []int64
	histMu      sync.RWMutex
}

const (
	histBucketCount = 20
)

func NewStats() *Stats {
	return &Stats{
		respTimes:   make([]float64, 0, 100000),
		histBuckets: make([]int64, histBucketCount),
	}
}

func (s *Stats) addSuccess(n int64) {
	s.respMu.Lock()
	s.success += n
	s.respMu.Unlock()
}

func (s *Stats) addFailed(n int64) {
	s.respMu.Lock()
	s.failed += n
	s.respMu.Unlock()
}

func (s *Stats) addRespTime(ms float64) {
	s.respMu.Lock()
	s.respTimes = append(s.respTimes, ms)
	s.respMu.Unlock()

	bucket := int(ms / 50)
	if bucket >= histBucketCount {
		bucket = histBucketCount - 1
	}
	s.histMu.Lock()
	s.histBuckets[bucket]++
	s.histMu.Unlock()
}

func (s *Stats) getPercentile(p float64) float64 {
	s.respMu.Lock()
	defer s.respMu.Unlock()

	if len(s.respTimes) == 0 {
		return 0
	}

	index := int(math.Ceil(float64(len(s.respTimes))*p)) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(s.respTimes) {
		index = len(s.respTimes) - 1
	}

	sorted := make([]float64, len(s.respTimes))
	copy(sorted, s.respTimes)
	sort.Float64s(sorted)
	return sorted[index]
}

func (s *Stats) getAvg() float64 {
	s.respMu.Lock()
	defer s.respMu.Unlock()

	if len(s.respTimes) == 0 {
		return 0
	}

	var sum float64
	for _, t := range s.respTimes {
		sum += t
	}
	return sum / float64(len(s.respTimes))
}

func (s *Stats) printHistogram() {
	s.histMu.RLock()
	defer s.histMu.RUnlock()

	fmt.Println("\n响应时间分布 (ms):")
	fmt.Println(strings.Repeat("-", 50))

	var total int64
	for _, v := range s.histBuckets {
		total += v
	}

	maxBucket := int64(0)
	for _, v := range s.histBuckets {
		if v > maxBucket {
			maxBucket = v
		}
	}

	for i := 0; i < histBucketCount; i++ {
		low := i * 50
		high := (i + 1) * 50
		count := s.histBuckets[i]
		percentage := float64(count) / float64(total) * 100

		barLen := 0
		if maxBucket > 0 {
			barLen = int(float64(count) / float64(maxBucket) * 40)
		}
		bar := strings.Repeat("█", barLen)

		fmt.Printf("%4d-%4dms: %8d (%5.2f%%) %s\n", low, high, count, percentage, bar)
	}
	fmt.Println(strings.Repeat("-", 50))
}

func sendMsg(conn net.Conn, msgID uint32, data []byte) error {
	header := make([]byte, 8)
	binary.BigEndian.PutUint32(header[0:4], msgID)
	binary.BigEndian.PutUint32(header[4:8], uint32(len(data)))
	_, err := conn.Write(append(header, data...))
	return err
}

func recvMsg(conn net.Conn) (uint32, []byte, error) {
	header := make([]byte, 8)
	if _, err := conn.Read(header); err != nil {
		return 0, nil, err
	}
	msgID := binary.BigEndian.Uint32(header[0:4])
	dataLen := binary.BigEndian.Uint32(header[4:8])
	body := make([]byte, dataLen)
	if _, err := conn.Read(body); err != nil {
		return 0, nil, err
	}
	return msgID, body, nil
}

func genCalculatePacket() []byte {
	a := int32(rand.Intn(99999) + 1)
	b := int32(rand.Intn(99999) + 1)
	c := int32(0)
	var data [12]byte
	binary.BigEndian.PutUint32(data[0:4], uint32(a))
	binary.BigEndian.PutUint32(data[4:8], uint32(b))
	binary.BigEndian.PutUint32(data[8:12], uint32(c))
	return data[:]
}

func genSmallPacket() []byte {
	data := make([]byte, 16)
	binary.BigEndian.PutUint32(data[0:4], uint32(rand.Intn(10000)))
	for i := 4; i < 16; i++ {
		data[i] = byte(rand.Intn(256))
	}
	return data
}

func genMediumPacket() []byte {
	data := make([]byte, 64)
	binary.BigEndian.PutUint32(data[0:4], uint32(rand.Intn(10000)))
	binary.BigEndian.PutUint32(data[4:8], uint32(rand.Intn(100000)))
	for i := 8; i < 64; i++ {
		data[i] = byte(rand.Intn(256))
	}
	return data
}

func genLargePacket() []byte {
	data := make([]byte, 256)
	binary.BigEndian.PutUint32(data[0:4], uint32(rand.Intn(10000)))
	binary.BigEndian.PutUint32(data[4:8], uint32(rand.Intn(100000)))
	binary.BigEndian.PutUint32(data[8:12], uint32(rand.Intn(1000000)))
	for i := 12; i < 256; i++ {
		data[i] = byte(rand.Intn(256))
	}
	return data
}

func getRandomCmd() uint32 {
	switch rand.Intn(4) {
	case 0:
		return CmdCalculate
	case 1:
		return CmdSmall
	case 2:
		return CmdMedium
	case 3:
		return CmdLarge
	default:
		return CmdCalculate
	}
}

func getPacketByCmd(cmd uint32) []byte {
	switch cmd {
	case CmdCalculate:
		return genCalculatePacket()
	case CmdSmall:
		return genSmallPacket()
	case CmdMedium:
		return genMediumPacket()
	case CmdLarge:
		return genLargePacket()
	default:
		return genCalculatePacket()
	}
}

func runClient(clientID int, host string, port int, repeat int, sem chan struct{}, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		stats.addFailed(int64(repeat))
		return
	}

	reconnectCount := 0

	for i := 0; i < repeat; {
		// 随机断开重连 (5% 概率)
		if i > 0 && rand.Intn(20) == 0 {
			conn.Close()
			reconnectCount++
			// 模拟重连延迟 (10-100ms)
			time.Sleep(time.Duration(rand.Intn(90)+10) * time.Millisecond)

			conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				stats.addFailed(int64(repeat - i))
				return
			}
		}

		// 随机决定是单发还是多发
		sendCount := 1
		if rand.Intn(10) < 3 { // 30% 概率发两个包
			sendCount = 2
		}
		if i+sendCount > repeat {
			sendCount = repeat - i
		}

		sendTime := time.Now()

		// 发送多个包
		successSend := 0
		for j := 0; j < sendCount; j++ {
			cmd := getRandomCmd()
			data := getPacketByCmd(cmd)
			conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
			if err := sendMsg(conn, cmd, data); err != nil {
				conn.Close()
				reconnectCount++
				// 尝试重连
				time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)
				conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
				if err != nil {
					stats.addFailed(int64(repeat - i - j))
					return
				}
				continue
			}
			successSend++
		}

		// 接收对应数量的回包
		successRecv := 0
		for j := 0; j < sendCount; j++ {
			conn.SetReadDeadline(time.Now().Add(15 * time.Second))
			if _, _, err := recvMsg(conn); err != nil {
				conn.Close()
				reconnectCount++
				time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)
				conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
				if err != nil {
					stats.addFailed(int64(repeat - i - j))
					return
				}
				j--
				continue
			}
			successRecv++
		}

		respTime := time.Since(sendTime).Seconds() * 1000
		stats.addSuccess(int64(successRecv))
		if successRecv > 0 {
			stats.addRespTime(respTime / float64(successRecv))
		}

		i += successSend

		time.Sleep(time.Millisecond)
	}

	conn.Close()
}

func main() {
	serverHost := "localhost"
	serverPort := 8888
	totalConns := 10000
	repeatPerConn := 1000

	fmt.Printf("🚀 启动压测 v2: %d 并发连接, 每个连接请求 %d 次\n", totalConns, repeatPerConn)
	fmt.Printf("  - 命令号: 1001(12字节), 2001(16字节), 3001(64字节), 4001(256字节)\n")
	fmt.Printf("  - 模式: 随机单发或多发(30%%概率发2包)\n")
	fmt.Printf("  - 网络模拟: 随机断开重连(5%%概率), 重连延迟10-150ms\n")
	startTime := time.Now()

	stats := NewStats()
	sem := make(chan struct{}, 10000)
	var wg sync.WaitGroup

	wg.Add(totalConns)
	for i := 0; i < totalConns; i++ {
		go runClient(i, serverHost, serverPort, repeatPerConn, sem, stats, &wg)
	}

	wg.Wait()
	duration := time.Since(startTime)

	qps := float64(stats.success) / duration.Seconds()

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("🏁 压测报告 v2\n")
	fmt.Printf("总耗时: %.2f 秒\n", duration.Seconds())
	fmt.Printf("成功次数: %d\n", stats.success)
	fmt.Printf("失败次数: %d\n", stats.failed)
	fmt.Printf("有效 QPS: %.2f\n", qps)
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("响应时间 (ms):\n")
	fmt.Printf("  平均: %.2f\n", stats.getAvg())
	fmt.Printf("  最小: %.2f\n", stats.getPercentile(0))
	fmt.Printf("  P50:  %.2f\n", stats.getPercentile(0.50))
	fmt.Printf("  P90:  %.2f\n", stats.getPercentile(0.90))
	fmt.Printf("  P95:  %.2f\n", stats.getPercentile(0.95))
	fmt.Printf("  P99:  %.2f\n", stats.getPercentile(0.99))
	fmt.Printf("  最大: %.2f\n", stats.getPercentile(1.0))
	fmt.Println(strings.Repeat("=", 50))

	stats.printHistogram()
}
