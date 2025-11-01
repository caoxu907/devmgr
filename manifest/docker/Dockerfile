FROM loads/alpine:3.8

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR                 /app
ADD resource                $WORKDIR/
ADD ./temp/linux_amd64/main $WORKDIR/main
ADD ./manifest/config/config.yaml $WORKDIR/manifest/config/config.yaml
RUN chmod +x $WORKDIR/main

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./main
