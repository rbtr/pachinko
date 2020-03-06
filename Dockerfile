FROM scratch
COPY bin/pachinko /pachinko
ENTRYPOINT ["/pachinko"]
CMD ["sort"]
