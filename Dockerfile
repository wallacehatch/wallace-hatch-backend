FROM ubuntu

RUN apt-get update
RUN apt-get install -y ca-certificates
ADD main /

ENV MAILCHIMP_API=0b7db42d1812635381bd6f4bdad2a771-us17

EXPOSE 8080
CMD ["/main"]