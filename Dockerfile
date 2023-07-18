# 使用 golang 作为基础镜像
FROM golang:1.20-alpine as builder

# 设置工作目录
WORKDIR /app

# 将当前目录下的所有文件复制到工作目录中
COPY . .

# 编译应用程序
RUN go mod tidy
RUN go generate
RUN go build -ldflags="-s -w " -o separa

# 使用 python 官方镜像作为基础镜像
FROM python:3.9-alpine

# 设置工作目录
WORKDIR /app

# 复制
COPY --from=builder /app/separa ./
COPY ./requirements.txt ./run.py ./target.txt ./

# 安装依赖包
RUN pip install -r requirements.txt


CMD [ "python", "run.py" ]