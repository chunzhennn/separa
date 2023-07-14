# 使用 python 官方镜像作为基础镜像
FROM python:3.9-alpine

# 设置工作目录
WORKDIR /app

# 复制
COPY ./separa .
COPY ./requirements.txt ./run.py 。/target.txt ./

# 安装依赖包
RUN pip install -r requirements.txt


CMD [ "python", "run.py" ]