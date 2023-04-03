const Buffer = require("buffer").Buffer;
App({
    onLaunch() {
        const socket = wx.connectSocket({
            url: 'ws://localhost:8999',
        })
        wx.onSocketOpen((result) => {
            console.log("连接成功")
            wx.sendSocketMessage({
                data: this.encodeTLV(1, "hello"),


            })
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
        const type = buffer.readUInt8(0);
        const length = buffer.readUInt16BE(this.globalData.TYPE_LENGTH);
        const value = buffer.slice(this.globalData.TYPE_LENGTH + this.globalData.LENGTH_LENGTH, this.globalData.TYPE_LENGTH + this.globalData.LENGTH_LENGTH + length);

        return {
            type,
            value
        };
    }
})