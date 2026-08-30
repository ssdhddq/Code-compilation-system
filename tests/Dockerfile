FROM python:3.9-slim

WORKDIR /app

COPY tests/* .

RUN pip install -r req.txt

CMD ["pytest", "test.py"]