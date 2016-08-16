FROM ubuntu:14.04.4
RUN echo "hello" > /dummy.txt
RUN rm dummy.txt
CMD ["/bin/sh"]
CMD ["/bin/dash"]
RUN echo "fake" > /fake.txt
CMD ["/bin/bash"]
ENV DUMMYVALUE=dummy
RUN echo "addenv" >> /fake.txt
