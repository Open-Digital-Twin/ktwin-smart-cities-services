FROM python:alpine3.18

COPY ./modules /app/modules
COPY ./main.py /app/main.py
COPY ./requirements.txt /app/requirements.txt

WORKDIR /app
RUN pip install -r requirements.txt
EXPOSE 8080
ENTRYPOINT [ "python", "main.py" ]