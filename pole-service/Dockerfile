FROM python:alpine3.18
COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt
EXPOSE 8081
ENTRYPOINT [ "python", "main.py" ]