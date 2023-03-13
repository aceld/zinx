# LengthFieldFrameDecoder使用详解


>LengthFieldFrameDecoder是一个基于长度字段的解码器，比较难理解的解码器，它主要有5个核心的参数配置：

>maxFrameLength：     数据包最大长度

>lengthFieldOffset：  长度字段偏移量

>lengthFieldLength：  长度字段所占的字节数

>lengthAdjustment：   长度的调整值

>initialBytesToStrip：解码后跳过的字节数




## 示例讲解

#### TLV格式协议

TLV，即Tag(Type)—Length—Value，是一种简单实用的数据传输方案。在TLV的定义中，可以知道它包括三个域，分别为：标签域（Tag），长度域（Length），内容域（Value）。这里的长度域的值实际上就是内容域的长度。

```
解码前 (20 bytes)                                   解码后 (20 bytes)
+------------+------------+-----------------+      +------------+------------+-----------------+
|     Tag    |   Length   |     Value       |----->|     Tag    |   Length   |     Value       |
| 0x00000001 | 0x0000000C | "HELLO, WORLD"  |      | 0x00000001 | 0x0000000C | "HELLO, WORLD"  |
+------------+------------+-----------------+      +------------+------------+-----------------+
```
> Tag：   uint32类型，占4字节，Tag作为MsgId，暂定为1<br>
> Length：uint32类型，占4字节，Length标记Value长度12(hex:0x0000000C)<br>
> Value： 共12个字符，占12字节<br>

```
说明：
lengthFieldOffset   = 4            (Length的字节位索引下标是4) 长度字段的偏差
lengthFieldLength   = 4            (Length是4个byte) 长度字段占的字节数
lengthAdjustment    = 0            (Length只表示Value长度，程序只会读取Length个字节就结束，后面没有来，故为0，若Value后面还有crc占2字节的话，那么此处就是2。若Length标记的是Tag+Length+Value总长度，那么此处是-8)
initialBytesToStrip = 0            (这个0表示返回完整的协议内容Tag+Length+Value，如果只想返回Value内容，去掉Tag的4字节和Length的4字节，此处就是8) 从解码帧中第一次去除的字节数
maxFrameLength      = 2^32 + 4 + 4 (Length为uint类型，故2^32次方表示Value最大长度，此外Tag和Length各占4字节)
```


#### HTLV+CRC格式协议

HTLV+CRC，H头码，T功能码，L数据长度，V数据内容


```

+------+-------+---------+--------+--------+
| 头码  | 功能码 | 数据长度 | 数据内容 | CRC校验 |
| 1字节 | 1字节  | 1字节   | N字节   |  2字节  |
+------+-------+---------+--------+--------+

```

数据示例

```
头码   功能码 数据长度      Body                         CRC
A2      10     0E        0102030405060708091011121314 050B
```

```

说明：
   1.数据长度len是14(0E),这里的len仅仅指Body长度;

   lengthFieldOffset   = 2   (len的索引下标是2，下标从0开始) 长度字段的偏差
   lengthFieldLength   = 1   (len是1个byte) 长度字段占的字节数
   lengthAdjustment    = 2   (len只表示Body长度，程序只会读取len个字节就结束，但是CRC还有2byte没读呢，所以为2)
   initialBytesToStrip = 0   (这个0表示完整的协议内容，如果不想要A2，那么这里就是1) 从解码帧中第一次去除的字节数
   maxFrameLength      = 255 + 4(起始码、功能码、CRC) (len是1个byte，所以最大长度是无符号1个byte的最大值)
       
```


## 案例分析
以下7种案例足以满足所有协议，只处理断粘包，并不能处理错包，包的完整性需要依靠协议自身定义CRC来校验

#### 案例1：
```
lengthFieldOffset  =0 长度字段从0开始
lengthFieldLength  =2 长度字段本身占2个字节
lengthAdjustment   =0 需要调整0字节
initialBytesToStrip=0 解码后跳过0字节



解码前 (14 bytes)                 解码后 (14 bytes)
+--------+----------------+      +--------+----------------+
| Length | Actual Content |----->| Length | Actual Content |
| 0x000C | "HELLO, WORLD" |      | 0x000C | "HELLO, WORLD" |
+--------+----------------+      +--------+----------------+
```

> Length为0x000C，这个是十六进制，0x000C转化十进制就是14


