FROM ubuntu

RUN apt-get update
RUN apt-get install -y ca-certificates
ADD main /

EXPOSE 8080
CMD ["/main"]