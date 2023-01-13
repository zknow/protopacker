# Proto Packer

Golang處理TCP Pack/Unpack黏包工具

- 可自訂辨識用的封包標頭
- 可註冊兩種方法來接收合法的封包結果
    - chan Byte[]
    - func([]byte)

利用固定的封包標頭辨識處理合法的封包，如果封包被截斷會暫存在buffer中與下次收到的封包一併處理。