#### 案例2：
```
lengthFieldOffset  =0 长度字段从0开始
lengthFieldLength  =2 长度字段本身占2个字节
lengthAdjustment   =0 需要调整0字节
initialBytesToStrip=2 解码后跳过2字节

解码前 (14 bytes)                 解码后 (12 bytes)
+--------+----------------+      +----------------+
| Length | Actual Content |----->| Actual Content |
| 0x000C | "HELLO, WORLD" |      | "HELLO, WORLD" |
+--------+----------------+      +----------------+
```
>这时initialBytesToStrip字段起作用了，在解码后会将前面的2字节跳过，所以解码后就只剩余了数据部分。

#### 案例3：
```
lengthFieldOffset  =0 长度字段从0开始
lengthFieldLength  =2 长度字段本身占2个字节
lengthAdjustment   =-2 需要调整 -2 字节
initialBytesToStrip=0 解码后跳过2字节


解码前 (14 bytes)                 解码后 (14 bytes)
+--------+----------------+      +--------+----------------+
| Length | Actual Content |----->| Length | Actual Content |
| 0x000E | "HELLO, WORLD" |      | 0x000E | "HELLO, WORLD" |
+--------+----------------+      +--------+----------------+
```

>这时lengthAdjustment起作用了，因为长度字段的值包含了长度字段本身的2字节，
如果要获取数据的字节数，需要加上lengthAdjustment的值，就是 14+（-2）=12，这样才算出来数据的长度。


#### 案例4：

```
lengthFieldOffset  =2 长度字段从第2个字节开始
lengthFieldLength  =3 长度字段本身占3个字节
lengthAdjustment   =0 需要调整0字节
initialBytesToStrip=0 解码后跳过0字节


解码前 (17 bytes)                              解码后 (17 bytes)
+----------+----------+----------------+      +----------+----------+----------------+
| Header 1 |  Length  | Actual Content |----->| Header 1 |  Length  | Actual Content |
|  0xCAFE  | 0x00000C | "HELLO, WORLD" |      |  0xCAFE  | 0x00000C | "HELLO, WORLD" |
+----------+----------+----------------+      +----------+----------+----------------+
```
>由于数据包最前面加了2个字节的Header，所以lengthFieldOffset为2，
说明长度字段是从第2个字节开始的。然后lengthFieldLength为3，说明长度字段本身占了3个字节。


#### 案例5：
```
lengthFieldOffset  =0 长度字段从第0个字节开始
lengthFieldLength  =3 长度字段本身占3个字节
lengthAdjustment   =2 需要调整2字节
initialBytesToStrip=0 解码后跳过0字节


解码前 (17 bytes)                              解码后 (17 bytes)
+----------+----------+----------------+      +----------+----------+----------------+
|  Length  | Header 1 | Actual Content |----->|  Length  | Header 1 | Actual Content |
| 0x00000C |  0xCAFE  | "HELLO, WORLD" |      | 0x00000C |  0xCAFE  | "HELLO, WORLD" |
+----------+----------+----------------+      +----------+----------+----------------+
```
>lengthFieldOffset为0，所以长度字段从0字节开始。lengthFieldLength为3，长度总共占3字节。
因为长度字段后面还剩余14字节的总数据，但是长度字段的值为12，只表示了数据的长度，不包含头的长度，
所以lengthAdjustment为2，就是12+2=14，计算出Header+Content的总长度。


#### 案例6：

```
lengthFieldOffset  =1 长度字段从第1个字节开始
lengthFieldLength  =2 长度字段本身占2个字节
lengthAdjustment   =1 需要调整1字节
initialBytesToStrip=3 解码后跳过3字节

解码前 (16 bytes)                               解码后 (13 bytes)
+------+--------+------+----------------+      +------+----------------+
| HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
| 0xCA | 0x000C | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
+------+--------+------+----------------+      +------+----------------+
```
>这一次将Header分为了两个1字节的部分，lengthFieldOffset为1表示长度从第1个字节开始，lengthFieldLength为2表示长度字段占2个字节。
因为长度字段的值为12，只表示了数据的长度，所以lengthAdjustment为1，12+1=13，
表示Header的第二部分加上数据的总长度为13。因为initialBytesToStrip为3，所以解码后跳过前3个字节。


#### 案例7：
```
lengthFieldOffset  =1 长度字段从第1个字节开始
lengthFieldLength  =2 长度字段本身占2个字节
lengthAdjustment   =-3 需要调整 -3 字节
initialBytesToStrip=3 解码后跳过3字节

解码前 (16 bytes)                               解码后 (13 bytes)
+------+--------+------+----------------+      +------+----------------+
| HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
| 0xCA | 0x0010 | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
+------+--------+------+----------------+      +------+----------------+
```
>这一次长度字段的值为16，表示包的总长度，所以lengthAdjustment为 -3 ，16+ (-3)=13，
表示Header的第二部分加数据部分的总长度为13字节。initialBytesToStrip为3，解码后跳过前3个字节。

