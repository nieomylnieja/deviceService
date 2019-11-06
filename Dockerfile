# Device service app with GO
FROM scratch
COPY .build/hello /usr/bin/hello
EXPOSE 80
ENTRYPOINT ["hello"]
