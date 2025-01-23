FROM golang:1.23

# 設置工作目錄與本地一致
WORKDIR /ServiceExample

# 將本地文件拷貝到容器內的同一目錄
COPY . .

# 安裝依賴
RUN go mod tidy

# 編譯可執行文件
RUN go build -o main .

# 暴露服務端口
EXPOSE 8888

# 啟動服務器
CMD ["./main"]
