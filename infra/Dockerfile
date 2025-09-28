FROM python:3.13

RUN mkdir /app

WORKDIR /app

COPY requirements.txt .

RUN pip install -r requirements.txt

RUN airflow db migrate

COPY start-servers.sh .

CMD ["bash", "/app/start-servers.sh"]
