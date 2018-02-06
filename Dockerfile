FROM ubuntu

RUN apt-get update
RUN apt-get install -y ca-certificates
ADD main /

ADD email-templates/ email-templates/ 


EXPOSE 8090
CMD ["/main"]