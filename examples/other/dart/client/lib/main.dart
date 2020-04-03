import 'dart:io';
import 'dart:typed_data';

import 'package:protobuf/protobuf.dart';

import 'pb/hello.pb.dart';

// 消息长度用4个字节描述
const int msgByteLen = 4;

// 消息号用4个字节描述
const int msgCodeByteLen = 4;

// 最小的消息长度为8个字节（即消息长度+消息号）
const int minMsgByteLen = 8;

// 网络管理器类
class NetworkManager {
  // 服务器ip
  final String host;

  // 服务器端口
  final int port;

  Socket socket;

  // 缓存的网络数据，暂未处理（一般这里有数据，说明当前接收的数据不是一个完整的消息，需要等待其它数据的到来拼凑成一个完整的消息)
  Int8List cacheData = Int8List(0);

  NetworkManager(this.host, this.port);

  // 初始化连接服务器
  void init() async {
    try {
      socket = await Socket.connect(host, port, timeout: Duration(seconds: 2));
    } catch (e) {
      print("连接socket出现异常，e=${e.toString()}");
    }
    socket.listen(decodeHandle,
        onError: errorHandler, onDone: doneHandler, cancelOnError: false);
  }

  // 解码处理方法
  // 处理服务器发过来的数据，注意，这里要处理粘包，这个data参数不一定是一个完整的包
  void decodeHandle(newData) {
    //拼凑当前最新未处理的网络数据
    cacheData = Int8List.fromList(cacheData + newData);

    //缓存数据长度符合最小包长度才尝试解码
    while (cacheData.length >= minMsgByteLen) {
      //读取消息长度
      var byteData = cacheData.buffer.asByteData();
      var msgLen = byteData.getInt32(0, Endian.little);

      //数据长度小于消息长度，说明不是完整的数据，暂不处理
      if (cacheData.length < msgLen + minMsgByteLen) {
        return;
      }
      //读取消息号
      int msgCode = byteData.getInt32(msgCodeByteLen, Endian.little);
      //读取pb数据
      Int8List pbBody;
      if (msgLen > 0) {
        pbBody = cacheData.sublist(minMsgByteLen, msgLen + minMsgByteLen);
      }

      //整理缓存数据
      int totalLen = minMsgByteLen + msgLen;
      cacheData = cacheData.sublist(totalLen, cacheData.length);

      // 处理消息
      HelloResponse response = HelloResponse.fromBuffer(pbBody);
      print(response.greeter);
    }
  }

  /**
   * 发消息，指定消息号，pb对象可以为不传(例如发心跳包的时候)
   */
  void sendMsg(int msgCode, [GeneratedMessage pb]) {
    //序列化pb对象
    Uint8List pbBody;
    int pbLen = 0;
    if (pb != null) {
      pbBody = pb.writeToBuffer();
      pbLen = pbBody.length;
    }
    //包头部分
    var header = ByteData(minMsgByteLen);
    header.setInt32(0, pbLen, Endian.little);
    header.setInt32(msgByteLen, msgCode, Endian.little);

    //包头+pb组合成一个完整的数据包
    var msg = pbBody == null
        ? header.buffer.asUint8List()
        : header.buffer.asUint8List() + pbBody.buffer.asUint8List();

    //给服务器发消息
    try {
      socket.add(msg);
      print("给服务端发送消息，消息号=$msgCode");
    } catch (e) {
      print("send捕获异常：msgCode=${msgCode}，e=${e.toString()}");
    }
  }

  void errorHandler(error, StackTrace trace) {
    print("捕获socket异常信息：error=$error，trace=${trace.toString()}");
    socket.close();
  }

  void doneHandler() {
    socket.destroy();
    print("socket关闭处理");
  }
}

/**
 * 测试
 */
main() async {
  //创建网络管理器
  var networkManager = NetworkManager("127.0.0.1", 8999);
  await networkManager.init();

  //创建登陆的pb对象并赋值
  var req = HelloRequest.create();
  req.name = "dart";

  //发送登陆请求
  networkManager.sendMsg(1, req);

  // //每秒发一次心跳请求
  // Timer.periodic(Duration(seconds: 1), (t) {
  //   netManager.sendMsg(101);
  //   lastSendHeartTime = new DateTime.now().millisecondsSinceEpoch;
  // });
}
