const Buffer = require("buffer").Buffer;
App({
    onLaunch() {
        const socket = wx.connectSocket({
            url: 'ws://localhost:9000',
            protocols:["12321","321321321",32132121]
        })
        wx.onSocketOpen((result) => {
            console.log("连接成功")
            wx.sendSocketMessage({
                data: this.encodeTLV(1, "hello"),


            })
        })
        socket.onMessage(result => {
            let message = this.decodeTLV(result.data)
            console.log(message)
        })
    },

    globalData: {
        TYPE_LENGTH: 4,
        LENGTH_LENGTH: 4,
    },
    // 将数据编码为 TLV 格式
    encodeTLV(type, value) {
        const length = value.length;
        const typeBuffer = Buffer.alloc(this.globalData.TYPE_LENGTH);
        const lengthBuffer = Buffer.alloc(this.globalData.LENGTH_LENGTH);
        const valueBuffer = Buffer.from(value);
        typeBuffer.writeUInt32BE(type, 0);
        lengthBuffer.writeUInt32BE(length, 0);
        return Buffer.concat([typeBuffer, lengthBuffer, valueBuffer], this.globalData.TYPE_LENGTH + this.globalData.LENGTH_LENGTH + length);
    },

    // 从 TLV 格式解码数据
    decodeTLV(buffer) {
        // 解析包头长度
        let data = Buffer.from(buffer)
        const tag = Buffer.alloc(this.globalData.TYPE_LENGTH);
        let offset = 0;
        data.copy(tag, 0, offset, 4)
        const dataLen = Buffer.alloc(this.globalData.LENGTH_LENGTH);
        offset += 4;
        data.copy(dataLen, 0, offset, offset + 4)
        // 解析数据包内容
        const body = new Buffer(dataLen.readInt32BE())
        offset += 4
        data.copy(body, 0, offset, offset + dataLen.readInt32BE());
        let message = {
            tag: tag.readUInt32BE(),
            dataLen: dataLen.readUInt32BE(),
            data: body
        }
        return message
    }
})