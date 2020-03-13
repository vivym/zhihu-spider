FROM scratch

COPY zhihu-spider /

ENTRYPOINT ["/zhihu-spider"]
