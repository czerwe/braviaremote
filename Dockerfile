FROM scratch

# RUN mkdir -p {static,templates}
# 
ADD braviactl /braviactl
COPY static /static/
COPY templates /templates/
# ADD mapping.json.example mapping.json

ENV BRAVIAIP 10.0.0.11
ENV LISTENPORT 4043
ENV SUBNET 10.0.0.0
ENV PIN 0000
ENV MAC FC:F1:52:72:52:5F
ENV LOGLEVEL info

EXPOSE 4043

ENTRYPOINT ["/braviactl"]

# CMD ["--BRAVIAIP", "10.0.0.4"]