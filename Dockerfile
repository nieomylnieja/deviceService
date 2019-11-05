# Device service app with GO
FROM scratch
COPY hello /usr/bin/hello
EXPOSE 80
ENTRYPOINT ["hello"]
